package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/StefanWin/file-utils/pkg"
)

const convSuffix = "_conv"

// outputName returns the name of the converted file for a video.
// ok is false for files that are themselves outputs of a previous run.
func outputName(name string) (out string, ok bool) {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	if strings.HasSuffix(base, convSuffix) {
		return "", false
	}
	return base + convSuffix + ".mp4", true
}

// ffmpegArgs returns the ffmpeg arguments to convert in to out with Nvidia's h264_nvenc encoder.
func ffmpegArgs(in, out string) []string {
	return []string{
		"-hwaccel", "cuda",
		"-hwaccel_output_format", "cuda",
		"-i", in,
		"-c:v", "h264_nvenc",
		"-c:a", "copy",
		out,
		"-y",
	}
}

func main() {
	start := time.Now()

	var quiet bool
	flag.BoolVar(&quiet, "quiet", false, "Disable output")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Enable per file logging output.")

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

	files, err := pkg.ListVideosInDirectory(".")
	if err != nil {
		log.Fatal(err)
	}

	logInfo("found %d video files\n", len(files))

	converted := 0
	for _, f := range files {
		out, ok := outputName(f.Name())
		if !ok {
			logVerbose("skipping already converted %s\n", f.Name())
			continue
		}
		cmd := exec.Command("ffmpeg", ffmpegArgs(f.Name(), out)...)
		logVerbose("running %s\n", strings.Join(cmd.Args, " "))
		if verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		converted++
	}

	elapsed := time.Since(start)

	logInfo("successfully converted %d videos in %s\n", converted, elapsed)
}
