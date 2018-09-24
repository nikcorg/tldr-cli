package archive

import (
	"os"
)

// Exists checks whether an archive file exists
func Exists(name string) bool {
	_, err := os.Stat(name)

	return err == nil
}
