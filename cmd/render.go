package cmd

import (
	"os"

	"github.com/nikcorg/tldr-cli/archive"
	"github.com/nikcorg/tldr-cli/config"
)

// Render is the render command
type Render struct {
	cfg *config.Config
}

// Configure configures the Render command
func (cmd Render) Configure(c *config.Config) RunnableCommand {
	return Render{
		c,
	}
}

// Run runs the Render command
func (cmd Render) Run(args ...string) {
	archives := []archive.Archive{}
	for _, file := range archive.Scan(cmd.cfg.Archive) {
		archives = append(archives, archive.Load(cmd.cfg.Archive+"/"+file))
	}

	archive.RenderToMarkdown(os.Stdout, archives)
}
