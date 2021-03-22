package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/StefanWin/file-util/v2/pkg"
)

func main() {
	urls := os.Args[1:]
	for _, url := range urls {
		filename := filepath.Base(url)
		data, err := pkg.DownloadFile(url)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to download : %s", url))
		}
		if err := os.WriteFile(filename, data, 0644); err != nil {
			log.Fatal(fmt.Errorf("failed to write %s", filename))
		}
	}
}
