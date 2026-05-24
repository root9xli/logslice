package index

import "time"

// Filter returns a subset of entries whose timestamps fall within
// [start, end] (both inclusive). A zero value for start or end means
// unbounded on that side.
func Filter(entries []Entry, start, end time.Time) []Entry {
	if len(entries) == 0 {
		return nil
	}

	var result []Entry
	for _, e := range entries {
		if inBound(e.Timestamp, start, end) {
			result = append(result, e)
		}
	}
	return result
}

// inBound reports whether t falls within [start, end].
// A zero start means no lower bound; a zero end means no upper bound.
func inBound(t, start, end time.Time) bool {
	if !start.IsZero() && t.Before(start) {
		return false
	}
	if !end.IsZero() && t.After(end) {
		return false
	}
	return true
}

// FilterByOffset returns entries whose byte offset is in the half-open
// interval [lo, hi). When hi is 0 the upper bound is ignored.
func FilterByOffset(entries []Entry, lo, hi int64) []Entry {
	var result []Entry
	for _, e := range entries {
		if e.Offset < lo {
			continue
		}
		if hi > 0 && e.Offset >= hi {
			continue
		}
		result = append(result, e)
	}
	return result
}
