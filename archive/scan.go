package archive

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type stringSlice []string

func (xs stringSlice) Len() int {
	return len(xs)
}

func (xs stringSlice) Swap(a, b int) {
	xs[a], xs[b] = xs[b], xs[a]
}

func (xs stringSlice) Less(a, b int) bool {
	return xs[a] < xs[b]
}

// Scan returns a list of yaml files in the archive directory
func Scan(archiveDir string) []string {
	if infos, err := ioutil.ReadDir(archiveDir); err != nil {
		panic(fmt.Errorf("Error reading archive dir: %+v", err))
	} else {
		names := []string{}
		for _, info := range infos {
			if name := info.Name(); strings.HasSuffix(name, "yaml") {
				names = append(names, name)
			}
		}

		sort.Sort(sort.Reverse(stringSlice(names)))

		return names
	}
}
