// Package index builds and caches byte-offset indexes for log files.
//
// # Overview
//
// Large log files are expensive to scan sequentially when only a small
// time-range slice is needed. This package solves the problem by building
// an in-memory Index — a sorted list of (byte-offset, timestamp) pairs —
// that allows binary-search lookups for any time range.
//
// # Usage
//
//	p, _ := timeparse.NewParser("")
//	b := index.NewBuilder(p)
//	idx, err := b.Build(readSeeker)
//
//	start := idx.SearchStart(from)
//	end   := idx.SearchEnd(to)
//	// seek to idx[start].Offset and read until idx[end].Offset
//
// # Caching
//
// Use Cache to avoid rebuilding indexes on repeated queries against the
// same file. The cache validates entries against file size and mtime so
// stale indexes are never returned after a log rotation or append.
package index
