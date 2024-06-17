package cache

import (
	"sync"
	"sync/atomic"
)

type Cache struct {
	sync.Mutex
	currentID atomic.Int64
	data      map[int64]string
}

// NewCache creates a new local cache
func NewCache() (*Cache, error) {
	return &Cache{
		data: make(map[int64]string),
	}, nil
}

// Add adds message to the cache
func (c *Cache) Add(message string) int64 {
	c.Lock()
	defer c.Unlock()
	id := c.getNextID()
	c.data[id] = message
	return id
}

// Get gets message from the cache
func (c *Cache) Get(id int64) (string, bool) {
	c.Lock()
	defer c.Unlock()
	message, ok := c.data[id]
	return message, ok
}

// getNextID gets the next id to use
func (c *Cache) getNextID() int64 {
	return c.currentID.Add(1)
}
