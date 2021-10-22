package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/nikcorg/tldr-cli/utils"
	log "github.com/sirupsen/logrus"
)

var (
	errInvalidArg        = fmt.Errorf("invalid argument")
	errExpectedNumberArg = fmt.Errorf("expected numeric argument")
)

type listCmd struct {
	num       int
	offset    int
	newerThan *time.Time
}

func (f *listCmd) Verbs() []string {
	return []string{"list", "show"}
}

func (f *listCmd) Init() {
	f.num = runtimeConfig.List.PageSize
	f.offset = 0
}

func (f *listCmd) ParseArgs(subcommand string, args ...string) error {
	argsCopy := args[0:]

	for len(argsCopy) > 0 {
		arg := argsCopy[0]

		switch strings.ToLower(arg) {
		case "-t", "--today":
			f.newerThan = utils.Today()
			log.Debugf("newer than %v", f.newerThan)

		case "-n", "--num":
			num, err := strconv.Atoi(argsCopy[1])
			if err != nil {
				return fmt.Errorf("%w: %s", errExpectedNumberArg, argsCopy[1])
			}
			f.num = num

			// Shift args
			argsCopy = argsCopy[1:]

		case "-o", "--offset", "--skip":
			offset, err := strconv.Atoi(argsCopy[1])
			if err != nil {
				return fmt.Errorf("%w: %s", errExpectedNumberArg, argsCopy[1])
			}
			f.offset = offset

			// Shift args
			argsCopy = argsCopy[1:]

		default:
			return fmt.Errorf("%w: %s", errInvalidArg, arg)
		}

		// Shift args
		argsCopy = argsCopy[1:]
	}

	if f.newerThan == nil && f.num < 0 {
		return fmt.Errorf("unlimited listing not yet supported")
	}

	return nil
}

func (f *listCmd) Execute(subcommand string, args ...string) error {
	log.Debugf("list:%s, args=%v", subcommand, args)

	source, err := stor.Load()
	if err != nil {
		return err
	}

	displayed := 0
	skipped := 0
	to_skip := f.num * f.offset

	for _, d := range source.Records {
		if f.newerThan != nil && !d.Date.Equal(*f.newerThan) && !d.Date.After(*f.newerThan) {
			log.Debugf("%v < %v", d.Date, f.newerThan)
			break
		}

		for i := len(d.Entries) - 1; i >= 0 && (f.num < 0 || displayed < f.num); i-- {
			e := d.Entries[i]
			if skipped < to_skip {
				skipped++
				continue
			}

			fmt.Printf("👉 %v, %+v\n", d.Date, e)

			displayed++
		}

		if f.num > -1 && displayed >= f.num {
			break
		}
	}

	return nil
}

func (f *listCmd) Help(subcommand string, args ...string) {
	fmt.Print(strings.Replace(heredoc.Doc(`
		Show previous entries

		__BINARY_NAME__ show [-n <n>] [-o <n>] [-t]
		__BINARY_NAME__ list [-n <n>] [-o <n>] [-t]

		-t, --today            entries added on the current date
		-o, --offset, --skip   skip <n> pages
		-n, --num              show <n> items per page
	`), "__BINARY_NAME__", binaryName, -1))
}
