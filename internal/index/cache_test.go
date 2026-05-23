package index_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/index"
)

var (
	fakeSize  int64 = 1024
	fakeMtime       = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	fakeIdx         = index.Index{{Offset: 0, Time: fakeMtime}}
)

func TestCache_PutAndGet(t *testing.T) {
	c := index.NewCache(0)
	c.Put("a.log", fakeIdx, fakeSize, fakeMtime)

	got, ok := c.Get("a.log", fakeSize, fakeMtime)
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(got) != len(fakeIdx) {
		t.Errorf("expected %d entries, got %d", len(fakeIdx), len(got))
	}
}

func TestCache_MissingKey(t *testing.T) {
	c := index.NewCache(0)
	_, ok := c.Get("missing.log", fakeSize, fakeMtime)
	if ok {
		t.Fatal("expected cache miss for unknown key")
	}
}

func TestCache_SizeMismatch(t *testing.T) {
	c := index.NewCache(0)
	c.Put("b.log", fakeIdx, fakeSize, fakeMtime)
	_, ok := c.Get("b.log", fakeSize+1, fakeMtime)
	if ok {
		t.Fatal("expected cache miss on size change")
	}
}

func TestCache_MtimeMismatch(t *testing.T) {
	c := index.NewCache(0)
	c.Put("c.log", fakeIdx, fakeSize, fakeMtime)
	_, ok := c.Get("c.log", fakeSize, fakeMtime.Add(time.Second))
	if ok {
		t.Fatal("expected cache miss on mtime change")
	}
}

func TestCache_TTLExpiry(t *testing.T) {
	c := index.NewCache(50 * time.Millisecond)
	c.Put("d.log", fakeIdx, fakeSize, fakeMtime)

	time.Sleep(80 * time.Millisecond)
	_, ok := c.Get("d.log", fakeSize, fakeMtime)
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c := index.NewCache(0)
	c.Put("e.log", fakeIdx, fakeSize, fakeMtime)
	c.Invalidate("e.log")
	_, ok := c.Get("e.log", fakeSize, fakeMtime)
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestCache_Len(t *testing.T) {
	c := index.NewCache(0)
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
	c.Put("f.log", fakeIdx, fakeSize, fakeMtime)
	c.Put("g.log", fakeIdx, fakeSize, fakeMtime)
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}
