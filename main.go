package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nikcorg/tldr-cli/cmd"
)

var config config.Config

func init() {
	config := config.Config{}

	config.Home = os.Getenv("TLDR_HOME")

	if config.home == "" {
		log.Fatal(fmt.Errorf("TLDR_HOME not found in environment"))
	}

	config.Archive = fmt.Sprintf("%s/archive", config.Home)
	config.Format
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		cmdArg := args[0]

		switch cmdArg {
		case "render":
			cmd.Render(config, args...)
			break
		default:
			log.Fatalf("Unknown command: %s", args[0])
		}
	} else {
		log.Fatalf("Interactive mode not yet implemented")
	}
}
