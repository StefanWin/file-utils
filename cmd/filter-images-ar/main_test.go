package main

import (
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func writePNG(t *testing.T, path string, w, h int) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, w, h))); err != nil {
		t.Fatal(err)
	}
}

func TestSortImages(t *testing.T) {
	src := t.TempDir()
	writePNG(t, filepath.Join(src, "wide.png"), 300, 100)
	writePNG(t, filepath.Join(src, "tall.PNG"), 100, 300)
	writePNG(t, filepath.Join(src, "square.png"), 50, 50)
	if err := os.WriteFile(filepath.Join(src, "notes.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := t.TempDir()
	wide, portrait := filepath.Join(out, "wide"), filepath.Join(out, "portrait")
	s, err := sortImages(src, wide, portrait, nil)
	if err != nil {
		t.Fatal(err)
	}
	if s.wide != 2 || s.portrait != 1 {
		t.Fatalf("wide=%d portrait=%d, want 2 and 1", s.wide, s.portrait)
	}
	// (3 + 1/3 + 1) / 3
	if want := (3.0 + 1.0/3 + 1) / 3; math.Abs(s.avgAspectRatio-want) > 1e-9 {
		t.Fatalf("avgAspectRatio = %f, want %f", s.avgAspectRatio, want)
	}
	for _, p := range []string{"wide/wide.png", "wide/square.png", "portrait/tall.PNG"} {
		if _, err := os.Stat(filepath.Join(out, p)); err != nil {
			t.Errorf("expected %s: %v", p, err)
		}
	}
}

func TestSortImagesEmptyDir(t *testing.T) {
	out := t.TempDir()
	wide := filepath.Join(out, "wide")
	s, err := sortImages(t.TempDir(), wide, filepath.Join(out, "portrait"), nil)
	if err != nil || s.wide+s.portrait != 0 {
		t.Fatalf("sortImages = %+v, %v", s, err)
	}
	if _, err := os.Stat(wide); !os.IsNotExist(err) {
		t.Fatalf("target dir should not be created for an empty source")
	}
}
