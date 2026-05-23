package index_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/index"
	"github.com/user/logslice/internal/timeparse"
)

const sampleLog = `2024-01-01T10:00:00Z INFO starting up
2024-01-01T10:01:00Z INFO connected
2024-01-01T10:02:00Z WARN high memory
2024-01-01T10:03:00Z ERROR crash
2024-01-01T10:04:00Z INFO recovered
`

func newBuilder(t *testing.T) *index.Builder {
	t.Helper()
	p, err := timeparse.NewParser("")
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	return index.NewBuilder(p)
}

func TestBuild_EntryCount(t *testing.T) {
	b := newBuilder(t)
	r := strings.NewReader(sampleLog)
	idx, err := b.Build(r)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if len(idx) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(idx))
	}
}

func TestBuild_OffsetsAscending(t *testing.T) {
	b := newBuilder(t)
	r := strings.NewReader(sampleLog)
	idx, err := b.Build(r)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	for i := 1; i < len(idx); i++ {
		if idx[i].Offset <= idx[i-1].Offset {
			t.Errorf("offsets not ascending at %d: %d <= %d", i, idx[i].Offset, idx[i-1].Offset)
		}
	}
}

func TestSearchStart(t *testing.T) {
	b := newBuilder(t)
	r := strings.NewReader(sampleLog)
	idx, _ := b.Build(r)

	target := time.Date(2024, 1, 1, 10, 2, 0, 0, time.UTC)
	pos := idx.SearchStart(target)
	if pos != 2 {
		t.Errorf("SearchStart: expected 2, got %d", pos)
	}
}

func TestSearchEnd(t *testing.T) {
	b := newBuilder(t)
	r := strings.NewReader(sampleLog)
	idx, _ := b.Build(r)

	target := time.Date(2024, 1, 1, 10, 2, 0, 0, time.UTC)
	pos := idx.SearchEnd(target)
	if pos != 3 {
		t.Errorf("SearchEnd: expected 3, got %d", pos)
	}
}

func TestSearchStart_BeforeAll(t *testing.T) {
	b := newBuilder(t)
	r := strings.NewReader(sampleLog)
	idx, _ := b.Build(r)

	target := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if pos := idx.SearchStart(target); pos != 0 {
		t.Errorf("expected 0, got %d", pos)
	}
}

func TestSearchEnd_AfterAll(t *testing.T) {
	b := newBuilder(t)
	r := strings.NewReader(sampleLog)
	idx, _ := b.Build(r)

	target := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if pos := idx.SearchEnd(target); pos != len(idx) {
		t.Errorf("expected %d, got %d", len(idx), pos)
	}
}
