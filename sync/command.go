package sync

import (
	"syscall"

	"github.com/nikcorg/tldr-cli/storage"
)

// WithCommand executes an external program to sync the local with a remote
func (s *Sync) WithCommand(sources []*storage.Source) error {
	var args []string = []string{s.config.Storage.Path}

	for _, source := range sources {
		args = append(args, source.SourceFile)
	}

	if err := syscall.Exec(s.config.Sync.Exec, args, []string{}); err != nil {
		return err
	}

	return nil
}
