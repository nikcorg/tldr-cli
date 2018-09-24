package archive

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

// Save serialises an archive to disk
func Save(filename string, a Archive) {
	yml, err := yaml.Marshal(a)

	if err != nil {
		panic(err)
	}

	ensureArchiveDir(path.Dir(filename))

	if err := ioutil.WriteFile(filename, yml, 0644); err != nil {
		panic(fmt.Errorf("error writing %s: %v", filename, err))
	}
}

func ensureArchiveDir(path string) {
	if info, err := os.Stat(path); err != nil {
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(err)
		}
	} else if !info.IsDir() {
		panic(fmt.Errorf("not a directory: %s", path))
	} else if info.Mode() != 0755 {
		panic(fmt.Errorf("wrong mode set on archive dir %s", info.Mode()))
	}
}
