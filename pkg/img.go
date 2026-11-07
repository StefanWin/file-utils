package pkg

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// GetImageBounds returns the width and height of the image at path.
// Only the image header is decoded, not the pixel data.
func GetImageBounds(path string) (int, int, error) {
	reader, err := os.Open(path)
	if err != nil {
		return -1, -1, err
	}
	defer reader.Close()
	cfg, _, err := image.DecodeConfig(reader)
	if err != nil {
		return -1, -1, err
	}
	return cfg.Width, cfg.Height, nil
}
