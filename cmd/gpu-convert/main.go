package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/StefanWin/file-util/v2/pkg"
)

func main() {
	start := time.Now()

	var quiet bool
	flag.BoolVar(&quiet, "quiet", false, "Disable output")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Enable per file logging output.")

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

	files, err := pkg.ListVideosInDirectory(".")
	if err != nil {
		log.Fatal(err)
	}

	count := len(files)

	logInfo("found %d video files\n", count)

	for _, f := range files {
		cmd := exec.Command(
			"ffmpeg",
			"-hwaccel", "cuda",
			"-hwaccel_output_format", "cuda",
			"-i", f.Name(),
			"-c:v", "h264_nvenc",
			"-c:a", "copy",
			strings.Replace(f.Name(), filepath.Ext(f.Name()), "_conv.mp4", -1),
			"-y",
		)
		logVerbose("running %s\n", strings.Join(cmd.Args, " "))
		if verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

	elapsed := time.Since(start)

	logInfo("successfully converted %d videos in %s\n", count, elapsed)

}
