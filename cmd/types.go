package cmd

import (
	"github.com/nikcorg/tldr-cli/config"

	"github.com/urfave/cli"
)

// CommandMap is a map of commands
type CommandMap map[string]Command

// Command is a container for configuration
type Command interface {
	Configure(*config.Config) cli.Command
}

// RunnableCommand is a command you can run
type RunnableCommand interface {
	Run(args ...string)
}
