package cache

// This is a custom cache used to store orders
// inspired by github.com/patrickmn/go-cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/GGP1/palo/pkg/shopping/ordering"
)

// Expiration ON/OFF.
const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

// Cache holds a default expiration, a map of items,
// the mutexes and its name.
type Cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	sync.RWMutex
}

// Item represents the order that the user stores
// in the cache, with their expiration.
type Item struct {
	Order      ordering.Order
	Expiration int64
}

// New creates a cache.
func New(defaultExpiration time.Duration) *Cache {
	items := make(map[string]Item)

	if defaultExpiration == 0 {
		defaultExpiration = -1
	}

	cache := &Cache{
		defaultExpiration: defaultExpiration,
		items:             items,
	}

	return cache
}

// Add checks if an item doesn't exist yet, if not, stores it.
func (c *Cache) Add(key string, order ordering.Order, defaultExpiration time.Duration) error {
	c.Lock()
	defer c.Unlock()

	_, found := c.Get(key)
	if found {
		c.Unlock()
		return fmt.Errorf("item %s already exists", key)
	}

	c.Set(key, order, defaultExpiration)

	return nil
}

// Delete an item.
func (c *Cache) Delete(key string) {
	c.Lock()
	delete(c.items, key)
	c.Unlock()
}

// Get an item.
func (c *Cache) Get(key string) (ordering.Order, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		c.RUnlock()
		return ordering.Order{}, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.RUnlock()
			return ordering.Order{}, false
		}
	}

	return item.Order, true
}

// Items returns a map copy of the items stored in the cache.
func (c *Cache) Items() (map[string]ordering.Order, error) {
	c.Lock()
	defer c.Unlock()

	newMap := make(map[string]ordering.Order, len(c.items))
	now := time.Now().UnixNano()

	for k, v := range c.items {
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		newMap[k] = v.Order
	}

	return newMap, nil
}

// Replace an item with a new one and store it.
func (c *Cache) Replace(key string, order ordering.Order, defaultExpiration time.Duration) error {
	c.Lock()
	defer c.Unlock()

	_, found := c.Get(key)
	if !found {
		c.Unlock()
		return fmt.Errorf("item %s doesn't exist", key)
	}
	c.Set(key, order, defaultExpiration)

	return nil
}

// Reset deletes all the items from the cache.
func (c *Cache) Reset() error {
	c.Lock()
	c.items = map[string]Item{}
	c.Unlock()

	return nil
}

// Set a new item to the cache.
func (c *Cache) Set(key string, order ordering.Order, defaultExpiration time.Duration) error {
	var expiration int64

	if defaultExpiration == DefaultExpiration {
		defaultExpiration = c.defaultExpiration
	}

	if defaultExpiration > 0 {
		expiration = time.Now().Add(defaultExpiration).UnixNano()
	}

	c.Lock()

	c.items[key] = Item{
		Order:      order,
		Expiration: expiration,
	}

	c.Unlock()

	return nil
}

// Size returns the number of items stored in the cache.
func (c *Cache) Size() (int, error) {
	c.RLock()
	defer c.RUnlock()

	return len(c.items), nil
}
