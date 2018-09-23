package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/nikcorg/tldr-cli/types"
	yaml "gopkg.in/yaml.v2"
)

func handleArchive(data []byte) types.Archive {
	var archive types.Archive
	yaml.Unmarshal(data, &archive)

	for x, entry := range archive.Entries {
		for y, related := range entry.Related {
			if v := types.MapRelatedEntry(related); v != nil {
				archive.Entries[x].Related[y] = v
			}
		}
	}

	return archive
}

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

func renderArchives(out io.Writer, archives []types.Archive) {
	funcMap := template.FuncMap{
		"Hello":     func() string { return "hello" },
		"IsSimple":  func(x types.RelatedEntry) bool { return x.IsSimple() },
		"IsConcise": func(x types.RelatedEntry) bool { return x.IsConcise() },
	}

	if err := template.Must(
		template.New("").Funcs(funcMap).ParseFiles("./templates/layout.tmpl"),
	).ExecuteTemplate(
		out,
		"Layout",
		struct{ Archives []types.Archive }{archives},
	); err != nil {
		log.Fatalf("Error render: %+v", err)
	}
}

const archiveDir = "../archive"

func main() {
	archives := []types.Archive{}
	for _, file := range readArchive(archiveDir) {
		if data, err := ioutil.ReadFile(archiveDir + "/" + file); err != nil {
			panic(fmt.Errorf("Error reading %s: %+v", file, err))
		} else {
			archives = append(archives, handleArchive(data))
		}
	}

	// fmt.Printf("%+v\n", archives)

	renderArchives(os.Stdout, archives)
	// for _, arch := range archives {
	// 	fmt.Println(arch.Title)
	// }
}
