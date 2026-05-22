package slicer_test

import (
	"strings"
	"testing"

	"github.com/user/logslice/internal/slicer"
)

func TestWrite_Lines(t *testing.T) {
	res := &slicer.Result{
		Lines: []string{"line one", "line two"},
		Count: 2,
	}

	var sb strings.Builder
	if err := slicer.Write(&sb, res, slicer.WriteOptions{}); err != nil {
		t.Fatalf("Write: %v", err)
	}

	out := sb.String()
	if !strings.Contains(out, "line one") || !strings.Contains(out, "line two") {
		t.Errorf("unexpected output: %q", out)
	}
	if strings.Contains(out, "matched") {
		t.Error("summary should not appear when ShowCount is false")
	}
}

func TestWrite_ShowCount(t *testing.T) {
	res := &slicer.Result{
		Lines: []string{"entry"},
		Count: 1,
	}

	var sb strings.Builder
	if err := slicer.Write(&sb, res, slicer.WriteOptions{ShowCount: true}); err != nil {
		t.Fatalf("Write: %v", err)
	}

	out := sb.String()
	if !strings.Contains(out, "1 line(s) matched") {
		t.Errorf("expected summary line, got: %q", out)
	}
}

func TestWrite_Empty(t *testing.T) {
	res := &slicer.Result{}

	var sb strings.Builder
	if err := slicer.Write(&sb, res, slicer.WriteOptions{ShowCount: true}); err != nil {
		t.Fatalf("Write: %v", err)
	}

	out := sb.String()
	if !strings.Contains(out, "0 line(s) matched") {
		t.Errorf("expected zero summary, got: %q", out)
	}
}
