// Package index provides facilities for building and searching a byte-offset
// index over timestamped log streams.
//
// # Overview
//
// A Builder scans a log stream line by line, extracts timestamps via a
// timeparse.Parser, and records each (offset, timestamp) pair as an Entry.
// The resulting Index can then be queried with Search to locate the contiguous
// Range of entries whose timestamps fall within a requested [from, to] window.
//
// StartOffset and EndOffset translate a Range back into byte offsets suitable
// for seeking within the original file, allowing the slicer to read only the
// relevant portion of a potentially large compressed log.
//
// # Caching
//
// The Cache type persists a serialised Index to disk keyed by file path,
// file size, and mtime, avoiding repeated full-file scans for unchanged logs.
package index
