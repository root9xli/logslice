package reader

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func writePlain(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "plain-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func writeGzip(t *testing.T, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "compressed.log.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create gz file: %v", err)
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	for _, l := range lines {
		gw.Write([]byte(l + "\n"))
	}
	gw.Close()
	return path
}

var sampleLines = []string{
	"2024-01-15T10:00:00Z INFO  server started",
	"2024-01-15T10:00:05Z DEBUG request received",
	"2024-01-15T10:00:10Z ERROR something failed",
}

func collectLines(t *testing.T, path string) []string {
	t.Helper()
	lr, err := New(path)
	if err != nil {
		t.Fatalf("New(%q): %v", path, err)
	}
	defer lr.Close()

	var got []string
	for lr.Scan() {
		got = append(got, lr.Text())
	}
	if err := lr.Err(); err != nil {
		t.Fatalf("scan error: %v", err)
	}
	return got
}

func TestNew_PlainFile(t *testing.T) {
	path := writePlain(t, sampleLines)
	got := collectLines(t, path)
	if len(got) != len(sampleLines) {
		t.Fatalf("expected %d lines, got %d", len(sampleLines), len(got))
	}
	for i, want := range sampleLines {
		if got[i] != want {
			t.Errorf("line %d: want %q, got %q", i, want, got[i])
		}
	}
}

func TestNew_GzipFile(t *testing.T) {
	path := writeGzip(t, sampleLines)
	got := collectLines(t, path)
	if len(got) != len(sampleLines) {
		t.Fatalf("expected %d lines, got %d", len(sampleLines), len(got))
	}
	for i, want := range sampleLines {
		if got[i] != want {
			t.Errorf("line %d: want %q, got %q", i, want, got[i])
		}
	}
}

func TestNew_NotFound(t *testing.T) {
	_, err := New("/nonexistent/path/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
