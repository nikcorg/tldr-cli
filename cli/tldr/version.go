package main

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
)

type versionCmd struct{}

const versionStr = `tldr-cli
 version : %s
 commit  : %s
 arch    : %s
 built at: %s
`

func (c *versionCmd) Verbs() []string {
	return []string{"version"}
}

func (c *versionCmd) Init() {}

func (c *versionCmd) Execute(subcommand string, args ...string) error {
	fmt.Printf(versionStr, buildVersion, buildCommit, buildArch, buildTime)

	return nil
}

func (c *versionCmd) ParseArgs(subcommand string, args ...string) error {
	return nil
}

func (c *versionCmd) Help(subcommand string, args ...string) {

	fmt.Print(strings.Replace(heredoc.Doc(`
		Show the version number and various build time details

		__BINARY_NAME__ version
	`), "__BINARY_NAME__", binaryName, -1))
}
