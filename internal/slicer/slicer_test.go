package slicer_test

import (
	"strings"
	"testing"

	"github.com/user/logslice/internal/slicer"
	"github.com/user/logslice/internal/timeparse"
)

func newParser(t *testing.T) *timeparse.Parser {
	t.Helper()
	p, err := timeparse.NewParser("")
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	return p
}

const sampleLog = `2024-01-10T10:00:00Z INFO  startup complete
2024-01-10T10:05:00Z DEBUG request received
2024-01-10T10:10:00Z INFO  processing done
2024-01-10T10:15:00Z WARN  high latency
2024-01-10T10:20:00Z ERROR timeout
`

func TestSlice_BasicRange(t *testing.T) {
	p := newParser(t)
	r := strings.NewReader(sampleLog)
	res, err := slicer.Slice(r, slicer.Options{
		Parser: p,
		From:   "2024-01-10T10:05:00Z",
		To:     "2024-01-10T10:15:00Z",
	})
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if res.Count != 3 {
		t.Errorf("expected 3 lines, got %d", res.Count)
	}
}

func TestSlice_MaxLines(t *testing.T) {
	p := newParser(t)
	r := strings.NewReader(sampleLog)
	res, err := slicer.Slice(r, slicer.Options{
		Parser:   p,
		From:     "2024-01-10T10:00:00Z",
		To:       "2024-01-10T10:20:00Z",
		MaxLines: 2,
	})
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if res.Count != 2 {
		t.Errorf("expected 2 lines, got %d", res.Count)
	}
}

func TestSlice_NoMatch(t *testing.T) {
	p := newParser(t)
	r := strings.NewReader(sampleLog)
	res, err := slicer.Slice(r, slicer.Options{
		Parser: p,
		From:   "2024-01-11T00:00:00Z",
		To:     "2024-01-11T01:00:00Z",
	})
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if res.Count != 0 {
		t.Errorf("expected 0 lines, got %d", res.Count)
	}
}

func TestSlice_EmptyInput(t *testing.T) {
	p := newParser(t)
	res, err := slicer.Slice(strings.NewReader(""), slicer.Options{
		Parser: p,
		From:   "2024-01-10T10:00:00Z",
		To:     "2024-01-10T11:00:00Z",
	})
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if res.Count != 0 {
		t.Errorf("expected 0 lines, got %d", res.Count)
	}
}
