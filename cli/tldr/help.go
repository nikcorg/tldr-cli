package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

type helpCmd struct {
	allVerbs []string
}

func (c *helpCmd) Verbs() []string {
	return []string{"help"}
}

func (c *helpCmd) Init() {
	c.allVerbs = []string{}

	for _, cc := range commands {
		c.allVerbs = append(c.allVerbs, cc.Verbs()[0])
	}

	sort.Strings(c.allVerbs)
}

func (c *helpCmd) Execute(subcommand string, args ...string) error {
	log.Debugf("help:%s, args=%v", subcommand, strings.Join(args, "|"))

	if len(args) > 0 && args[0] != "" {
		firstArg := args[0]
		runnableCommand, helpFocus, helpFocusSubcommand, restArgs := runnableForCommand(firstArg, args[1:])

		log.Debugf("focused help: %s:%s, %+v", helpFocus, helpFocusSubcommand, runnableCommand)

		runnableCommand.Help(helpFocusSubcommand, restArgs...)
	} else {
		c.Help("")
	}

	return nil
}

func (c *helpCmd) ParseArgs(subcommand string, args ...string) error {
	return nil
}

func (c *helpCmd) Help(subcommand string, args ...string) {
	fmt.Printf("Usage: %s <options> <command>\n\n", binaryName)
	fmt.Printf(
		"Available commands: %s\n\n",
		strings.Join(c.allVerbs, ", "),
	)
	fmt.Println(strings.ReplaceAll(
		"Use __BINARY_NAME__ help <command> for more information on each command",
		"__BINARY_NAME__",
		binaryName,
	))

	fmt.Println()
	fmt.Printf("Options for %s\n", binaryName)
	flag.PrintDefaults()
}
