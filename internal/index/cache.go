package index

import (
	"sync"
	"time"
)

// CacheEntry wraps an Index with metadata used for invalidation.
type CacheEntry struct {
	idx       Index
	builtAt   time.Time
	fileSize  int64
	fileMtime time.Time
}

// Cache is a thread-safe, in-memory store for built indexes keyed by file path.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewCache returns a Cache that evicts entries older than ttl.
// Pass 0 to disable TTL-based eviction (only size/mtime invalidation applies).
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

// Get retrieves a cached Index for path if it is still valid given the
// current file size and modification time.
func (c *Cache) Get(path string, size int64, mtime time.Time) (Index, bool) {
	c.mu.RLock()
	e, ok := c.entries[path]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if e.fileSize != size || !e.fileMtime.Equal(mtime) {
		return nil, false
	}
	if c.ttl > 0 && time.Since(e.builtAt) > c.ttl {
		return nil, false
	}
	return e.idx, true
}

// Put stores idx in the cache for path with the given file metadata.
func (c *Cache) Put(path string, idx Index, size int64, mtime time.Time) {
	e := &CacheEntry{
		idx:       idx,
		builtAt:   time.Now(),
		fileSize:  size,
		fileMtime: mtime,
	}
	c.mu.Lock()
	c.entries[path] = e
	c.mu.Unlock()
}

// Invalidate removes the cached entry for path, if any.
func (c *Cache) Invalidate(path string) {
	c.mu.Lock()
	delete(c.entries, path)
	c.mu.Unlock()
}

// Len returns the number of entries currently in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	n := len(c.entries)
	c.mu.RUnlock()
	return n
}
