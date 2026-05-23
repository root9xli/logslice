package index

import (
	"testing"
	"time"
)

func makeIndex(times []time.Time) *Index {
	entries := make([]Entry, len(times))
	for i, t := range times {
		entries[i] = Entry{Time: t, Offset: int64(i * 100)}
	}
	return &Index{Entries: entries}
}

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestSearch_FullRange(t *testing.T) {
	times := []time.Time{
		base,
		base.Add(1 * time.Minute),
		base.Add(2 * time.Minute),
		base.Add(3 * time.Minute),
	}
	idx := makeIndex(times)

	result := Search(idx, base, base.Add(3*time.Minute))
	if result.NoMatch() {
		t.Fatal("expected a match")
	}
	if result.StartOffset != 0 {
		t.Errorf("StartOffset = %d, want 0", result.StartOffset)
	}
	if result.EndOffset != -1 {
		t.Errorf("EndOffset = %d, want -1 (EOF)", result.EndOffset)
	}
}

func TestSearch_SubRange(t *testing.T) {
	times := []time.Time{
		base,
		base.Add(1 * time.Minute),
		base.Add(2 * time.Minute),
		base.Add(3 * time.Minute),
	}
	idx := makeIndex(times)

	result := Search(idx, base.Add(1*time.Minute), base.Add(2*time.Minute))
	if result.NoMatch() {
		t.Fatal("expected a match")
	}
	if result.StartOffset != 100 {
		t.Errorf("StartOffset = %d, want 100", result.StartOffset)
	}
	if result.EndOffset != 300 {
		t.Errorf("EndOffset = %d, want 300", result.EndOffset)
	}
}

func TestSearch_NoMatch_BeforeAll(t *testing.T) {
	times := []time.Time{
		base.Add(5 * time.Minute),
		base.Add(6 * time.Minute),
	}
	idx := makeIndex(times)

	result := Search(idx, base, base.Add(1*time.Minute))
	if !result.NoMatch() {
		t.Errorf("expected NoMatch, got %+v", result)
	}
}

func TestSearch_EmptyIndex(t *testing.T) {
	idx := &Index{Entries: []Entry{}}
	result := Search(idx, base, base.Add(time.Hour))
	if result.StartOffset != 0 || result.EndOffset != -1 {
		t.Errorf("unexpected result for empty index: %+v", result)
	}
}

func TestSearch_NoMatch_AfterAll(t *testing.T) {
	times := []time.Time{
		base,
		base.Add(1 * time.Minute),
	}
	idx := makeIndex(times)

	result := Search(idx, base.Add(10*time.Minute), base.Add(20*time.Minute))
	if !result.NoMatch() {
		t.Errorf("expected NoMatch, got %+v", result)
	}
}
