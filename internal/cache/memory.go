package cache

import (
	"sync"
	"time"
)

type Item struct {
	Value     any
	ExpiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]Item),
	}
}

func (c *MemoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Value, true
}

func (c *MemoryCache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	c.items[key] = Item{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()
}
