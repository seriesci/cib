package language

import (
	"errors"
	"os"
)

// Language enumeration.
const (
	JavaScript = iota
)

// Detect tries to detect the programming language of the repository.
// It's currently pretty simple and checks for the existence of various files.
func Detect() (int, error) {
	if _, err := os.Stat("package.json"); err == nil {
		return JavaScript, nil
	}
	return -1, errors.New("cannot detect language")
}
