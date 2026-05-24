package index

import (
	"sort"
	"time"
)

// Entry represents a single indexed log line with its timestamp and byte offset.
type Entry struct {
	Timestamp time.Time
	Offset    int64
}

// Merge combines two sorted index entry slices into a single sorted slice.
// Duplicate offsets are removed, keeping the first occurrence.
func Merge(a, b []Entry) []Entry {
	if len(a) == 0 {
		return cloneEntries(b)
	}
	if len(b) == 0 {
		return cloneEntries(a)
	}

	merged := make([]Entry, 0, len(a)+len(b))
	merged = append(merged, a...)
	merged = append(merged, b...)

	sort.Slice(merged, func(i, j int) bool {
		if merged[i].Timestamp.Equal(merged[j].Timestamp) {
			return merged[i].Offset < merged[j].Offset
		}
		return merged[i].Timestamp.Before(merged[j].Timestamp)
	})

	return deduplicate(merged)
}

// deduplicate removes entries with duplicate offsets, keeping the first.
func deduplicate(entries []Entry) []Entry {
	if len(entries) == 0 {
		return entries
	}
	seen := make(map[int64]struct{}, len(entries))
	out := entries[:0]
	for _, e := range entries {
		if _, ok := seen[e.Offset]; ok {
			continue
		}
		seen[e.Offset] = struct{}{}
		out = append(out, e)
	}
	return out
}

func cloneEntries(src []Entry) []Entry {
	dup := make([]Entry, len(src))
	copy(dup, src)
	return dup
}
