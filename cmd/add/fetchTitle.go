package add

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func retrieveTitle(href string) (string, error) {
	resp, err := http.Get(href)
	if err != nil {
		panic(fmt.Errorf("unable to fetch title from %s: %+v", href, err))
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(fmt.Errorf("unable to parse the html from %s: %v", href, err))
	}

	title := titleNode(doc)

	if title == "" {
		return "", fmt.Errorf("unable to retrieve title for %s", href)
	}

	return title, nil
}

func titleNode(n *html.Node) string {
	if n == nil {
		return ""
	} else if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tit := titleNode(c)
		if tit != "" {
			return tit
		}
	}

	return ""
}
