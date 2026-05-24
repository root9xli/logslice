package index

import (
	"testing"
	"time"
)

func makeSamplerEntries(n int) []Entry {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := make([]Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = Entry{
			Offset:    int64(i * 100),
			Timestamp: base.Add(time.Duration(i) * time.Minute),
		}
	}
	return entries
}

func TestNewSampler_StepOne(t *testing.T) {
	entries := makeSamplerEntries(10)
	s := NewSampler(entries, 1)
	if s.Len() != 10 {
		t.Fatalf("expected 10 samples, got %d", s.Len())
	}
}

func TestNewSampler_StepThree(t *testing.T) {
	entries := makeSamplerEntries(10)
	s := NewSampler(entries, 3)
	// indices 0,3,6,9 => 4 samples
	if s.Len() != 4 {
		t.Fatalf("expected 4 samples, got %d", s.Len())
	}
}

func TestNewSampler_ZeroStep(t *testing.T) {
	entries := makeSamplerEntries(5)
	s := NewSampler(entries, 0)
	if s.Step() != 1 {
		t.Fatalf("expected step=1 for zero input, got %d", s.Step())
	}
	if s.Len() != 5 {
		t.Fatalf("expected 5 samples, got %d", s.Len())
	}
}

func TestNewSampler_Empty(t *testing.T) {
	s := NewSampler([]Entry{}, 2)
	if s.Len() != 0 {
		t.Fatalf("expected 0 samples, got %d", s.Len())
	}
	if got := s.NearestOffset(time.Now()); got != -1 {
		t.Fatalf("expected -1 for empty sampler, got %d", got)
	}
}

func TestNearestOffset_ExactMatch(t *testing.T) {
	entries := makeSamplerEntries(5)
	s := NewSampler(entries, 1)
	target := entries[3].Timestamp
	got := s.NearestOffset(target)
	if got != entries[3].Offset {
		t.Fatalf("expected offset %d, got %d", entries[3].Offset, got)
	}
}

func TestNearestOffset_BeforeAll(t *testing.T) {
	entries := makeSamplerEntries(5)
	s := NewSampler(entries, 1)
	before := entries[0].Timestamp.Add(-time.Hour)
	got := s.NearestOffset(before)
	if got != entries[0].Offset {
		t.Fatalf("expected first offset %d, got %d", entries[0].Offset, got)
	}
}

func TestNearestOffset_AfterAll(t *testing.T) {
	entries := makeSamplerEntries(5)
	s := NewSampler(entries, 1)
	after := entries[4].Timestamp.Add(time.Hour)
	got := s.NearestOffset(after)
	if got != entries[4].Offset {
		t.Fatalf("expected last offset %d, got %d", entries[4].Offset, got)
	}
}

func TestSamples_ReturnsCopy(t *testing.T) {
	entries := makeSamplerEntries(3)
	s := NewSampler(entries, 1)
	copy1 := s.Samples()
	copy1[0].Offset = 9999
	copy2 := s.Samples()
	if copy2[0].Offset == 9999 {
		t.Fatal("Samples should return an independent copy")
	}
}
