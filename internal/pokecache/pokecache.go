package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mut      sync.Mutex
	entries  map[string]cacheEntry
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	data      []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Get(key *string) ([]byte, bool) {
	if key == nil {
		return nil, false
	}

	c.mut.Lock()
	defer c.mut.Unlock()

	entry, ok := c.entries[*key]
	if !ok {
		return nil, false
	}

	return entry.data, true
}

func (c *Cache) Add(key *string, data []byte) error {

	if key == nil || *key == "" {
		return fmt.Errorf("key is nil")
	}

	c.mut.Lock()
	defer c.mut.Unlock()

	if len(data) == 0 || data == nil {
		return fmt.Errorf("data is nil or empty")
	}

	c.entries[*key] = cacheEntry{
		createdAt: time.Now().UTC(),
		data:      data,
	}
	return nil
}

// reap removes entries that are older than the interval
// that was set when the cache was created
func (c *Cache) reap() {
	c.mut.Lock()
	defer c.mut.Unlock()

	for key, entry := range c.entries {
		// createdAt is UTC time, and time.Since() returns
		// UTC time as well, so we don't need to convert to UTC
		if time.Since(entry.createdAt) > c.interval {
			delete(c.entries, key)
		}
	}
}

// reapLoop loops forever, reaping the cache of old entries
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval / 2)
	defer ticker.Stop()

	for range ticker.C {
		c.reap()
	}
}
