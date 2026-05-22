package timeparse_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timeparse"
)

func TestExtract_RFC3339(t *testing.T) {
	p := timeparse.NewParser(nil, time.UTC)
	line := "2024-03-15T08:23:01Z some log message here"
	got, offset, err := p.Extract(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 15, 8, 23, 1, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("time mismatch: got %v, want %v", got, want)
	}
	if offset == 0 {
		t.Error("offset should be non-zero")
	}
}

func TestExtract_SpaceSeparated(t *testing.T) {
	p := timeparse.NewParser(nil, time.UTC)
	line := "2024-03-15 08:23:01 INFO starting server"
	got, _, err := p.Extract(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 15, 8, 23, 1, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("time mismatch: got %v, want %v", got, want)
	}
}

func TestExtract_NoTimestamp(t *testing.T) {
	p := timeparse.NewParser(nil, time.UTC)
	_, _, err := p.Extract("this line has no timestamp at all")
	if err != timeparse.ErrNoTimestamp {
		t.Errorf("expected ErrNoTimestamp, got %v", err)
	}
}

func TestExtract_CustomFormat(t *testing.T) {
	custom := []string{"2006/01/02 15:04:05"}
	p := timeparse.NewParser(custom, time.UTC)
	line := "2024/03/15 08:23:01 custom format log"
	got, _, err := p.Extract(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 15, 8, 23, 1, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("time mismatch: got %v, want %v", got, want)
	}
}

func TestInRange(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	cases := []struct {
		name string
		t    time.Time
		want bool
	}{
		{"before range", time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC), false},
		{"at start", from, true},
		{"in middle", time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC), true},
		{"at end", to, true},
		{"after range", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := timeparse.InRange(tc.t, from, to); got != tc.want {
				t.Errorf("InRange(%v) = %v, want %v", tc.t, got, tc.want)
			}
		})
	}
}
