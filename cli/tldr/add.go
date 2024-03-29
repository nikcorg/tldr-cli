package main

import (
	"fmt"
	"strings"

	"github.com/nikcorg/tldr-cli/fetch"
	"github.com/nikcorg/tldr-cli/input/entry"
	"github.com/nikcorg/tldr-cli/storage"
	"github.com/nikcorg/tldr-cli/touchable"
	"github.com/nikcorg/tldr-cli/utils"

	"github.com/MakeNowJust/heredoc"

	log "github.com/sirupsen/logrus"
)

var (
	errURLNotFound error = fmt.Errorf("no URL argument found")
)

type addCmd struct {
	interactive *touchable.Bool
	url         *touchable.String
	relatedURLs []string
	sourceURL   *touchable.String
	title       *touchable.String
	unread      *touchable.Bool
}

func (c *addCmd) Verbs() []string {
	return []string{"add", "amend"}
}

func (c *addCmd) Execute(subcommand string, args ...string) error {
	log.Debugf("add:%s, args=%v", subcommand, args)

	source, err := stor.Load()
	if err != nil {
		return err
	}

	switch subcommand {
	case "title":
		err = amendTitle(source, c.url.Val())

	case "source":
		err = amendSource(source, c.url.Val())

	case "related":
		err = amendRelated(source, c.url.Val())

	default:
		err = c.addEntry(source)
	}

	if err != nil {
		return err
	}

	log.Debugf("before save: %+v", source)

	err = stor.Save(source)
	if err != nil {
		return err
	}

	return nil
}

func (c *addCmd) Init() {
	c.interactive = touchable.NewBool(false)
	c.title = touchable.NewString("")
	c.sourceURL = touchable.NewString("")
	c.relatedURLs = []string{}
	c.unread = touchable.NewBool(true)
	c.url = touchable.NewString("")
}

func (c *addCmd) Help(subcommand string, args ...string) {
	fmt.Print(strings.Replace(heredoc.Doc(`
		Add a new url to the log or amend the previous added entry

		__BINARY_NAME__ add <url> [-i] [-x] [-s <url>] [-r <url>]

		-i, --interactive   interactive mode
		-x, --read          mark as read
		-s, --source        source url
		-r, --related       related url, can be supplied multiple times

		Add (as amend) also provides a non-interactive form of amending the previous entry

		__BINARY_NAME__ amend:related <url>
		__BINARY_NAME__ amend:source <url>
		__BINARY_NAME__ amend:title <title>
	`), "__BINARY_NAME__", binaryName, -1))
}

func (c *addCmd) ParseArgs(subcommand string, args ...string) error {
	urlFound := false
	argsCopy := args[0:]

	shift := func() {
		argsCopy = argsCopy[1:]
	}

	for len(argsCopy) > 0 {
		arg := argsCopy[0]

		if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			switch arg {
			case "-i", "--interactive":
				c.interactive.Set(true)

			case "-r", "-rel", "--related":
				shift()
				url := argsCopy[0]
				log.Debugf("found -r: %s", url)
				if url == "" {
					return errURLNotFound
				} else if !strings.HasPrefix(url, "http") {
					return fmt.Errorf("%w: %s", errInvalidArg, url)
				}
				c.relatedURLs = append(c.relatedURLs, url)

			case "-s", "--source":
				shift()
				url := argsCopy[0]
				log.Debugf("found -s: %s", url)
				if url == "" {
					return errURLNotFound
				} else if !strings.HasPrefix(url, "http") {
					return fmt.Errorf("%w: %s", errInvalidArg, url)
				}
				c.sourceURL.Set(url)

			case "-t", "--title":
				c.title.Set(arg)

			case "-x", "--read":
				log.Debug("setting unread=false")
				c.unread.Set(false)

			default:
				return fmt.Errorf("%w: %s", errUnknownArg, arg)
			}
		} else if strings.HasPrefix(arg, "http") {
			urlFound = true
			c.url.Set(arg)
		}

		shift()
	}

	if !urlFound && subcommand == "" {
		return errURLNotFound
	}

	return nil
}

///

func (c *addCmd) addEntry(source *storage.Source) error {
	log.Debugf("Fetching %v", c.url.Val())

	var res *fetch.Details
	var err error
	if res, err = fetch.Preview(c.url.Val()); err != nil {
		return fmt.Errorf("error fetching (%s): %w", c.url.Val(), err)
	}
	log.Debugf("Fetch result: %+v", res)

	if len(res.Titles) > 0 {
		c.title.SetUnlessTouched(res.Titles[0])
	}

	var newEntry = &storage.Entry{
		URL:         res.URL,
		Title:       c.title.Val(),
		Unread:      c.unread.ValOrDefault(true),
		RelatedURLs: c.relatedURLs,
		SourceURL:   c.sourceURL.Val(),
	}

	if c.interactive.ValOrDefault(false) {
		entry.Create(newEntry, &entry.EditContext{Titles: res.Titles})
	} else {
		fmt.Printf("Title: %s\n", newEntry.Title)
		fmt.Printf("URL: %s\n", newEntry.URL)
		fmt.Printf("Unread: %v\n", newEntry.Unread)
		if newEntry.SourceURL != "" {
			fmt.Printf("Source: %s\n", newEntry.SourceURL)
		}
		if len(newEntry.RelatedURLs) > 0 {
			fmt.Println("Related:")
			for _, u := range newEntry.RelatedURLs {
				fmt.Printf("- %s\n", u)
			}
		}
	}

	addEntryToTLDR(newEntry, source)

	return nil
}

func addEntryToTLDR(newEntry *storage.Entry, source *storage.Source) {
	today := *utils.Today()

	if source.Size() == 0 {
		source.Records = []*storage.Record{
			{
				Date:    today,
				Entries: []storage.Entry{*newEntry},
			},
		}
		return
	}

	lastEntryDate := source.Records[0].Date

	if lastEntryDate.Equal(today) {
		log.Debug("Entry for today already exists, appending")
		source.Records[0].Entries = append(source.Records[0].Entries, *newEntry)
	} else {
		log.Debug("Entry for today doesn't exist, creating")
		newRecords := append([]*storage.Record{
			{
				Date:    today,
				Entries: []storage.Entry{*newEntry},
			},
		}, source.Records...)
		source.Records = newRecords
	}
}

///

func amendRelated(source *storage.Source, url string) error {
	r := source.Records[0]
	e := &r.Entries[len(r.Entries)-1]

	e.RelatedURLs = append(e.RelatedURLs, url)

	return nil
}

///

func amendSource(source *storage.Source, url string) error {
	r := source.Records[0]
	e := &r.Entries[len(r.Entries)-1]

	e.SourceURL = url

	return nil
}

func amendTitle(source *storage.Source, title string) error {
	r := source.Records[0]
	e := &r.Entries[len(r.Entries)-1]

	e.Title = title

	return nil
}
