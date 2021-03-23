package pkg

import (
	"fmt"
	"io"
	"net/http"
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

// CopyFile copies the src file to the dst.
// adopted from https://opensource.com/article/18/6/copying-files-go
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// DownloadFile downloads the given url.
func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
