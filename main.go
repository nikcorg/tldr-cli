package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nikcorg/tldr-cli/cmd"
	"github.com/nikcorg/tldr-cli/config"
)

var cfg *config.Config

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

		switch cmdArg {
		case "render":
			cmd.Render(cfg, args[1:]...)
			break
		default:
			log.Fatalf("Unknown command: %s", args[0])
		}
	} else {
		log.Fatalf("Interactive mode not yet implemented")
	}
}
