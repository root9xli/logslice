package main

import (
	"testing"
	"time"
)

func TestParseTime_RFC3339(t *testing.T) {
	t.Parallel()
	got, err := parseTime("2024-03-15T10:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseTime_SpaceSeparated(t *testing.T) {
	t.Parallel()
	got, err := parseTime("2024-03-15 10:00:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseTime_DateOnly(t *testing.T) {
	t.Parallel()
	got, err := parseTime("2024-03-15")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseTime_Invalid(t *testing.T) {
	t.Parallel()
	_, err := parseTime("not-a-date")
	if err == nil {
		t.Fatal("expected error for invalid timestamp, got nil")
	}
}

func TestParseTime_Formats(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input string
		want  time.Time
	}{
		{"2024-01-02T15:04:05Z", time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)},
		{"2024-01-02 15:04:05", time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)},
		{"2024-01-02", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			got, err := parseTime(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
