package main

import "fmt"

type versionCmd struct{}

const versionStr = `tldr-cli
 version : %s
 built on: %s
 commit  : %s
 arch    : %s
`

func (c *versionCmd) Init() {}

func (c *versionCmd) Execute(subcommand string, args ...string) error {
	fmt.Printf(versionStr, buildVersion, buildDate, buildCommit, buildArch)

	return nil
}

func (c *versionCmd) ParseArgs(subcommand string, args ...string) error {
	return nil
}

func (c *versionCmd) Help(subcommand string, args ...string) {}
