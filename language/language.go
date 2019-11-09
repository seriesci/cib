package language

import (
	"errors"
	"os"
	"path/filepath"
)

// Language enumeration.
const (
	JavaScript = iota
	Go
)

// Detect tries to detect the programming language of the repository.
// It's currently pretty simple and checks for the existence of various files.
func Detect(path string) (int, error) {
	if _, err := os.Stat(filepath.Join(path, "package.json")); err == nil {
		return JavaScript, nil
	}
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		return Go, nil
	}
	return -1, errors.New("cannot detect language")
}
