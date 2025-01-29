package localfs

import (
	"os"
	"path/filepath"
)

func FileExists(dir, filename string) (bool, error) {
	// Construct the full path to the file
	fullPath := filepath.Join(dir, filename)

	// Use os.Stat to get the file information
	_, err := os.Stat(fullPath)

	// If the error is nil, the file exists
	if err == nil {
		return true, nil
	}

	// If the error is of type *os.PathError, the file does not exist
	if os.IsNotExist(err) {
		return false, nil
	}

	// For any other error, return false and the error
	return false, err
}
