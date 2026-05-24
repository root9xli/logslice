package index

import (
	"fmt"
	"time"
)

// Stats holds summary statistics about an index.
type Stats struct {
	// EntryCount is the total number of indexed log entries.
	EntryCount int

	// FirstTime is the timestamp of the earliest log entry.
	FirstTime time.Time

	// LastTime is the timestamp of the most recent log entry.
	LastTime time.Time

	// SpanDuration is the total time span covered by the index.
	SpanDuration time.Duration

	// UniqueOffsets is the number of distinct byte offsets in the index.
	UniqueOffsets int
}

// Compute derives Stats from the provided index entries.
// It assumes entries are sorted in ascending time order.
// Returns an empty Stats and an error if entries is nil or empty.
func Compute(entries []Entry) (Stats, error) {
	if len(entries) == 0 {
		return Stats{}, fmt.Errorf("index: cannot compute stats from empty entry list")
	}

	seen := make(map[int64]struct{}, len(entries))
	for _, e := range entries {
		seen[e.Offset] = struct{}{}
	}

	first := entries[0].Time
	last := entries[len(entries)-1].Time

	return Stats{
		EntryCount:    len(entries),
		FirstTime:     first,
		LastTime:      last,
		SpanDuration:  last.Sub(first),
		UniqueOffsets: len(seen),
	}, nil
}

// String returns a human-readable summary of the Stats.
func (s Stats) String() string {
	return fmt.Sprintf(
		"entries=%d unique_offsets=%d first=%s last=%s span=%s",
		s.EntryCount,
		s.UniqueOffsets,
		s.FirstTime.Format(time.RFC3339),
		s.LastTime.Format(time.RFC3339),
		s.SpanDuration.Round(time.Second),
	)
}
