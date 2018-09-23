package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nikcorg/tldr-cli/archive"
)

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
	for _, file := range archive.ScanArchive(archiveDir) {
		archives = append(archives, archive.LoadFile(archiveDir+"/"+file))
	}

	archive.RenderToMarkdown(os.Stdout, archives)
}
