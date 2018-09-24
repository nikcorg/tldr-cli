package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/nikcorg/tldr-cli/archive"
	"github.com/nikcorg/tldr-cli/config"
)

// Add is the Add command
type Add struct {
	cfg *config.Config
}

// Configure configures the Add command
func (cmd Add) Configure(c *config.Config) RunnableCommand {
	return Add{c}
}

// Run runs the Add command
func (cmd Add) Run(args ...string) {
	fmt.Printf("Add %+v, args %+v", cmd.cfg, args)
	filename := fmt.Sprintf("%s/%s.yaml", cmd.cfg.Archive, time.Now().Format(cmd.cfg.Format))

	a := archive.Archive{
		Title:   time.Now().Format(cmd.cfg.TitleFormat),
		Entries: []archive.Entry{},
	}

	if archive.Exists(filename) {
		a = archive.Load(filename)
	}

	a.Entries = append(a.Entries, archive.Entry{
		Title:  strings.Join(args[1:], " "),
		URL:    args[0],
		Unread: true,
	})

	archive.Save(filename, a)
}
