package pkg

import (
	"fmt"
	"os"
)

// check if the given value is in the given slice
func contains(arr []string, v string) bool {
	for _, s := range arr {
		if s == v {
			return true
		}
	}
	return false
}

// EnsureDir creates a directory if it does not exist.
func EnsureDir(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.Mkdir(directory, 0644)
		if err != nil {
			return fmt.Errorf("failed to create directory : %s", directory)
		}
	}
	return nil
}

// MoveFile moves the src file to the dst.
func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("failed to move file : %s -> %s", src, dst)
	}
	return nil
}
