package add

import (
	"fmt"
	"strings"
	"time"

	"github.com/nikcorg/tldr-cli/archive"
	"github.com/nikcorg/tldr-cli/config"

	"github.com/urfave/cli"
)

// Add is the Add command
type Add struct {
	cfg   *config.Config
	flags flags
}

type flags struct {
	Unread  bool
	Source  string
	Related []string
}

// Configure configures the Add command
func (cmd Add) Configure(c *config.Config) cli.Command {
	confd := Add{cfg: c}

	return cli.Command{
		Name:      "add",
		Aliases:   []string{"a"},
		ArgsUsage: "URL Anything that follows becomes the title",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "unread, u", Usage: "mark entry as unread"},
			cli.StringFlag{Name: "source, s", Usage: "the source `URL`"},
			cli.StringSliceFlag{Name: "related, r", Usage: "a related `URL`, as many as you like"},
		},
		Action: func(ctx *cli.Context) {
			confd.flags = flags{
				Unread:  ctx.Bool("unread"),
				Source:  ctx.String("source"),
				Related: ctx.StringSlice("related"),
			}

			href := assertIsHref(ctx.Args().First())
			title := strings.Join(ctx.Args().Tail(), " ")

			if len(title) == 0 {
				fetchedTitle, err := retrieveTitle(href)
				if err != nil {
					panic(err)
				}

				title = fetchedTitle
			}

			confd.Run(href, title)
		},
	}
}

// Run runs the Add command
func (cmd Add) Run(url string, title string) {
	filename := fmt.Sprintf("%s/%s.yaml", cmd.cfg.Archive, time.Now().Format(cmd.cfg.Format))

	a := archive.Archive{
		Title:   time.Now().Format(cmd.cfg.TitleFormat),
		Entries: []archive.Entry{},
	}

	if archive.Exists(filename) {
		a = archive.Load(filename)
	}

	entry := archive.Entry{
		Title:   title,
		URL:     url,
		Unread:  cmd.flags.Unread,
		Related: toSimpleEntrySlice(cmd.flags.Related),
		Source:  cmd.flags.Source,
	}
	a.Entries = append(a.Entries, entry)

	archive.Save(filename, a)

	fmt.Printf("Title: %s\n  URL: %s\n", entry.Title, entry.URL)
}

func toSimpleEntrySlice(xs []string) []archive.RelatedEntry {
	es := make([]archive.RelatedEntry, len(xs))
	for i, url := range xs {
		es[i] = archive.SimpleEntry{URL: url}
	}
	return es
}
