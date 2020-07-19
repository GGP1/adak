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
	defaultExpiration time.Duration
	items             map[string]Item
	sync.RWMutex
}

// Item represents the object that the user stores
// in the cache, with its expiration
type Item struct {
	Object     interface{}
	Expiration int64
}

// New creates a cache
// TODO: Create a file that contains the cache map
// and set 2 functions: load and save
func New(defaultExpiration time.Duration) *Cache {
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
	c.Lock()
	defer c.Unlock()

	_, found := c.Get(key)
	if found {
		c.Unlock()
		return fmt.Errorf("item %s already exists", key)
	}

	c.Set(key, obj, defaultExpiration)

	return nil
}

// Delete an item
func (c *Cache) Delete(key string) {
	c.Lock()
	delete(c.items, key)
	c.Unlock()
}

// Get an item
func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		c.RUnlock()
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.RUnlock()
			return nil, false
		}
	}

	return item.Object, true
}

// Items returns a map copy of the items stored in the cache
func (c *Cache) Items() (map[string]Item, error) {
	c.Lock()
	defer c.Unlock()

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
func (c *Cache) Replace(key string, obj interface{}, defaultExpiration time.Duration) error {
	c.Lock()
	defer c.Unlock()

	_, found := c.Get(key)
	if !found {
		c.Unlock()
		return fmt.Errorf("item %s doesn't exist", key)
	}
	c.Set(key, obj, defaultExpiration)

	return nil
}

// Reset deletes all the items from the cache
func (c *Cache) Reset() error {
	c.Lock()
	c.items = map[string]Item{}
	c.Unlock()

	return nil
}

// Set a new item to the cache
func (c *Cache) Set(key string, obj interface{}, defaultExpiration time.Duration) error {
	var expiration int64

	if defaultExpiration == DefaultExpiration {
		defaultExpiration = c.defaultExpiration
	}

	if defaultExpiration > 0 {
		expiration = time.Now().Add(defaultExpiration).UnixNano()
	}

	c.Lock()

	c.items[key] = Item{
		Object:     obj,
		Expiration: expiration,
	}

	c.Unlock()

	return nil
}

// Size returns the number of items stored in the cache
func (c *Cache) Size() (int, error) {
	c.RLock()
	defer c.RUnlock()

	return len(c.items), nil
}
