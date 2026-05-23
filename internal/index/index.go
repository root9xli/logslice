// Package index provides byte-offset indexing for log files to enable
// fast binary-search based time-range lookups without full sequential scans.
package index

import (
	"bufio"
	"io"
	"time"

	"github.com/user/logslice/internal/timeparse"
)

// Entry holds a byte offset and the parsed timestamp of a log line.
type Entry struct {
	Offset int64
	Time   time.Time
}

// Index is a sorted slice of offset/time entries for a log file.
type Index []Entry

// Builder constructs an Index by scanning a ReadSeeker.
type Builder struct {
	parser *timeparse.Parser
}

// NewBuilder returns a Builder that uses the given time parser.
func NewBuilder(p *timeparse.Parser) *Builder {
	return &Builder{parser: p}
}

// Build scans r from the current position and returns an Index.
// r must be seekable; the caller is responsible for seeking to the
// desired start position before calling Build.
func (b *Builder) Build(r io.ReadSeeker) (Index, error) {
	var idx Index
	var offset int64

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if t, ok := b.parser.Extract(line); ok {
			idx = append(idx, Entry{Offset: offset, Time: t})
		}
		// +1 for the newline byte consumed by the scanner
		offset += int64(len(scanner.Bytes())) + 1
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return idx, nil
}

// SearchStart returns the index of the first entry whose time is >= t.
// Returns len(idx) if no such entry exists.
func (idx Index) SearchStart(t time.Time) int {
	lo, hi := 0, len(idx)
	for lo < hi {
		mid := (lo + hi) / 2
		if idx[mid].Time.Before(t) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// SearchEnd returns the index one past the last entry whose time is <= t.
// Returns 0 if no such entry exists.
func (idx Index) SearchEnd(t time.Time) int {
	lo, hi := 0, len(idx)
	for lo < hi {
		mid := (lo + hi) / 2
		if idx[mid].Time.After(t) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}
