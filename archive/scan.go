package archive

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type filesByName []os.FileInfo

func (x filesByName) Len() int {
	return len(x)
}

func (x filesByName) Swap(a, b int) {
	x[a], x[b] = x[b], x[a]
}

func (x filesByName) Less(a, b int) bool {
	return x[a].Name() < x[b].Name()
}

// ScanArchive returns a list of yaml files in the archive directory
func ScanArchive(archiveDir string) []string {
	if infos, err := ioutil.ReadDir(archiveDir); err != nil {
		panic(fmt.Errorf("Error reading archive dir: %+v", err))
	} else {
		sort.Sort(sort.Reverse(filesByName(infos)))
		names := []string{}
		for _, info := range infos {
			if name := info.Name(); strings.HasSuffix(name, "yaml") {
				names = append(names, name)
			}
		}
		return names
	}
}
