package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var search string

// Delete will delete all files that match a specified path.
func Delete(path string, info os.FileInfo, err error) error {
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if strings.Contains(p, search) {

		fmt.Printf("delete: %s\n", p)

		return os.Remove(p)
	}

	return nil
}

// SetSearch sets the match string to use for strings.Contains.
func SetSearch(s string) {
	search = s
}
