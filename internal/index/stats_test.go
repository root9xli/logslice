package index

import (
	"testing"
	"time"
)

func makeStatsEntries() []Entry {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return []Entry{
		{Time: base, Offset: 0},
		{Time: base.Add(5 * time.Minute), Offset: 512},
		{Time: base.Add(10 * time.Minute), Offset: 1024},
		{Time: base.Add(15 * time.Minute), Offset: 2048},
		{Time: base.Add(20 * time.Minute), Offset: 2048}, // duplicate offset
	}
}

func TestCompute_Basic(t *testing.T) {
	entries := makeStatsEntries()
	s, err := Compute(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.EntryCount != 5 {
		t.Errorf("EntryCount: want 5, got %d", s.EntryCount)
	}
	if s.UniqueOffsets != 4 {
		t.Errorf("UniqueOffsets: want 4, got %d", s.UniqueOffsets)
	}
}

func TestCompute_TimeRange(t *testing.T) {
	entries := makeStatsEntries()
	s, err := Compute(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantSpan := 20 * time.Minute
	if s.SpanDuration != wantSpan {
		t.Errorf("SpanDuration: want %s, got %s", wantSpan, s.SpanDuration)
	}

	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !s.FirstTime.Equal(base) {
		t.Errorf("FirstTime: want %s, got %s", base, s.FirstTime)
	}
	if !s.LastTime.Equal(base.Add(20 * time.Minute)) {
		t.Errorf("LastTime: want %s, got %s", base.Add(20*time.Minute), s.LastTime)
	}
}

func TestCompute_Empty(t *testing.T) {
	_, err := Compute(nil)
	if err == nil {
		t.Fatal("expected error for nil entries, got nil")
	}
	_, err = Compute([]Entry{})
	if err == nil {
		t.Fatal("expected error for empty entries, got nil")
	}
}

func TestStats_String(t *testing.T) {
	entries := makeStatsEntries()
	s, err := Compute(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := s.String()
	if len(out) == 0 {
		t.Error("String() returned empty output")
	}
	for _, want := range []string{"entries=5", "unique_offsets=4", "span="} {
		if !containsSubstr(out, want) {
			t.Errorf("String() missing %q in output: %s", want, out)
		}
	}
}

func containsSubstr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
