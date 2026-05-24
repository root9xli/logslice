package index

import (
	"testing"
	"time"
)

func makeEntries(offsets []int64, base time.Time, step time.Duration) []Entry {
	entries := make([]Entry, len(offsets))
	for i, off := range offsets {
		entries[i] = Entry{
			Timestamp: base.Add(time.Duration(i) * step),
			Offset:    off,
		}
	}
	return entries
}

func TestMerge_DisjointSets(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a := makeEntries([]int64{0, 100, 200}, base, time.Second)
	b := makeEntries([]int64{300, 400, 500}, base.Add(10*time.Second), time.Second)

	result := Merge(a, b)
	if len(result) != 6 {
		t.Fatalf("expected 6 entries, got %d", len(result))
	}
	for i := 1; i < len(result); i++ {
		if result[i].Timestamp.Before(result[i-1].Timestamp) {
			t.Errorf("entries not sorted at index %d", i)
		}
	}
}

func TestMerge_DuplicateOffsets(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a := makeEntries([]int64{0, 100, 200}, base, time.Second)
	b := makeEntries([]int64{100, 200, 300}, base.Add(time.Second), time.Second)

	result := Merge(a, b)
	seen := make(map[int64]int)
	for _, e := range result {
		seen[e.Offset]++
		if seen[e.Offset] > 1 {
			t.Errorf("duplicate offset %d in merged result", e.Offset)
		}
	}
}

func TestMerge_EmptyA(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	b := makeEntries([]int64{10, 20}, base, time.Second)

	result := Merge(nil, b)
	if len(result) != len(b) {
		t.Fatalf("expected %d entries, got %d", len(b), len(result))
	}
}

func TestMerge_EmptyB(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a := makeEntries([]int64{10, 20}, base, time.Second)

	result := Merge(a, nil)
	if len(result) != len(a) {
		t.Fatalf("expected %d entries, got %d", len(a), len(result))
	}
}

func TestMerge_BothEmpty(t *testing.T) {
	result := Merge(nil, nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(result))
	}
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a := makeEntries([]int64{0, 100}, base, time.Second)
	b := makeEntries([]int64{50, 150}, base.Add(500*time.Millisecond), time.Second)

	aOrig := cloneEntries(a)
	bOrig := cloneEntries(b)
	Merge(a, b)

	for i := range a {
		if a[i] != aOrig[i] {
			t.Errorf("input a was mutated at index %d", i)
		}
	}
	for i := range b {
		if b[i] != bOrig[i] {
			t.Errorf("input b was mutated at index %d", i)
		}
	}
}
