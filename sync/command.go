package sync

import (
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/nikcorg/tldr-cli/storage"
	log "github.com/sirupsen/logrus"
)

// Known error outcomes
var (
	errSyncCommandNotConfigured = errors.New("no sync command is configured")
)

// WithCommand executes an external program to sync the local with a remote
func (s *Sync) WithCommand(sources []*storage.Source) error {
	if s.config.Sync.Exec == "" {
		return errSyncCommandNotConfigured
	}

	var args = []string{path.Base(s.config.Sync.Exec), s.config.Storage.Path}

	for _, source := range sources {
		args = append(args, source.SourceFile)
	}

	log.Debugf("running %s with args %v", s.config.Sync.Exec, args)

	if err := syscall.Exec(s.config.Sync.Exec, args, os.Environ()); err != nil {
		return fmt.Errorf("running sync command failed: %w", err)
	}

	return nil
}
