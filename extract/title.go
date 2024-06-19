package extract

import (
	"errors"
	"sort"
	"strings"

	"github.com/andybalholm/cascadia"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var (
	ErrNoTitles = errors.New("no titles found")
)

var selectors []titlePuller = []titlePuller{
	{"og:title", cascadia.MustCompile("meta[property=\"og:title\"]"), attrValueFor("content"), 1},
	{"twitter:title", cascadia.MustCompile("meta[property=\"twitter:title\"]"), attrValueFor("content"), 1},
	// also handle the incorrect [meta name=""] tags, as they seem fairly prevalent
	{"og:title", cascadia.MustCompile("meta[name=\"og:title\"]"), attrValueFor("content"), 1},
	{"twitter:title", cascadia.MustCompile("meta[name=\"twitter:title\"]"), attrValueFor("content"), 1},
	// -----
	{"title", cascadia.MustCompile("title"), textContentFor, 2},
	{"h1", cascadia.MustCompile("h1"), textContentFor, 3},
	{"h2", cascadia.MustCompile("h2"), textContentFor, 0},
	{"h3", cascadia.MustCompile("h3"), textContentFor, 0},
	{".title", cascadia.MustCompile(".title"), textContentFor, 0},
}

// Titles find title candidates from a root html.Node,
// and returns a ranked, unique set
func Titles(root *html.Node) ([]string, error) {
	var (
		titleCandidates                    []rankedTitle
		rankedCandidates, uniqueCandidates []string
	)

	titleCandidates, _ = getTitleCandidates(root)
	rankedCandidates, _ = rankTitleCandidates(titleCandidates)
	uniqueCandidates = uniqueTitles(rankedCandidates)

	log.Debugf("Found %d overall and %d ranked title candidates", len(titleCandidates), len(uniqueCandidates))
	for _, title := range uniqueCandidates {
		log.Debugf("- %s", title)
	}

	return uniqueCandidates, nil
}

func getTitleCandidates(res *html.Node) ([]rankedTitle, error) {
	var titles []rankedTitle = []rankedTitle{}

	for _, sel := range selectors {
		titleNode := cascadia.Query(res, sel.Selector)

		if titleNode != nil {
			titleText, err := sel.Extractor(titleNode)

			if err != nil {
				log.Debugf("Error extracting title using %s: %s", sel.Name, err.Error())
				continue
			}

			trimmedTitle := strings.TrimSpace(titleText)
			if len(trimmedTitle) > 0 {
				titles = append(titles, rankedTitle{trimmedTitle, sel.BaseScore})
				log.Debugf("Found title using %s: '%s'", sel.Name, trimmedTitle)
			}
		}
	}

	log.Debugf("Returning candidates: %v", titles)

	return titles, nil
}

const (
	exactMatch        = 3
	includesAnother   = 2 // this is most likely a site name suffixed/prefix title
	includedByAnother = 1 // this is probably the site name
)

func rankTitleCandidates(titles []rankedTitle) ([]string, error) {
	var scoredTitles []rankedTitle = []rankedTitle{}

	if len(titles) == 0 {
		return nil, ErrNoTitles
	}

	for i, t := range titles {
		rt := rankedTitle{t.Title, t.Score}
		for j, t2 := range titles {
			if i == j {
				continue
			}
			// Increase a title's rank when:
			// - It exactly matches another title
			if t == t2 {
				rt.Score += exactMatch
			}
			// - It is contained by another title
			if strings.Contains(t2.Title, t.Title) {
				rt.Score += includedByAnother
			}
			// - It contains another title
			if strings.Contains(t.Title, t2.Title) {
				rt.Score += includesAnother
			}
		}
		scoredTitles = append(scoredTitles, rt)
	}

	sort.SliceStable(scoredTitles, func(a, b int) bool {
		// Return a > b for descending rank order
		return scoredTitles[a].Score > scoredTitles[b].Score
	})

	log.Debugf("Titles scored: %+v", scoredTitles)

	// Return titles only
	rankedTitles := []string{}
	for _, t := range scoredTitles {
		rankedTitles = append(rankedTitles, t.Title)
	}

	return rankedTitles, nil
}

func uniqueTitles(allTitles []string) []string {
	if len(allTitles) == 0 {
		return allTitles
	}

	titles := []string{}
	for _, title := range allTitles {
		seen := false
		for _, t2 := range titles {
			if t2 == title {
				seen = true
				break
			}
		}
		if !seen {
			titles = append(titles, title)
		}
	}

	return titles
}
