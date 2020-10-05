package main

import "fmt"

type versionCmd struct{}

func (c *versionCmd) Init() {}

func (c *versionCmd) Execute(subcommand string, args ...string) error {
	fmt.Printf("tldr-cli version %s, built on %s\n", buildVersion, buildDate)

	return nil
}

func (c *versionCmd) ParseArgs(subcommand string, args ...string) error {
	return nil
}

func (c *versionCmd) Help(subcommand string, args ...string) {}
