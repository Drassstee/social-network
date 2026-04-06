// Package utils provides common helper functions and utilities used across the application,
// including caching, database transaction management, and environment variable handling.
package utils

import (
	"sync"
	"time"
)

//--------------------------------------------------------------------------------------|

// CacheItem represents an item stored in the cache with an expiration time.
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

//--------------------------------------------------------------------------------------|

// Cache is a simple in-memory key-value store with TTL.
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

//--------------------------------------------------------------------------------------|

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

//--------------------------------------------------------------------------------------|

// Set adds an item to the cache with a specified TTL.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	}
}

//--------------------------------------------------------------------------------------|

// Get retrieves an item from the cache. Returns (value, found).
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Value, true
}

//--------------------------------------------------------------------------------------|

// Delete removes an item from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

//--------------------------------------------------------------------------------------|

// Clear removes all items from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}
