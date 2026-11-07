package main

import (
	"slices"
	"testing"
)

func TestOutputName(t *testing.T) {
	cases := []struct {
		in, out string
		ok      bool
	}{
		{"clip.mkv", "clip_conv.mp4", true},
		{"clip.mp4", "clip_conv.mp4", true},
		// only the trailing extension is replaced
		{"my.mkv.backup.mkv", "my.mkv.backup_conv.mp4", true},
		{"clip_conv.mp4", "", false},
		{"clip_conv.MKV", "", false},
	}
	for _, c := range cases {
		out, ok := outputName(c.in)
		if out != c.out || ok != c.ok {
			t.Errorf("outputName(%q) = %q, %v; want %q, %v", c.in, out, ok, c.out, c.ok)
		}
	}
}

func TestFFmpegArgs(t *testing.T) {
	args := ffmpegArgs("in.mkv", "in_conv.mp4")
	i := slices.Index(args, "-i")
	if i < 0 || args[i+1] != "in.mkv" {
		t.Fatalf("input not passed via -i: %v", args)
	}
	if !slices.Contains(args, "in_conv.mp4") || !slices.Contains(args, "h264_nvenc") {
		t.Fatalf("unexpected args: %v", args)
	}
}
