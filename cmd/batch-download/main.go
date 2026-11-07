package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/StefanWin/file-utils/pkg"
)

// parseURLs splits the url file into urls: one per line, blank lines and \r\n handled.
func parseURLs(data []byte) []string {
	return strings.Fields(string(data))
}

// downloadAll downloads every url into outputDir using the given number of workers.
// onDone is called after each successful download. It returns one error per failed url.
func downloadAll(urls []string, outputDir string, workers int, onDone func(url, path string)) []error {
	workers = max(workers, 1)
	jobs := make(chan string)
	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)
	for range workers {
		wg.Go(func() {
			for url := range jobs {
				path := filepath.Join(outputDir, filepath.Base(url))
				err := download(url, path)
				mu.Lock()
				if err != nil {
					errs = append(errs, err)
				} else if onDone != nil {
					onDone(url, path)
				}
				mu.Unlock()
			}
		})
	}
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)
	wg.Wait()
	return errs
}

func download(url, path string) error {
	data, err := pkg.DownloadFile(url)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", url, err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}

func main() {

	start := time.Now()

	var quiet bool
	flag.BoolVar(&quiet, "quiet", false, "Disable logging output.")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Enable per file logging output.")

	var filePath string
	flag.StringVar(&filePath, "file", "./urls", "path to the file containing the urls. new line separated.")

	var outputDir string
	flag.StringVar(&outputDir, "output", "./", "path to the output directory")

	var threadCount int
	flag.IntVar(&threadCount, "threads", 3, "number of threads to use.")

	flag.Parse()

	logInfo := func(format string, v ...any) {
		if !quiet {
			log.Printf(format, v...)
		}
	}

	logVerbose := func(format string, v ...any) {
		if !quiet && verbose {
			log.Printf(format, v...)
		}
	}

	logInfo("downloading to %s\n", outputDir)

	if err := pkg.EnsureDir(outputDir); err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load file %s: %w", filePath, err))
	}

	urls := parseURLs(data)
	count := len(urls)

	logInfo("found %d urls in %s\n", count, filePath)

	errs := downloadAll(urls, outputDir, threadCount, func(url, path string) {
		logVerbose("downloaded %s => %s\n", url, path)
	})
	for _, err := range errs {
		log.Println(err)
	}

	elapsed := time.Since(start)

	logInfo("downloaded %d of %d files in %s\n", count-len(errs), count, elapsed)
	if len(errs) > 0 {
		os.Exit(1)
	}
}
