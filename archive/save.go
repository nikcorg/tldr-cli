package archive

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Save serialises an archive to disk
func Save(filename string, a Archive) {
	yml, err := yaml.Marshal(a)

	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filename, yml, 0644); err != nil {
		panic(fmt.Errorf("Error writing %s: %v", filename, err))
	}
}
