package index

import (
	"sort"
	"time"
)

// SearchResult holds the byte offsets bounding a time range within a log file.
type SearchResult struct {
	// StartOffset is the byte offset of the first log line at or after the start time.
	StartOffset int64
	// EndOffset is the byte offset just past the last log line at or before the end time.
	// A value of -1 means read until EOF.
	EndOffset int64
}

// Search returns the byte range within the index that covers [start, end].
// It uses binary search over the index entries.
func Search(idx *Index, start, end time.Time) SearchResult {
	if len(idx.Entries) == 0 {
		return SearchResult{StartOffset: 0, EndOffset: -1}
	}

	// Find the first entry whose time >= start.
	startIdx := sort.Search(len(idx.Entries), func(i int) bool {
		return !idx.Entries[i].Time.Before(start)
	})

	// Find the last entry whose time <= end.
	endIdx := sort.Search(len(idx.Entries), func(i int) bool {
		return idx.Entries[i].Time.After(end)
	}) - 1

	if startIdx >= len(idx.Entries) || endIdx < 0 || startIdx > endIdx {
		return SearchResult{StartOffset: -1, EndOffset: -1}
	}

	startOffset := idx.Entries[startIdx].Offset

	var endOffset int64
	if endIdx+1 < len(idx.Entries) {
		endOffset = idx.Entries[endIdx+1].Offset
	} else {
		endOffset = -1 // read until EOF
	}

	return SearchResult{
		StartOffset: startOffset,
		EndOffset:   endOffset,
	}
}

// NoMatch returns true when the SearchResult indicates no lines were found.
func (r SearchResult) NoMatch() bool {
	return r.StartOffset == -1
}
