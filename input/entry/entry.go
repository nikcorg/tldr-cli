package entry

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nikcorg/tldr-cli/storage"
)

type entryMode int

const (
	addEntry entryMode = iota
	editEntry
)

// EditContext represents additional data useful for the editing context
type EditContext struct {
	Titles []string
}

// Edit takes an entry and presents it for editing via user input
func Edit(newEntry *storage.Entry, ctx *EditContext) error {
	return edit(newEntry, ctx, editEntry)
}

// Create creates a new entry
func Create(newEntry *storage.Entry, ctx *EditContext) error {
	return edit(newEntry, ctx, addEntry)
}

func edit(newEntry *storage.Entry, ctx *EditContext, mode entryMode) error {
	const (
		listTitles   = "L"
		customTitle  = "T"
		toggleUnread = "X"
		sourceURL    = "S"
		relatedURL   = "R"
		quit         = "Q"
		delete       = "D"
	)

	var err error
	var deleted = false

	title := "Add New Entry"
	if mode == editEntry {
		title = "Edit Entry"
	}

	reader := bufio.NewReader(os.Stdin)

	for true {
		fmt.Printf("=== %s ===\n", title)
		fmt.Printf("Title: %s\n", newEntry.Title)
		fmt.Printf("URL: %s\n", newEntry.URL)
		fmt.Printf("Unread: %v\n", newEntry.Unread)
		if len(newEntry.SourceURL) > 0 {
			fmt.Printf("Source: %s\n", newEntry.SourceURL)
		}
		if len(newEntry.RelatedURLs) > 0 {
			fmt.Println("Related:")
			for _, u := range newEntry.RelatedURLs {
				fmt.Printf("- %s\n", u)
			}
		}

		deleteVerb := "delete"
		if deleted {
			deleteVerb = "undelete"
		}

		fmt.Println("---")
		fmt.Println("Press Enter to accept entry and exit")
		fmt.Printf("%s) toggle unread status\n", toggleUnread)
		fmt.Printf("%s) list titles\n", listTitles)
		fmt.Printf("%s) custom title\n", customTitle)
		fmt.Printf("%s) source URL\n", sourceURL)
		fmt.Printf("%s) related URL\n", relatedURL)
		if mode == editEntry {
			fmt.Printf("%s) %s entry\n", delete, deleteVerb)
		}
		fmt.Printf("%s) quit without saving\n", quit)
		fmt.Print("Your selection: ")

		var selection string
		selection, err = reader.ReadString('\n')
		if err != nil {
			return err
		}

		selection = strings.ToUpper(strings.TrimSpace(selection))

		if len(selection) == 0 {
			break
		}

		switch selection {
		case listTitles:
			for n, t := range ctx.Titles {
				fmt.Printf("%d) %s\n", n, t)
			}
			fmt.Print("Select title or press enter to keep current title: ")
			selection, _ = reader.ReadString('\n')
			selection = strings.TrimSpace(selection)
			if len(selection) > 0 {
				n, _ := strconv.Atoi(selection)
				newEntry.Title = ctx.Titles[n]
			}

		case toggleUnread:
			newEntry.Unread = !newEntry.Unread

		case customTitle:
			fmt.Printf("Enter title: ")
			selection, _ = reader.ReadString('\n')
			newEntry.Title = strings.TrimSpace(selection)

		case sourceURL:
			fmt.Printf("Enter source: ")
			selection, _ = reader.ReadString('\n')
			newEntry.SourceURL = strings.TrimSpace(selection)

		case relatedURL:
			fmt.Printf("Enter related: ")
			selection, _ = reader.ReadString('\n')
			newEntry.RelatedURLs = append(newEntry.RelatedURLs, strings.TrimSpace(selection))

		case quit:
			fmt.Println("Ok, quitting without saving.")
			os.Exit(0)

		case delete:
			if mode == editEntry {
				deleted = !deleted
			} else {
				fmt.Printf("Delete is not applicable for new entries, use %s to discard the entry\n", quit)
			}

		default:
			fmt.Printf("I'm sorry, I don't understand '%s'. Please try again.\n", selection)
		}
	}

	newEntry.SetDeleted(deleted)

	return nil
}
