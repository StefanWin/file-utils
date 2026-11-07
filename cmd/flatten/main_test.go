package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestParseMode(t *testing.T) {
	cases := []struct {
		name                      string
		images, videos, imagesSet bool
		want                      mode
		wantErr                   bool
	}{
		{"default", true, false, false, modeImages, false},
		{"-videos", true, true, false, modeVideos, false},
		{"-images=false", false, false, true, modeAll, false},
		{"-images -videos", true, true, true, 0, true},
		{"-images=false -videos", false, true, true, modeVideos, false},
	}
	for _, c := range cases {
		got, err := parseMode(c.images, c.videos, c.imagesSet)
		if (err != nil) != c.wantErr || (!c.wantErr && got != c.want) {
			t.Errorf("%s: parseMode = %v, %v; want %v, err=%v", c.name, got, err, c.want, c.wantErr)
		}
	}
}

func TestCollect(t *testing.T) {
	root := t.TempDir()
	for _, p := range []string{
		"a/x.png", "a/v.mp4", "a/notes.txt",
		"b/y.JPG", "b/w.MKV",
		"a/nested/deep.png",  // only depth 1 is flattened
		"images/already.png", // target dir is skipped
		"top.png",            // files in root are ignored
	} {
		full := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	cases := map[mode][]string{
		modeImages: {"a/x.png", "b/y.JPG"},
		modeVideos: {"a/v.mp4", "b/w.MKV"},
		// "all" target is "all", so images/ is a normal source dir here
		modeAll: {"a/v.mp4", "a/x.png", "b/w.MKV", "b/y.JPG", "images/already.png"},
	}
	for m, want := range cases {
		got, err := collect(root, m)
		if err != nil {
			t.Fatal(err)
		}
		for i := range got {
			rel, _ := filepath.Rel(root, got[i])
			got[i] = filepath.ToSlash(rel)
		}
		slices.Sort(got)
		if !slices.Equal(got, want) {
			t.Errorf("collect(%s) = %v, want %v", m.targetDir(), got, want)
		}
	}
}
