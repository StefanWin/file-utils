package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"

	"github.com/StefanWin/file-util/v2/pkg"
)

func main() {

	start := time.Now()

	var quiet bool
	flag.BoolVar(&quiet, "quiet", false, "Disable output")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Enable per file logging output.")

	var directory string
	flag.StringVar(&directory, "dir", "./images", "The directory to move the images to. Will create if not exists.")
	flag.Parse()

	logInfo := func(format string, v ...interface{}) {
		if !quiet {
			log.Printf(format, v...)
		}
	}

	logVerbose := func(format string, v ...interface{}) {
		if !quiet && verbose {
			log.Printf(format, v...)
		}
	}

	logInfo("target directory => %s\n", directory)

	if err := pkg.EnsureDir(directory); err != nil {
		log.Fatal(err)
	}

	logInfo("created directory => %s\n", directory)

	files, err := pkg.ListImagesInDirectory(".")
	if err != nil {
		log.Fatal(err)
	}

	count := len(files)

	logInfo("found %d image files\n", count)

	for _, f := range files {
		src := filepath.Join("./", f.Name())
		dst := filepath.Join(directory, f.Name())
		if err := pkg.MoveFile(src, dst); err != nil {
			log.Fatal(err)
		}
		logVerbose("moved :: %s => %s\n", src, dst)
	}

	elapsed := time.Since(start)

	logInfo("successfully moved %d images to %s in %s\n", count, directory, elapsed)

}
