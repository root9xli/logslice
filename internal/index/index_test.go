package index

import (
	"strings"
	"testing"

	"github.com/user/logslice/internal/timeparse"
)

func newBuilder() *Builder {
	p := timeparse.NewParser("")
	return NewBuilder(p)
}

func TestBuild_EntryCount(t *testing.T) {
	input := strings.NewReader(
		"2024-01-01T00:00:00Z line one\n" +
			"2024-01-01T01:00:00Z line two\n" +
			"no timestamp here\n" +
			"2024-01-01T02:00:00Z line three\n",
	)
	b := newBuilder()
	idx, err := b.Build(input)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(idx.Entries); got != 3 {
		t.Fatalf("expected 3 entries, got %d", got)
	}
}

func TestBuild_OffsetsAscending(t *testing.T) {
	input := strings.NewReader(
		"2024-01-01T00:00:00Z alpha\n" +
			"2024-01-01T01:00:00Z beta\n" +
			"2024-01-01T02:00:00Z gamma\n",
	)
	b := newBuilder()
	idx, err := b.Build(input)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(idx.Entries); i++ {
		if idx.Entries[i].Offset <= idx.Entries[i-1].Offset {
			t.Fatalf("offsets not strictly ascending at index %d", i)
		}
	}
}

func TestSearchStart(t *testing.T) {
	b := newBuilder()
	idx, _ := b.Build(strings.NewReader(
		"2024-01-01T00:00:00Z a\n" +
			"2024-01-01T06:00:00Z b\n" +
			"2024-01-01T12:00:00Z c\n",
	))
	if off := StartOffset(idx, Range{Start: 1, End: 3}); off != idx.Entries[1].Offset {
		t.Fatalf("unexpected start offset %d", off)
	}
}

func TestSearchEnd(t *testing.T) {
	b := newBuilder()
	idx, _ := b.Build(strings.NewReader(
		"2024-01-01T00:00:00Z a\n" +
			"2024-01-01T06:00:00Z b\n" +
			"2024-01-01T12:00:00Z c\n",
	))
	if off := EndOffset(idx, Range{Start: 0, End: 2}); off != idx.Entries[2].Offset {
		t.Fatalf("unexpected end offset %d", off)
	}
}

func TestBuild_Reset(t *testing.T) {
	b := newBuilder()
	b.Build(strings.NewReader("2024-01-01T00:00:00Z first\n")) //nolint
	b.Reset()
	idx, _ := b.Build(strings.NewReader("2024-01-01T01:00:00Z second\n"))
	if got := len(idx.Entries); got != 1 {
		t.Fatalf("expected 1 entry after reset, got %d", got)
	}
}
