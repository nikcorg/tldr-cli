package cmd

import (
	"os"

	"github.com/nikcorg/tldr-cli/archive"
	"github.com/nikcorg/tldr-cli/config"
)

// Render beep boop
func Render(cfg *config.Config, args ...string) {
	archives := []archive.Archive{}
	for _, file := range archive.Scan(cfg.Archive) {
		archives = append(archives, archive.Load(cfg.Archive+"/"+file))
	}

	archive.RenderToMarkdown(os.Stdout, archives)
}
