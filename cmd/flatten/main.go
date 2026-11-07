package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/StefanWin/file-utils/pkg"
)

type mode int

const (
	modeAll mode = iota
	modeImages
	modeVideos
)

// parseMode resolves the -images/-videos flags.
// -videos on its own implies -images=false, since -images defaults to true.
func parseMode(images, videos, imagesSet bool) (mode, error) {
	if videos && !imagesSet {
		images = false
	}
	switch {
	case images && videos:
		return 0, errors.New("can not use both 'images' and 'videos'")
	case images:
		return modeImages, nil
	case videos:
		return modeVideos, nil
	default:
		return modeAll, nil
	}
}

// targetDir is the name of the directory the files are copied to.
func (m mode) targetDir() string {
	switch m {
	case modeImages:
		return "images"
	case modeVideos:
		return "videos"
	default:
		return "all"
	}
}

// keep reports whether the entry is a media file this mode copies.
func (m mode) keep(e os.DirEntry) bool {
	ext := pkg.GetDirEntryEXT(e)
	switch m {
	case modeImages:
		return pkg.IsImageEXT(ext)
	case modeVideos:
		return pkg.IsVideoEXT(ext)
	default:
		return pkg.IsImageEXT(ext) || pkg.IsVideoEXT(ext)
	}
}

// collect returns the paths of all matching files in the direct subdirectories of root,
// skipping the target directory itself.
func collect(root string, m mode) ([]string, error) {
	directories, err := pkg.ListDirectoriesInDirectory(root)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, dir := range directories {
		if dir.Name() == m.targetDir() {
			continue
		}
		subDirPath := filepath.Join(root, dir.Name())
		files, err := pkg.ListFilesInDirectoryOptions(subDirPath, &pkg.FilterOptions{
			FilterDirectories: true,
			FilterFuncs:       []pkg.DirEntryFilterFunc{m.keep},
		})
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			paths = append(paths, filepath.Join(subDirPath, f.Name()))
		}
	}
	return paths, nil
}

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

	imagesSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "images" {
			imagesSet = true
		}
	})

	m, err := parseMode(images, videos, imagesSet)
	if err != nil {
		log.Fatal(err)
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

	targetDirectory := filepath.Join(".", m.targetDir())

	logInfo("target directory => %s\n", targetDirectory)

	if err := pkg.EnsureDir(targetDirectory); err != nil {
		log.Fatal(err)
	}

	paths, err := collect(".", m)
	if err != nil {
		log.Fatal(err)
	}

	logInfo("found %d total target files\n", len(paths))

	for _, src := range paths {
		dst := filepath.Join(targetDirectory, filepath.Base(src))
		if err := pkg.CopyFile(src, dst); err != nil {
			log.Fatal(err)
		}
		logVerbose("copied :: %s => %s\n", src, dst)
	}

	elapsed := time.Since(start)

	logInfo("successfully copied %d files to %s in %s\n", len(paths), targetDirectory, elapsed)
}
