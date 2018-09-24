package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/nikcorg/tldr-cli/archive"
	"github.com/nikcorg/tldr-cli/config"
)

// Add will create a new archive or amend the day's archive if it exists
func Add(cfg *config.Config, args ...string) {
	fmt.Printf("Add %+v, args %+v", cfg, args)
	filename := fmt.Sprintf("%s/%s.yaml", cfg.Archive, time.Now().Format(cfg.Format))

	a := archive.Archive{
		Title:   time.Now().Format("2006-01-02"),
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
