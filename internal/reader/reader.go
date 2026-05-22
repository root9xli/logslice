package reader

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// LineReader wraps a file (plain or gzip-compressed) and exposes a
// line-by-line scanner so callers never need to care about compression.
type LineReader struct {
	file    *os.File
	gz      *gzip.Reader
	scanner *bufio.Scanner
}

// New opens the given path for reading. Files ending in ".gz" are
// transparently decompressed.
func New(path string) (*LineReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("reader: open %q: %w", path, err)
	}

	lr := &LineReader{file: f}

	var src io.Reader = f
	if strings.HasSuffix(path, ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("reader: gzip init %q: %w", path, err)
		}
		lr.gz = gr
		src = gr
	}

	lr.scanner = bufio.NewScanner(src)
	lr.scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return lr, nil
}

// Scan advances to the next line. Returns true while lines are available.
func (lr *LineReader) Scan() bool {
	return lr.scanner.Scan()
}

// Text returns the current line without the trailing newline.
func (lr *LineReader) Text() string {
	return lr.scanner.Text()
}

// Err returns the first non-EOF error encountered by the scanner.
func (lr *LineReader) Err() error {
	return lr.scanner.Err()
}

// Close releases all underlying resources.
func (lr *LineReader) Close() error {
	if lr.gz != nil {
		if err := lr.gz.Close(); err != nil {
			return err
		}
	}
	return lr.file.Close()
}
