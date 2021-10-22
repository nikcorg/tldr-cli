package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/nikcorg/tldr-cli/config"
	"github.com/nikcorg/tldr-cli/storage"

	log "github.com/sirupsen/logrus"
)

// Build variables
var (
	buildArch    = "<none>"
	buildCommit  = "<none>"
	buildTime    = "<none>"
	buildVersion = "<none>"
	binaryName   = "tldr"
)

var (
	configFile = ""
	sourceDir  = ""
	sourceFile = ""

	debugLogging   = false
	verboseLogging = false

	configWasLoadedFromDisk = false
	runtimeConfig           *config.Settings

	stor *storage.Storage

	cmdAdd   = &addCmd{}
	cmdHelp  = &helpCmd{}
	commands = []runnable{
		cmdAdd,
		cmdHelp,
		&configCmd{},
		&editCmd{},
		&findCmd{},
		&listCmd{},
		&versionCmd{},
		&syncCmd{},
	}
)

func main() {
	handleFlags()
	setLogLevel()

	runtimeConfig = config.NewWithDefaults()
	stor = storage.New(runtimeConfig)

	if err := mainWithErr(flag.Args()...); err != nil {
		log.Fatalf("Error running cmd: %s", err.Error())
	}
}

func splitCommand(cmd string) (string, string) {
	if cmd == "" {
		return "", ""
	} else if !strings.Contains(cmd, ":") {
		return cmd, ""
	}
	cmds := strings.SplitN(cmd, ":", 2)

	return cmds[0], cmds[1]
}

func runnableForCommand(firstArg string, args []string) (runnable, string, string, []string) {
	var (
		runnableCommand runnable
		nextArgs        = args
	)

	command, subcommand := splitCommand(firstArg)

	for _, c := range commands {
		for _, v := range c.Verbs() {
			if v == command {
				runnableCommand = c
				break
			}
		}
	}

	if runnableCommand == nil {
		subcommand = ""
		runnableCommand, nextArgs = defaultRunnable(firstArg, args)
	}

	return runnableCommand, command, subcommand, nextArgs
}

func defaultRunnable(firstArg string, args []string) (runnable, []string) {
	if strings.HasPrefix(firstArg, "http") {
		return cmdAdd, append([]string{firstArg}, args...)
	}

	return cmdHelp, append([]string{firstArg}, args...)
}

func mainWithErr(args ...string) error {
	var err error
	if err = runtimeConfig.Load(configFile); err != nil && err != config.ErrConfigFileNotFound {
		return err
	}

	configWasLoadedFromDisk = err != config.ErrConfigFileNotFound

	log.Debugf("Runtime config after Load (from disk? %v) %+v", configWasLoadedFromDisk, runtimeConfig)

	firstArg := ""
	restArgs := []string{}

	if len(args) > 0 {
		firstArg = args[0]
	}

	if len(args) > 1 {
		restArgs = args[1:]
	}

	runnableCommand, command, subcommand, cmdArgs := runnableForCommand(firstArg, restArgs)

	runnableCommand.Init()

	if err = runnableCommand.ParseArgs(subcommand, cmdArgs...); err != nil {
		return fmt.Errorf("%w: %s", errInvalidArg, err)
	}

	if err = runnableCommand.Execute(subcommand, cmdArgs...); err != nil {
		if subcommand != "" {
			return fmt.Errorf("error running %s:%s: %w", command, subcommand, err)

		}
		return fmt.Errorf("error running %s: %w", command, err)
	}

	return nil
}

func handleFlags() {
	flag.StringVar(&configFile, "c", "", "Override config file")
	flag.StringVar(&sourceDir, "d", "", "Override storage location")
	flag.StringVar(&sourceFile, "f", "tldr.yaml", "Override storage file name (stem)")
	flag.BoolVar(&verboseLogging, "v", false, "Show verbose output")
	flag.BoolVar(&debugLogging, "vv", false, "Show debug output")
	flag.Parse()
}

func setLogLevel() {
	log.SetLevel(log.ErrorLevel)
	if debugLogging {
		log.SetLevel(log.DebugLevel)
	} else if verboseLogging {
		log.SetLevel(log.InfoLevel)
	}
}
