package frame

import (
	"errors"

	"github.com/orsinium-labs/wypes"
)

// Cache is a cache for frames and other data that needs to be shared between modules, or the guest and host.
type Cache struct {
	frameCache map[wypes.UInt32]Frame

	// ReturnDataPtr is a pointer to a linear memory buffer for return values.
	ReturnDataPtr uint32
}

// NewCache creates a new Cache.
func NewCache() *Cache {
	return &Cache{
		frameCache: map[wypes.UInt32]Frame{},
	}
}

// Get returns a frame from the cache.
func (c *Cache) Get(id wypes.UInt32) (Frame, bool) {
	frame, ok := c.frameCache[id]
	return frame, ok
}

// Set sets a frame in the cache.
func (c *Cache) Set(frame Frame) error {
	if frame.ID == 0 {
		return errors.New("frame ID is 0")
	}

	c.frameCache[frame.ID] = frame
	return nil
}

// Delete deletes a frame from the cache.
func (c *Cache) Delete(id wypes.UInt32) {
	delete(c.frameCache, id)
}

// Close closes all frames in the cache and also deletes them.
func (c *Cache) Close() {
	for _, frame := range c.frameCache {
		frame.Close()
		c.Delete(frame.ID)
	}
}
