package pkg

import (
	"slices"
	"strings"
)

var (
	videoExts = []string{".mp4", ".avi", ".webm", ".mkv", ".flv", ".wmv", ".ts"}
	imageExts = []string{".jpg", ".jpeg", ".png", ".bmp", ".gif"}
)

// IsImageEXT reports whether ext (e.g. ".jpg", ".JPG") is an image extension.
func IsImageEXT(ext string) bool {
	return slices.Contains(imageExts, strings.ToLower(ext))
}

// IsVideoEXT reports whether ext (e.g. ".mp4", ".MP4") is a video extension.
func IsVideoEXT(ext string) bool {
	return slices.Contains(videoExts, strings.ToLower(ext))
}
