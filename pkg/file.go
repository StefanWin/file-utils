package pkg

import (
	"regexp"
)

var (
	videoExts = []string{".mp4", ".avi", ".webm", ".mkv", ".flv", ".wmv"}
	imageExts = []string{".jpg", ".jpeg", ".png", ".bmp", ".gif"}

	extRegex = regexp.MustCompile(".[a-zA-Z0-9]+")
)

// Check if an extension is a image extension.
func IsImageEXT(ext string) bool {
	if extRegex.MatchString(ext) {
		return contains(imageExts, ext)
	}
	return false
}

// Check if an extension is a video extension.
func IsVideoEXT(ext string) bool {
	if extRegex.MatchString(ext) {
		return contains(videoExts, ext)
	}
	return false
}
