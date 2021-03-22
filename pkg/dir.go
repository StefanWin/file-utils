package pkg

import (
	"os"
	"path/filepath"
)

// A filter function for os.DirEntry.
type DirEntryFilterFunc = func(os.DirEntry) bool

// Options to filter certain entries from files in a directory.
// Filters are applied in order (dir -> image -> video -> additional)
// TODO: add recursive option
type FilterOptions struct {
	// Filter all directories.
	FilterDirectories bool
	// Filter all images.
	FilterImages bool
	// Filter all videos.
	FilterVideos bool
	// Customizable filter functions.
	FilterFuncs []DirEntryFilterFunc
}

// GetDirEntryEXT returns the file extension of the given DirEntry.
// Will return empty string if it's a directory (i.e. no '.' in the filename).
func GetDirEntryEXT(dirEntry os.DirEntry) string {
	return filepath.Ext(dirEntry.Name())
}

// filterDirEntries filters the given slice of DirEntries with given filter function.
func filterDirEntries(entries []os.DirEntry, test DirEntryFilterFunc) []os.DirEntry {
	result := make([]os.DirEntry, 0)
	for _, entry := range entries {
		if test(entry) {
			result = append(result, entry)
		}
	}
	return result
}

// A function for filtering directories for use in filterDirEntries.
func filterDirFunc(dirEntry os.DirEntry) bool {
	return !dirEntry.IsDir()
}

// Inverse of filterDirFunc
func onlyDirFunc(dirEntry os.DirEntry) bool {
	return !filterDirFunc(dirEntry)
}

// A function for filtering images for use in filterDirEntries.
func filterImagesFunc(dirEntry os.DirEntry) bool {
	return !IsImageEXT(GetDirEntryEXT(dirEntry))
}

// Inverse for filterImagesFunc.
func onlyImagesFunc(dirEntry os.DirEntry) bool {
	return !filterImagesFunc(dirEntry)
}

// A function for filtering videos for use in filterDirEntries.
func filterVideosFunc(dirEntry os.DirEntry) bool {
	return !IsVideoEXT(GetDirEntryEXT(dirEntry))
}

// Inverse for filterVideosFunc.
func onlyVideosFunc(dirEntry os.DirEntry) bool {
	return !filterVideosFunc(dirEntry)
}

// ListFilesInDirectory lists all files in the given directory
// and applies the given filter options.
func ListFilesInDirectoryOptions(directory string, options *FilterOptions) ([]os.DirEntry, error) {
	if options == nil {
		return os.ReadDir(directory)
	}
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	// TODO: probably better to use intersections for filter results
	if options.FilterDirectories {
		files = filterDirEntries(files, filterDirFunc)
	}
	if options.FilterImages {
		files = filterDirEntries(files, filterImagesFunc)
	}
	if options.FilterVideos {
		files = filterDirEntries(files, filterVideosFunc)
	}
	// len(arr) where arr == nil returns 0
	if len(options.FilterFuncs) > 0 {
		for _, f := range options.FilterFuncs {
			files = filterDirEntries(files, f)
		}
	}
	return files, nil
}

// ListFilesInDirectory lists all files/directories in the given directory.
func ListFilesInDirectory(directory string) ([]os.DirEntry, error) {
	return ListFilesInDirectoryOptions(directory, nil)
}

// ListDirectories lists all sub-directories for the given directory.
func ListDirectoriesInDirectory(directory string) ([]os.DirEntry, error) {
	return ListFilesInDirectoryOptions(directory, &FilterOptions{
		FilterFuncs: []DirEntryFilterFunc{onlyDirFunc},
	})
}

// ListImagesInDirectory lists all images in the given directory.
func ListImagesInDirectory(directory string) ([]os.DirEntry, error) {
	return ListFilesInDirectoryOptions(directory, &FilterOptions{
		FilterDirectories: true,
		FilterFuncs:       []DirEntryFilterFunc{onlyImagesFunc},
	})
}

// ListVideosInDirectory lists all videos in the given directory.
func ListVideosInDirectory(directory string) ([]os.DirEntry, error) {
	return ListFilesInDirectoryOptions(directory, &FilterOptions{
		FilterDirectories: true,
		FilterFuncs:       []DirEntryFilterFunc{onlyVideosFunc},
	})
}
