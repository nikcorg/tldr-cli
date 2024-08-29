package extract

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func attrValueFor(name string) func(n *html.Node) (string, error) {
	return func(n *html.Node) (string, error) {
		for _, attr := range n.Attr {
			if attr.Key == name {
				return attr.Val, nil
			}
		}
		return "", fmt.Errorf("Missing value attribute: %s", name)
	}
}

func textContentFor(n *html.Node) (string, error) {
	var fragments []string

	for currentNode := n.FirstChild; currentNode != nil; currentNode = currentNode.NextSibling {
		switch currentNode.Type {
		case html.ElementNode:
			sub, err := textContentFor(currentNode)
			if err != nil {
				return "", err
			}
			fragments = append(fragments, sub)

		case html.TextNode:
			fragments = append(fragments, strings.TrimSpace(currentNode.Data))
		}
	}

	if fragments == nil || len(fragments) == 0 {
		return "", errNoTextNodes
	}

	return strings.Join(compact(fragments), " "), nil
}

func compact(xs []string) []string {
	nxs := []string{}

	for _, x := range xs {
		if len(x) > 0 {
			nxs = append(nxs, x)
		}
	}

	return nxs
}
