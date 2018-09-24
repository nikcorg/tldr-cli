package cmd

import (
	"os"

	"github.com/nikcorg/tldr-cli/archive"
	"github.com/nikcorg/tldr-cli/config"

	"github.com/urfave/cli"
)

// Render is the render command
type Render struct {
	cfg *config.Config
}

// Configure configures the Render command
func (cmd Render) Configure(c *config.Config) cli.Command {
	confd := Render{c}

	return cli.Command{
		Name:    "render",
		Aliases: []string{"r"},
		Action: func(ctx *cli.Context) {
			confd.Run(ctx.Args()...)
		},
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
