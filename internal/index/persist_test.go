package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveLoad_RoundTrip(t *testing.T) {
	entries := []Entry{
		{Offset: 0, Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Offset: 512, Time: time.Date(2024, 1, 1, 0, 1, 0, 0, time.UTC)},
		{Offset: 1024, Time: time.Date(2024, 1, 1, 0, 2, 0, 0, time.UTC)},
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "test.json")

	if err := Save(entries, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("saved file not found: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(got) != len(entries) {
		t.Fatalf("entry count: want %d, got %d", len(entries), len(got))
	}
	for i, e := range entries {
		if got[i].Offset != e.Offset {
			t.Errorf("[%d] offset: want %d, got %d", i, e.Offset, got[i].Offset)
		}
		if !got[i].Time.Equal(e.Time) {
			t.Errorf("[%d] time: want %v, got %v", i, e.Time, got[i].Time)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Fatal("expected error loading missing file, got nil")
	}
}

func TestSave_CreatesIntermediateDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "index.json")

	if err := Save([]Entry{}, path); err != nil {
		t.Fatalf("Save with nested dirs: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}
