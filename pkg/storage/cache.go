package storage

// This is a simplified cache inspired on
// github.com/patrickmn/go-cache

import (
	"fmt"
	"sync"
	"time"
)

// Expiration ON/OFF
const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

// Item represents the object that the user stores
// in the cache, with its expiration
type Item struct {
	Object     interface{}
	Expiration int64
}

// Cache holds a default expiration, a map of items
// and the mutexes
type Cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
}

// NewCache creates a cache
func NewCache(defaultExpiration time.Duration) *Cache {
	items := make(map[string]Item)

	if defaultExpiration == 0 {
		defaultExpiration = -1
	}

	newCache := &Cache{
		defaultExpiration: defaultExpiration,
		items:             items,
	}

	return newCache
}

// Add checks if an item doesn't exist yet, if not, stores it
func (c *Cache) Add(key string, obj interface{}, defaultExpiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.Get(key)
	if found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s already exists", key)
	}
	c.Set(key, obj, defaultExpiration)

	return nil
}

// Delete an item
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Get an item
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}

	return item.Object, true
}

// Items returns a map copy of the items stored in the cache
func (c *Cache) Items() map[string]Item {
	c.mu.Lock()
	defer c.mu.Unlock()

	newMap := make(map[string]Item, len(c.items))
	now := time.Now().UnixNano()

	for k, v := range c.items {
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		newMap[k] = v
	}
	return newMap
}

// Replace an item with a new one and store it
func (c *Cache) Replace(key string, obj interface{}, defaultExpiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.Get(key)
	if !found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s doesn't exist", key)
	}
	c.Set(key, obj, defaultExpiration)

	return nil
}

// Set a new item to the cache
func (c *Cache) Set(key string, obj interface{}, defaultExpiration time.Duration) {
	var expiration int64

	if defaultExpiration == DefaultExpiration {
		defaultExpiration = c.defaultExpiration
	}

	if defaultExpiration > 0 {
		expiration = time.Now().Add(defaultExpiration).UnixNano()
	}

	c.mu.Lock()

	c.items[key] = Item{
		Object:     obj,
		Expiration: expiration,
	}

	c.mu.Unlock()
}
