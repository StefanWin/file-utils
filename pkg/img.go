package pkg

import (
	"image"
	_ "image/jpeg"
	"os"
)

// Get the image bounds from a path.
func GetImageBounds(path string) (int, int, error) {
	reader, err := os.Open(path)
	if err != nil {
		return -1, -1, err
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		return -1, -1, err
	}
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	return width, height, nil
}
