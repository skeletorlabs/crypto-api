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

func Get[T any](c *MemoryCache, key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok || time.Now().After(item.ExpiresAt) {
		var zero T
		return zero, false
	}

	val, ok := item.Value.(T)
	return val, ok
}

// Eg: cache.Set[models.PriceResponse](marketCache, key, data, ttl)
func Set[T any](c *MemoryCache, key string, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}
