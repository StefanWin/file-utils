package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"

	"github.com/StefanWin/file-utils/pkg"
)

type stats struct {
	wide, portrait int
	avgAspectRatio float64
}

// sortImages copies every image in srcDir into wideDir (width >= height) or portraitDir.
// onCopy is called for each copied image. The target directories are only created if there are images.
func sortImages(srcDir, wideDir, portraitDir string, onCopy func(w, h int, src, dst string)) (stats, error) {
	var s stats
	files, err := pkg.ListImagesInDirectory(srcDir)
	if err != nil || len(files) == 0 {
		return s, err
	}
	for _, dir := range []string{portraitDir, wideDir} {
		if err := pkg.EnsureDir(dir); err != nil {
			return s, err
		}
	}

	aspectRatioSum := 0.0
	for _, f := range files {
		src := filepath.Join(srcDir, f.Name())
		w, h, err := pkg.GetImageBounds(src)
		if err != nil {
			return s, err
		}
		aspectRatioSum += float64(w) / float64(h)
		dst := filepath.Join(portraitDir, f.Name())
		if w >= h {
			dst = filepath.Join(wideDir, f.Name())
			s.wide++
		} else {
			s.portrait++
		}
		if err := pkg.CopyFile(src, dst); err != nil {
			return s, err
		}
		if onCopy != nil {
			onCopy(w, h, src, dst)
		}
	}
	s.avgAspectRatio = aspectRatioSum / float64(len(files))
	return s, nil
}

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

	s, err := sortImages(".", wideScreenDir, portraitDir, func(w, h int, src, dst string) {
		logVerbose("copied :: (%dx%d) %s => %s\n", w, h, src, dst)
	})
	if err != nil {
		log.Fatal(err)
	}

	count := s.wide + s.portrait
	if count == 0 {
		logInfo("found no image files\n")
		return
	}

	logInfo("average aspect ratio : %f\n", s.avgAspectRatio)

	elapsed := time.Since(start)
	logInfo("successfully filtered %d (%d widescreen, %d portrait) in %s\n", count, s.wide, s.portrait, elapsed)
}
