package sync

import (
	"github.com/nikcorg/tldr-cli/config"
)

// Sync represents a mechanism for syncing the local database with a remote
type Sync struct {
	config *config.Settings
}

// NewSync constructs a Sync
func NewSync(config *config.Settings) Sync {
	return Sync{config}
}
