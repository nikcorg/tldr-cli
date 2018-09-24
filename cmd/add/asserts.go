package add

import (
	"fmt"
	"net/url"
)

func assertIsHref(href string) string {
	if v, e := url.Parse(href); e != nil || !v.IsAbs() {
		panic(fmt.Errorf("Trying to add non-absolute URL: %s", href))
	} else {
		return v.String()
	}
}
