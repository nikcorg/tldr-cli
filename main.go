package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nikcorg/tldr-cli/cmd"
	"github.com/nikcorg/tldr-cli/config"

	"github.com/urfave/cli"
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
		Format:      "2006-01-02",
		TitleFormat: "2006-01-02",
	}

	cfg.Home = os.Getenv("TLDR_HOME")

	if cfg.Home == "" {
		log.Fatal(fmt.Errorf("TLDR_HOME not found in environment"))
	}

	cfg.Archive = fmt.Sprintf("%s/archive", cfg.Home)
}

func main() {
	if len(os.Args) > 0 {
		app := cli.NewApp()
		app.Commands = []cli.Command{}
		for _, cmd := range commands {
			app.Commands = append(app.Commands, cmd.Configure(cfg))
		}
		app.Run(os.Args)
	} else {
		log.Fatalf("Interactive mode not yet implemented")
	}
}
