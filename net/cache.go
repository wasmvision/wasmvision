package net

import (
	"errors"

	"github.com/orsinium-labs/wypes"
)

// Cache is a cache for Net.
type Cache struct {
	netCache map[wypes.UInt32]Net
}

// NewCache creates a new Cache.
func NewCache() *Cache {
	return &Cache{
		netCache: map[wypes.UInt32]Net{},
	}
}

// Get returns a Net from the cache.
func (c *Cache) Get(id wypes.UInt32) (Net, bool) {
	frame, ok := c.netCache[id]
	return frame, ok
}

// Set sets a Net in the cache.
func (c *Cache) Set(n Net) error {
	if n.ID == 0 {
		return errors.New("net ID is 0")
	}

	c.netCache[n.ID] = n
	return nil
}

// Delete deletes a Net from the cache.
func (c *Cache) Delete(id wypes.UInt32) {
	delete(c.netCache, id)
}

// Close closes all Net in the cache.
func (c *Cache) Close() {
	for _, n := range c.netCache {
		n.Close()
	}
}
