package net

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/orsinium-labs/wypes"
)

// Cache is a cache for Net.
type Cache struct {
	netCache map[wypes.UInt32]Net

	ModelsDir string
}

// NewCache creates a new Cache.
func NewCache() *Cache {
	return &Cache{
		netCache: map[wypes.UInt32]Net{},
	}
}

// Get returns a Net from the cache.
func (c *Cache) Get(id wypes.UInt32) (Net, bool) {
	net, ok := c.netCache[id]
	return net, ok
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

// Close closes all Net in the cache and deletes them.
func (c *Cache) Close() {
	for _, n := range c.netCache {
		n.Close()
		c.Delete(n.ID)
	}
}

// ModelFile gets the model file path name for the Net.
func (c *Cache) ModelFileName(model string) string {
	if km, ok := KnownModels[model]; ok {
		return filepath.Join(c.ModelsDir, km.Filename)
	}

	return filepath.Join(c.ModelsDir, model)
}

func (c *Cache) ModelExists(model string) bool {
	if _, err := os.Stat(model); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (c *Cache) ModelWellKnown(model string) bool {
	if _, ok := KnownModels[model]; ok {
		return true
	}

	return false
}

func (c *Cache) DownloadModel(model string) error {
	km, ok := KnownModels[model]
	if !ok {
		return errors.New("model not found")
	}

	fmt.Printf("Downloading model %s...\n", km.Alias)

	return DownloadModel(km, c.ModelsDir)
}
