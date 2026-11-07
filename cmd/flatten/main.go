package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	var images bool
	flag.BoolVar(&images, "images", true, "Only use images (default unless -videos is set)")

	var videos bool
	flag.BoolVar(&videos, "videos", false, "Only use videos")

	flag.Parse()

	setFlags := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { setFlags[f.Name] = true })
	// -videos on its own implies -images=false
	if videos && !setFlags["images"] {
		images = false
	}

	if images && videos {
		log.Fatal(fmt.Errorf("can not use both 'images' and 'videos'"))
	}

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

	targetDirectory := "./all"
	if images {
		targetDirectory = "./images"
	}
	if videos {
		targetDirectory = "./videos"
	}

	logInfo("target directory => %s\n", targetDirectory)

	if err := pkg.EnsureDir(targetDirectory); err != nil {
		log.Fatal(err)
	}

	logInfo("created directory => %s\n", targetDirectory)

	directories, err := pkg.ListDirectoriesInDirectory(".")
	if err != nil {
		log.Fatal(err)
	}

	logInfo("found %d subdirectories\n", len(directories))

	// keep only the media files the selected mode asks for
	keep := func(e os.DirEntry) bool {
		ext := pkg.GetDirEntryEXT(e)
		switch {
		case images:
			return pkg.IsImageEXT(ext)
		case videos:
			return pkg.IsVideoEXT(ext)
		default:
			return pkg.IsImageEXT(ext) || pkg.IsVideoEXT(ext)
		}
	}

	var paths []string

	for _, dir := range directories {
		subDirPath := filepath.Join("./", dir.Name())
		// avoid duplicates if target already exists
		if fmt.Sprintf("./%s", subDirPath) == targetDirectory {
			logInfo("skipping target directory since already exists")
			continue
		}
		files, err := pkg.ListFilesInDirectoryOptions(subDirPath, &pkg.FilterOptions{
			FilterDirectories: true,
			FilterFuncs:       []pkg.DirEntryFilterFunc{keep},
		})
		if err != nil {
			log.Fatal(err)
		}
		logInfo("found %d target files in %s\n", len(files), dir.Name())
		for _, f := range files {
			paths = append(paths, filepath.Join("./", dir.Name(), f.Name()))
		}
	}

	logInfo("found %d total target files\n", len(paths))

	for _, p := range paths {
		src := p
		dst := filepath.Join(targetDirectory, filepath.Base(src))
		if err := pkg.CopyFile(src, dst); err != nil {
			log.Fatal(err)
		}
		logVerbose("copied :: %s => %s\n", src, dst)
	}

	elapsed := time.Since(start)

	logInfo("successfully moved %d to %s in %s\n", len(paths), targetDirectory, elapsed)
}
