package cache

import (
	"testing"
	"time"
)

func TestMemoryCache_Generics(t *testing.T) {
	c := NewMemoryCache()

	// This test ensures that we can store and retrieve different types without conflicts
	Set(c, "price", 69000.50, time.Minute)
	Set(c, "count", 100, time.Minute)

	valPrice, okPrice := Get[float64](c, "price")
	if !okPrice || valPrice != 69000.50 {
		t.Errorf("failed to get float64 from cache")
	}

	valCount, okCount := Get[int](c, "count")
	if !okCount || valCount != 100 {
		t.Errorf("failed to get int from cache")
	}

	// Testing failure of type assertion
	_, okWrong := Get[string](c, "count")
	if okWrong {
		t.Errorf("should have failed to get int as string")
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	c := NewMemoryCache()

	// This test ensures that items expire correctly after their TTL
	Set(c, "expired", "bye", time.Millisecond)
	time.Sleep(2 * time.Millisecond)

	_, ok := Get[string](c, "expired")
	if ok {
		t.Errorf("item should have expired")
	}
}
