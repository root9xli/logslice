package index

import (
	"bufio"
	"io"
	"time"

	"github.com/user/logslice/internal/timeparse"
)

// Entry represents a single indexed log line with its byte offset and timestamp.
type Entry struct {
	Offset    int64
	Timestamp time.Time
}

// Index holds an ordered slice of entries built from a log stream.
type Index struct {
	Entries []Entry
}

// Builder incrementally constructs an Index by scanning a log stream.
type Builder struct {
	parser  *timeparse.Parser
	entries []Entry
}

// NewBuilder creates a Builder using the provided time parser.
func NewBuilder(p *timeparse.Parser) *Builder {
	return &Builder{parser: p}
}

// Build reads all lines from r, recording byte offsets and parsed timestamps.
func (b *Builder) Build(r io.Reader) (*Index, error) {
	var offset int64
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		ts, ok := b.parser.Extract(line)
		if ok {
			b.entries = append(b.entries, Entry{
				Offset:    offset,
				Timestamp: ts,
			})
		}
		offset += int64(len(line)) + 1 // +1 for newline
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &Index{Entries: b.entries}, nil
}

// Reset clears accumulated entries so the Builder can be reused.
func (b *Builder) Reset() {
	b.entries = b.entries[:0]
}
