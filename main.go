package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/nikcorg/tldr-cli/archive"
)

type filesByName []os.FileInfo

func (x filesByName) Len() int {
	return len(x)
}

func (x filesByName) Swap(a, b int) {
	x[a], x[b] = x[b], x[a]
}

func (x filesByName) Less(a, b int) bool {
	return x[a].Name() < x[b].Name()
}

func readArchive(archiveDir string) []string {
	if infos, err := ioutil.ReadDir(archiveDir); err != nil {
		panic(fmt.Errorf("Error reading archive dir: %+v", err))
	} else {
		sort.Sort(sort.Reverse(filesByName(infos)))
		names := []string{}
		for _, info := range infos {
			names = append(names, info.Name())
		}
		return names
	}
}

func renderArchives(out io.Writer, archives []archive.Archive) {
	funcMap := template.FuncMap{
		"IsSimple":  func(x archive.RelatedEntry) bool { return x.IsSimple() },
		"IsConcise": func(x archive.RelatedEntry) bool { return x.IsConcise() },
	}

	if err := template.Must(
		template.New("").Funcs(funcMap).ParseFiles("./templates/layout.tmpl"),
	).ExecuteTemplate(
		out,
		"Layout",
		struct{ Archives []archive.Archive }{archives},
	); err != nil {
		log.Fatalf("Error render: %+v", err)
	}
}

var (
	tldrHome   string
	archiveDir string
)

func init() {
	tldrHome := os.Getenv("TLDR_HOME")

	if tldrHome == "" {
		log.Fatal(fmt.Errorf("TLDR_HOME not found in environment"))
	}

	archiveDir = tldrHome + "/archive"
}

func main() {
	archives := []archive.Archive{}
	for _, file := range readArchive(archiveDir) {
		archives = append(archives, archive.LoadFile(archiveDir+"/"+file))
	}

	renderArchives(os.Stdout, archives)
}
