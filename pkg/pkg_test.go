package pkg

import (
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestIsImageAndVideoEXT(t *testing.T) {
	cases := []struct {
		ext          string
		image, video bool
	}{
		{".jpg", true, false},
		{".JPG", true, false},
		{".Png", true, false},
		{".mp4", false, true},
		{".MKV", false, true},
		{".txt", false, false},
		{"", false, false},
		{"jpg", false, false},
	}
	for _, c := range cases {
		if got := IsImageEXT(c.ext); got != c.image {
			t.Errorf("IsImageEXT(%q) = %v, want %v", c.ext, got, c.image)
		}
		if got := IsVideoEXT(c.ext); got != c.video {
			t.Errorf("IsVideoEXT(%q) = %v, want %v", c.ext, got, c.video)
		}
	}
}

func TestEnsureDirIsWritable(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b")
	if err := EnsureDir(dir); err != nil {
		t.Fatal(err)
	}
	// calling again on an existing dir is fine
	if err := EnsureDir(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "f"), []byte("x"), 0o644); err != nil {
		t.Fatalf("created directory is not writable: %v", err)
	}
}

func TestCopyFileReportsCreateError(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := CopyFile(src, filepath.Join(tmp, "missing", "dst")); err == nil {
		t.Fatal("expected error when destination directory does not exist")
	}
	dst := filepath.Join(tmp, "dst")
	if err := CopyFile(src, dst); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(dst); string(b) != "hello" {
		t.Fatalf("dst = %q, want %q", b, "hello")
	}
}

func TestDownloadFileRejectsErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ok" {
			w.Write([]byte("data"))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	if b, err := DownloadFile(srv.URL + "/ok"); err != nil || string(b) != "data" {
		t.Fatalf("DownloadFile(/ok) = %q, %v", b, err)
	}
	if _, err := DownloadFile(srv.URL + "/missing"); err == nil {
		t.Fatal("expected error for 404")
	}
}

func TestGetImageBounds(t *testing.T) {
	path := filepath.Join(t.TempDir(), "img.png")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, 30, 20))); err != nil {
		t.Fatal(err)
	}
	f.Close()

	w, h, err := GetImageBounds(path)
	if err != nil || w != 30 || h != 20 {
		t.Fatalf("GetImageBounds = %d, %d, %v; want 30, 20, nil", w, h, err)
	}
	// file handle must be released (on Windows, removing an open file fails)
	if err := os.Remove(path); err != nil {
		t.Fatalf("file still open: %v", err)
	}
}
