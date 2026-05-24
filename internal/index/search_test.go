package index

import (
	"testing"
	"time"
)

func makeIndex(timestamps ...string) *Index {
	var entries []Entry
	var offset int64
	for _, ts := range timestamps {
		t, _ := time.Parse(time.RFC3339, ts)
		entries = append(entries, Entry{Offset: offset, Timestamp: t})
		offset += 50
	}
	return &Index{Entries: entries}
}

// mustParseTime parses an RFC3339 timestamp and panics on failure.
// Intended for use in tests only.
func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic("invalid RFC3339 timestamp in test: " + s)
	}
	return t
}

func TestSearch_FullRange(t *testing.T) {
	idx := makeIndex("2024-01-01T00:00:00Z", "2024-01-01T01:00:00Z", "2024-01-01T02:00:00Z")
	from := mustParseTime("2024-01-01T00:00:00Z")
	to := mustParseTime("2024-01-01T02:00:00Z")
	r := Search(idx, from, to)
	if r.Start != 0 || r.End != 3 {
		t.Fatalf("expected [0,3), got [%d,%d)", r.Start, r.End)
	}
}

func TestSearch_SubRange(t *testing.T) {
	idx := makeIndex("2024-01-01T00:00:00Z", "2024-01-01T01:00:00Z", "2024-01-01T02:00:00Z")
	from := mustParseTime("2024-01-01T01:00:00Z")
	to := mustParseTime("2024-01-01T01:30:00Z")
	r := Search(idx, from, to)
	if r.Start != 1 || r.End != 2 {
		t.Fatalf("expected [1,2), got [%d,%d)", r.Start, r.End)
	}
}

func TestSearch_NoMatch_BeforeAll(t *testing.T) {
	idx := makeIndex("2024-01-01T10:00:00Z", "2024-01-01T11:00:00Z")
	from := mustParseTime("2024-01-01T00:00:00Z")
	to := mustParseTime("2024-01-01T09:00:00Z")
	r := Search(idx, from, to)
	if !r.Empty() {
		t.Fatalf("expected empty range, got [%d,%d)", r.Start, r.End)
	}
}

func TestSearch_EmptyIndex(t *testing.T) {
	idx := &Index{}
	from := mustParseTime("2024-01-01T00:00:00Z")
	to := mustParseTime("2024-01-01T23:59:59Z")
	r := Search(idx, from, to)
	if !r.Empty() {
		t.Fatalf("expected empty range for empty index")
	}
}

func TestSearch_EndOffset_EOF(t *testing.T) {
	idx := makeIndex("2024-01-01T00:00:00Z", "2024-01-01T01:00:00Z")
	from := mustParseTime("2024-01-01T00:00:00Z")
	to := mustParseTime("2024-01-01T02:00:00Z")
	r := Search(idx, from, to)
	if off := EndOffset(idx, r); off != -1 {
		t.Fatalf("expected -1 (EOF), got %d", off)
	}
}
