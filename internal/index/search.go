package index

import (
	"sort"
	"time"
)

// Range describes a contiguous span of entries within an Index.
type Range struct {
	Start int // inclusive index into Index.Entries
	End   int // exclusive index into Index.Entries
}

// Empty reports whether the range contains no entries.
func (r Range) Empty() bool { return r.Start >= r.End }

// Search returns the Range of entries whose timestamps fall within [from, to].
// Both endpoints are inclusive. If no entries match, the returned Range is empty.
func Search(idx *Index, from, to time.Time) Range {
	if len(idx.Entries) == 0 {
		return Range{}
	}

	// Find first entry with Timestamp >= from.
	start := sort.Search(len(idx.Entries), func(i int) bool {
		return !idx.Entries[i].Timestamp.Before(from)
	})

	// Find first entry with Timestamp > to.
	end := sort.Search(len(idx.Entries), func(i int) bool {
		return idx.Entries[i].Timestamp.After(to)
	})

	return Range{Start: start, End: end}
}

// StartOffset returns the byte offset of the first entry in r, or 0 if empty.
func StartOffset(idx *Index, r Range) int64 {
	if r.Empty() {
		return 0
	}
	return idx.Entries[r.Start].Offset
}

// EndOffset returns the byte offset just past the last entry in r.
// Returns -1 to signal "read to EOF" when the range extends to the last entry.
func EndOffset(idx *Index, r Range) int64 {
	if r.Empty() {
		return 0
	}
	if r.End >= len(idx.Entries) {
		return -1
	}
	return idx.Entries[r.End].Offset
}
