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

func render() {
	archives := []archive.Archive{}
	for _, file := range archive.ScanArchive(archiveDir) {
		archives = append(archives, archive.LoadFile(archiveDir+"/"+file))
	}

	archive.RenderToMarkdown(os.Stdout, archives)
}

func init() {
	tldrHome := os.Getenv("TLDR_HOME")

	if tldrHome == "" {
		log.Fatal(fmt.Errorf("TLDR_HOME not found in environment"))
	}

	archiveDir = tldrHome + "/archive"
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		cmd := args[0]

		switch cmd {
		case "render":
			render()
			break
		default:
			log.Fatalf("Unknown command: %s", args[0])
		}
	} else {
		log.Fatalf("Interactive mode not yet implemented")
	}
}
