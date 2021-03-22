package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/StefanWin/file-util/v2/pkg"
)

func worker(id int, jobs <-chan string, results chan<- error, output string, logger func(format string, v ...interface{})) {
	for url := range jobs {
		filename := filepath.Base(url)
		path := filepath.Join(output, filename)
		data, err := pkg.DownloadFile(url)
		if err != nil {
			results <- fmt.Errorf("tread<%d> failed to download : %s", id, url)
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			results <- fmt.Errorf("thread<%d> failed to write %s", id, path)
		}
		logger("thread<%d> downloaded %s => %s\n", id, url, path)
		results <- nil
	}
}

func main() {

	start := time.Now()

	var quiet bool
	flag.BoolVar(&quiet, "quiet", false, "Disable logging output.")

	var filePath string
	flag.StringVar(&filePath, "file", "", "path to the file containing the urls. new line separated.")

	var outputDir string
	flag.StringVar(&outputDir, "output", "./", "path to the output directory")

	var threadCount int
	flag.IntVar(&threadCount, "threads", 3, "number of threads to use. default 3.")

	flag.Parse()

	logInfo := func(format string, v ...interface{}) {
		if !quiet {
			log.Printf(format, v...)
		}
	}

	logInfo("downloading to %s\n", outputDir)

	if err := pkg.EnsureDir(outputDir); err != nil {
		log.Fatal(err)
	}

	logInfo("created %s\n", outputDir)

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load file %s", filePath))
	}

	logInfo("reading urls from %s\n", filePath)

	str := strings.Replace(string(data), "\r\n", "\n", -1) // fuck windows
	urls := strings.Split(str, "\n")

	count := len(urls)

	logInfo("found %d urls in %s\n", count, filePath)

	jobs := make(chan string, count)
	results := make(chan error, count)

	for i := 1; i <= threadCount; i++ {
		go worker(i, jobs, results, outputDir, logInfo)
	}

	logInfo("created %d threads\n", threadCount)

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	for j := 1; j <= count; j++ {
		if err := <-results; err != nil {
			log.Println(err)
		}
	}
	elapsed := time.Since(start)

	logInfo("downloaded %d files in %s\n", count, elapsed)

}
