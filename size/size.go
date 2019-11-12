package size

import (
	"os"
	"path/filepath"
)

// Directory returns size in kilobytes for given directory path.
func Directory(path string) (float64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return float64(size) / 1024.0, err
}
