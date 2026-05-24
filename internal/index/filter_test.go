package index

import (
	"testing"
	"time"
)

func makeFilterEntries() []Entry {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return []Entry{
		{Timestamp: base.Add(0 * time.Hour), Offset: 0},
		{Timestamp: base.Add(1 * time.Hour), Offset: 100},
		{Timestamp: base.Add(2 * time.Hour), Offset: 200},
		{Timestamp: base.Add(3 * time.Hour), Offset: 300},
		{Timestamp: base.Add(4 * time.Hour), Offset: 400},
	}
}

func TestFilter_FullRange(t *testing.T) {
	entries := makeFilterEntries()
	got := Filter(entries, time.Time{}, time.Time{})
	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
}

func TestFilter_SubRange(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := makeFilterEntries()
	start := base.Add(1 * time.Hour)
	end := base.Add(3 * time.Hour)
	got := Filter(entries, start, end)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	if got[0].Offset != 100 || got[2].Offset != 300 {
		t.Errorf("unexpected offsets: %v", got)
	}
}

func TestFilter_NoMatch(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := makeFilterEntries()
	start := base.Add(10 * time.Hour)
	end := base.Add(12 * time.Hour)
	got := Filter(entries, start, end)
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestFilter_Empty(t *testing.T) {
	got := Filter(nil, time.Time{}, time.Time{})
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestFilterByOffset_Basic(t *testing.T) {
	entries := makeFilterEntries()
	got := FilterByOffset(entries, 100, 300)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Offset != 100 || got[1].Offset != 200 {
		t.Errorf("unexpected offsets: %v", got)
	}
}

func TestFilterByOffset_NoUpperBound(t *testing.T) {
	entries := makeFilterEntries()
	got := FilterByOffset(entries, 200, 0)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}
