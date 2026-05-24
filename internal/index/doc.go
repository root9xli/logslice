// Package index provides an offset-based index for large log files.
//
// The index maps timestamps to byte offsets within a (possibly compressed)
// log file, enabling fast binary-search seeks into the file without
// scanning every line.
//
// # Building
//
// Use [NewBuilder] to create an index from a channel of (offset, time) pairs
// produced while scanning a log file.
//
// # Searching
//
// [Search] returns the sub-slice of index entries that fall within a given
// time range. [StartOffset] and [EndOffset] extract the byte offsets that
// bracket the matching region.
//
// # Caching
//
// [NewCache] wraps an on-disk JSON cache keyed by file path, size, and mtime
// so that the index is rebuilt only when the source file changes.
package index
