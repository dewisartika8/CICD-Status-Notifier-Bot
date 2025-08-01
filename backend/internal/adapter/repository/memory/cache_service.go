package memory

import (
	"sync"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/port"
)

// item represents a cached item with expiration
type item struct {
	value      interface{}
	expiration int64
}

// expired checks if the item has expired
func (i item) expired() bool {
	return i.expiration > 0 && time.Now().UnixNano() > i.expiration
}

// InMemoryCache is a simple in-memory cache implementation
type InMemoryCache struct {
	items sync.Map
	mu    sync.RWMutex
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache() port.CacheServiceInterface {
	cache := &InMemoryCache{}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a value from cache
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	obj, found := c.items.Load(key)
	if !found {
		return nil, false
	}

	item := obj.(item)
	if item.expired() {
		c.items.Delete(key)
		return nil, false
	}

	return item.value, true
}

// Set stores a value in cache with TTL
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	c.items.Store(key, item{
		value:      value,
		expiration: expiration,
	})
}

// Delete removes a value from cache
func (c *InMemoryCache) Delete(key string) {
	c.items.Delete(key)
}

// cleanup removes expired items periodically
func (c *InMemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.items.Range(func(key, value interface{}) bool {
				item := value.(item)
				if item.expired() {
					c.items.Delete(key)
				}
				return true
			})
		}
	}
}
