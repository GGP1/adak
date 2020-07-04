package cache

// This is a custom cache inspired on
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

// Cache holds a default expiration, a map of items,
// the mutexes and its name
type Cache struct {
	name              string
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
}

// Item represents the object that the user stores
// in the cache, with its expiration
type Item struct {
	Object     interface{}
	Expiration int64
}

// NewCache creates a cache
// TODO: Create a file that contains the cache map
// and set 2 functions: load and save
func NewCache(name string, defaultExpiration time.Duration) *Cache {
	items := make(map[string]Item)

	if defaultExpiration == 0 {
		defaultExpiration = -1
	}

	newCache := &Cache{
		name:              name,
		defaultExpiration: defaultExpiration,
		items:             items,
	}

	return newCache
}

// Add checks if an item doesn't exist yet, if not, stores it
func (c *Cache) Add(name, key string, obj interface{}, defaultExpiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.Get(name, key)
	if found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s already exists", key)
	}
	if name != c.name {
		return fmt.Errorf("Cache %s does not exist", name)
	}

	c.Set(name, key, obj, defaultExpiration)

	return nil
}

// Delete an item
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Get an item
func (c *Cache) Get(name, key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if name != c.name {
		return nil, false
	}

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

// ItemCount returns the number of items stored in the cache
func (c *Cache) ItemCount(name string) (int, error) {
	if name != c.name {
		return 0, fmt.Errorf("Cache %s does not exist", name)
	}

	c.mu.RLock()
	count := len(c.items)
	c.mu.RUnlock()

	return count, nil
}

// Items returns a map copy of the items stored in the cache
func (c *Cache) Items(name string) (map[string]Item, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if name != c.name {
		return nil, fmt.Errorf("Cache %s does not exist", name)
	}

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
	return newMap, nil
}

// Replace an item with a new one and store it
func (c *Cache) Replace(name, key string, obj interface{}, defaultExpiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if name != c.name {
		return fmt.Errorf("Cache %s does not exist", name)
	}

	_, found := c.Get(name, key)
	if !found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s doesn't exist", key)
	}
	c.Set(name, key, obj, defaultExpiration)

	return nil
}

// Reset deletes all the items from the cache
func (c *Cache) Reset(name string) error {
	if name != c.name {
		return fmt.Errorf("Cache %s does not exist", name)
	}

	c.mu.Lock()
	c.items = map[string]Item{}
	c.mu.Unlock()

	return nil
}

// Set a new item to the cache
func (c *Cache) Set(name, key string, obj interface{}, defaultExpiration time.Duration) error {
	var expiration int64

	if name != c.name {
		return fmt.Errorf("Cache %s does not exist", name)
	}

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

	return nil
}
