package frame

import "github.com/orsinium-labs/wypes"

type Cache struct {
	frameCache map[wypes.UInt32]Frame
}

func NewCache() *Cache {
	return &Cache{
		frameCache: map[wypes.UInt32]Frame{},
	}
}

func (c *Cache) Get(id wypes.UInt32) (Frame, bool) {
	frame, ok := c.frameCache[id]
	return frame, ok
}

func (c *Cache) Set(id wypes.UInt32, frame Frame) {
	c.frameCache[id] = frame
}

func (c *Cache) Delete(id wypes.UInt32) {
	delete(c.frameCache, id)
}

func (c *Cache) Close() {
	for _, frame := range c.frameCache {
		frame.Close()
	}
}
