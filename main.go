package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nikcorg/tldr-cli/cmd"
	"github.com/nikcorg/tldr-cli/config"
)

var (
	cfg      *config.Config
	commands = cmd.CommandMap{
		"render": cmd.Render{},
		"add":    cmd.Add{},
	}
)

func init() {
	cfg = &config.Config{
		Format: "2006-01-02",
	}

	cfg.Home = os.Getenv("TLDR_HOME")

	if cfg.Home == "" {
		log.Fatal(fmt.Errorf("TLDR_HOME not found in environment"))
	}

	cfg.Archive = fmt.Sprintf("%s/archive", cfg.Home)
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		cmdArg := args[0]

		if command, found := commands[cmdArg]; found {
			command.Configure(cfg).Run(args[1:]...)
		} else {
			fmt.Printf("Command not found: %s", cmdArg)
			os.Exit(1)
		}
	} else {
		log.Fatalf("Interactive mode not yet implemented")
	}
}
