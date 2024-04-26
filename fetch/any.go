package fetch

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/html"

	"github.com/nikcorg/tldr-cli/extract"

	log "github.com/sirupsen/logrus"
)

const userAgent = "TL;DR cli/1.0"

func any(url string) (*Details, error) {
	var (
		err error
		res *http.Response
		req *http.Request
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, _ = http.NewRequestWithContext(ctx, "GET", url, nil)

	req.Header.Add("User-Agent", userAgent)

	if res, err = http.DefaultClient.Do(req); err != nil {
		log.Debugf("failed fetching %s: %v", url, err)
		return nil, err
	} else if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed fetching %s: %v", url, res.StatusCode)
	}

	log.Debugf("Fetched URL: %s -> %s\n", url, res.Request.URL)

	var body *html.Node
	body, err = html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	var titles, _ = extract.Titles(body)

	return &Details{res.Request.URL.String(), titles}, nil
}
