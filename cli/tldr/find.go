package main

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/nikcorg/tldr-cli/storage"

	log "github.com/sirupsen/logrus"
)

type findCmd struct {
	filters     []findFilter
	needle      string
	showRelated bool
}

type findFilter func(entry *storage.Entry) bool

func unreadFilter(state bool) findFilter {
	return func(e *storage.Entry) bool {
		return e.Unread == state
	}
}

var (
	errUnreadReadConflict  = fmt.Errorf("cannot use both --unread and --read")
	errArgsAfterSearchTerm = fmt.Errorf("cannot pass flags after the search term")
	errUnknownArg          = fmt.Errorf("unknown argument")
	errInvalidArgument     = fmt.Errorf("invalid argument")
	errJunkAfterNeedle     = fmt.Errorf("found junk after search term")
	errMissingNeedle       = fmt.Errorf("no search term found")
)

func (c *findCmd) Verbs() []string {
	return []string{"find"}
}

func (c *findCmd) Init() {
	c.filters = []findFilter{}
	c.needle = ""
	c.showRelated = false
}

func (c *findCmd) ParseArgs(subcommand string, args ...string) error {
	seenUnread := false
	seenRead := false
	needleFound := false

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			if needleFound {
				return errArgsAfterSearchTerm
			}

			switch arg {
			case "-u", "--unread":
				if seenRead {
					return errUnreadReadConflict
				}
				c.filters = append(c.filters, unreadFilter(true))
				seenUnread = true

			case "-r", "--read":
				if seenUnread {
					return errUnreadReadConflict
				}
				c.filters = append(c.filters, unreadFilter(false))
				seenUnread = true

			case "-rel", "--related":
				log.Debugf("related, Set `showRelated=true`")
				c.showRelated = true

			default:
				log.Debugf("Unknown argument: %s", arg)
				return fmt.Errorf("%w: %s", errUnknownArg, arg)
			}
		} else if !needleFound {
			c.needle = strings.ToLower(arg)
			log.Debugf("found needle %s, %s", arg, c.needle)
			needleFound = true
		} else {
			return fmt.Errorf("%w: %s", errInvalidArgument, arg)
		}
	}

	if !needleFound {
		return errMissingNeedle
	}

	return nil
}

func (c *findCmd) Execute(subcommand string, args ...string) error {
	log.Debugf("Needle: %v", c.needle)
	log.Debugf("Find filters: %+v", c.filters)

	source, err := stor.Load()
	if err != nil {
		log.Debugf("Unhandled error loading storage: %v", err)
		return err
	}

	switch subcommand {
	case "one", "first":
		c.findFirst(source)

	case "all":
		fallthrough

	default:
		c.findAll(source)
	}

	return nil
}

func (c *findCmd) Help(subcommand string, args ...string) {
	fmt.Print(strings.Replace(heredoc.Doc(`
		Find an existing entry

		__BINARY_NAME__ find <search term>
		__BINARY_NAME__ find:all <search term>

		Return all entries matching the search term

		__BINARY_NAME__ find:one <search term>
		__BINARY_NAME__ find:first <search term>

		Return only the first match.
	`), "__BINARY_NAME__", binaryName, -1))
}

type searchResult struct {
	Entry  *storage.Entry
	Record *storage.Record
}

func filtersMatch(entry *storage.Entry, filters *[]findFilter) bool {
	if filters == nil || len(*filters) == 0 {
		return true
	}

	for _, f := range *filters {
		if !f(entry) {
			return false
		}
	}

	return true
}

func locateMatches(stor []*storage.Record, needle string, stopAfter int, filters *[]findFilter) []searchResult {
	searched := 0
	results := []searchResult{}
	for _, record := range stor {
		for _, entry := range record.Entries {
			searched++
			if entry.Contains(needle) && filtersMatch(&entry, filters) {
				log.Debugf("Found needle (%s) in %+v added on %v", needle, entry, record.Date)
				e := entry
				r := record
				results = append(results, searchResult{&e, r})

				if stopAfter > 0 && len(results) >= stopAfter {
					return results
				}
			}
		}
	}

	if len(results) == 0 {
		log.Debugf("Searched %d records, but found no match for needle '%s'", searched, needle)
	} else {
		log.Debugf("Searched %d records, found %d matches", searched, len(results))
	}

	return results
}

func (c *findCmd) noMatches(needle string) {
	fmt.Printf("No match found for \"%s\"\n", needle)
}

func (c *findCmd) oneMatch(sr searchResult, needle string) {
	entry := sr.Entry

	log.Debugf("Showing entry: %+v, Related: %v", entry, c.showRelated)

	if !entry.Unread {
		fmt.Printf("[x] %s\n%s\n", entry.Title, entry.URL)
	} else {
		fmt.Printf("[ ] %s\n%s\n", entry.Title, entry.URL)
	}

	if c.showRelated && len(entry.RelatedURLs) > 0 {
		fmt.Print("See also:\n")
		for _, rel := range entry.RelatedURLs {
			fmt.Printf("- %s\n", rel)
		}
	}
}

func (c *findCmd) findAll(source *storage.Source) {
	needle := c.needle
	filters := &c.filters

	results := locateMatches(source.Records, needle, 0, filters)
	if len(results) == 0 {
		c.noMatches(c.needle)
		return
	}

	matches := "match"
	if len(results) > 1 {
		matches = "matches"
	}

	fmt.Printf("Found %d %s for \"%s\"\n", len(results), matches, needle)
	for _, rs := range results {
		c.oneMatch(rs, needle)
	}
}

func (c *findCmd) findFirst(source *storage.Source) {
	needle := c.needle
	filters := &c.filters

	results := locateMatches(source.Records, needle, 1, filters)
	if len(results) == 0 {
		c.noMatches(needle)
		return
	}

	fmt.Printf("Found match for \"%s\" from %s\n", needle, results[0].Record.Date)
	c.oneMatch(results[0], needle)
}
