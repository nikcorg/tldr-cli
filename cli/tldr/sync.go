package main

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/nikcorg/tldr-cli/config/rotation"
	"github.com/nikcorg/tldr-cli/storage"
	"github.com/nikcorg/tldr-cli/sync"
)

type syncCmd struct{}

func (s *syncCmd) Init() {}

func (s *syncCmd) Verbs() []string {
	return []string{"sync"}
}

func (s *syncCmd) ParseArgs(subcommand string, args ...string) error {
	return nil
}

func (s *syncCmd) Help(subcommand string, args ...string) {
	fmt.Printf(strings.Replace(heredoc.Doc(`
		Sync the tldr log with a remote peer

		__BINARY_NAME__ sync

		Sync runs an external command and provides to it as arguments:
		- the directory containing the tldr log files
		- the log file to sync

		Currently only syncing a single log file is implemented
	`), "__BINARY_NAME__", binaryName, -1))
}

func (s *syncCmd) Execute(subcommand string, args ...string) error {
	if runtimeConfig.Sync.Exec == "" && runtimeConfig.Sync.Remote != "" {
		return fmt.Errorf("Git sync not yet implemented")
	}

	if runtimeConfig.Rotation == rotation.None {
		return s.simpleSync()
	}

	return s.multiSync()
}

func (s *syncCmd) simpleSync() error {
	syncer := sync.NewSync(runtimeConfig)

	source, err := stor.Load()
	if err != nil {
		return err
	}

	return syncer.WithCommand([]*storage.Source{source})
}

func (s *syncCmd) multiSync() error {
	return fmt.Errorf("Sync with storage rotation not yet implemented")
}
