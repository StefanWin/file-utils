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

	var portraitDir string
	flag.StringVar(&portraitDir, "portraitDir", "./portrait", "Directory for portrait images.")

	var wideScreenDir string
	flag.StringVar(&wideScreenDir, "wideScreenDir", "./widescreen", "Directory for widescreen images.")

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

	files, err := pkg.ListImagesInDirectory(".")
	if err != nil {
		log.Fatal(err)
	}

	count := len(files)
	logInfo("found %d image files\n", count)

	if err := pkg.EnsureDir(portraitDir); err != nil {
		log.Fatal(err)
	}

	logInfo("created directorty %s\n", portraitDir)

	if err := pkg.EnsureDir(wideScreenDir); err != nil {
		log.Fatal(err)
	}

	logInfo("created directorty %s\n", wideScreenDir)

	wideScreenCount := 0
	portraitCount := 0
	aspectRationSum := 0.0

	for _, f := range files {
		src := filepath.Join("./", f.Name())
		dst := ""
		w, h, err := pkg.GetImageBounds(src)
		if err != nil {
			log.Fatal(err)
		}
		aspectRationSum += float64(w / h)
		if w >= h {
			dst = filepath.Join("./", wideScreenDir, f.Name())
			wideScreenCount += 1
		} else {
			dst = filepath.Join("./", portraitDir, f.Name())
			portraitCount += 1
		}
		if err := pkg.CopyFile(src, dst); err != nil {
			log.Fatal(err)
		}
		logVerbose("copied :: (%dx%d) %s => %s\n", w, h, src, dst)
	}

	avgAspectRatio := aspectRationSum / float64(count)
	logInfo("average aspect ratio : %f\n", avgAspectRatio)

	elapsed := time.Since(start)
	logInfo("successfully filtered %d (%d widescreen, %d portrait) in %s\n", count, wideScreenCount, portraitCount, elapsed)

}
