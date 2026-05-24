package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Save writes the index entries to path as a JSON file, creating any
// intermediate directories as needed.
func Save(entries []Entry, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("index persist: mkdir %s: %w", filepath.Dir(path), err)
	}

	f, err := os.CreateTemp(filepath.Dir(path), ".logslice-idx-*")
	if err != nil {
		return fmt.Errorf("index persist: create temp: %w", err)
	}
	tmpName := f.Name()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "")
	if err := enc.Encode(entries); err != nil {
		f.Close()
		os.Remove(tmpName)
		return fmt.Errorf("index persist: encode: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("index persist: close temp: %w", err)
	}

	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("index persist: rename to %s: %w", path, err)
	}
	return nil
}

// Load reads index entries previously written by Save.
func Load(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("index persist: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("index persist: decode %s: %w", path, err)
	}
	return entries, nil
}
