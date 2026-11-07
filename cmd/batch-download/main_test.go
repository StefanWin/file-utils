package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestParseURLs(t *testing.T) {
	got := parseURLs([]byte("https://a/1.jpg\r\n\r\nhttps://a/2.jpg\n  \nhttps://a/3.jpg\n"))
	want := []string{"https://a/1.jpg", "https://a/2.jpg", "https://a/3.jpg"}
	if !slices.Equal(got, want) {
		t.Fatalf("parseURLs = %v, want %v", got, want)
	}
	if got := parseURLs(nil); len(got) != 0 {
		t.Fatalf("parseURLs(nil) = %v", got)
	}
}

func TestDownloadAll(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing.jpg" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "content of %s", r.URL.Path)
	}))
	defer srv.Close()

	out := t.TempDir()
	urls := []string{srv.URL + "/a.jpg", srv.URL + "/missing.jpg", srv.URL + "/b.jpg"}
	var done []string
	errs := downloadAll(urls, out, 0, func(url, path string) { done = append(done, filepath.Base(path)) })

	if len(errs) != 1 {
		t.Fatalf("got %d errors (%v), want 1", len(errs), errs)
	}
	slices.Sort(done)
	if want := []string{"a.jpg", "b.jpg"}; !slices.Equal(done, want) {
		t.Fatalf("done = %v, want %v", done, want)
	}
	if b, _ := os.ReadFile(filepath.Join(out, "a.jpg")); string(b) != "content of /a.jpg" {
		t.Fatalf("a.jpg = %q", b)
	}
	if _, err := os.Stat(filepath.Join(out, "missing.jpg")); !os.IsNotExist(err) {
		t.Fatalf("missing.jpg should not be written, stat err = %v", err)
	}
}
