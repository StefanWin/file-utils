package pkg

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// EnsureDir creates a directory (and any missing parents) if it does not exist.
func EnsureDir(directory string) error {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", directory, err)
	}
	return nil
}

// MoveFile moves the src file to the dst.
func MoveFile(src, dst string) error {
	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("failed to move file %s -> %s: %w", src, dst, err)
	}
	return nil
}

// CopyFile copies the src file to the dst.
func CopyFile(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := dstFile.Close(); err == nil {
			err = cerr
		}
	}()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// DownloadFile downloads the given url.
// Non-2xx responses are returned as errors.
func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}
