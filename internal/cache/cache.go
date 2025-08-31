package cache

import (
	"sync"
	"time"
)

type Entry struct {
	Data      interface{}
	Timestamp time.Time
	TTL       time.Duration
}

type Cache struct {
	data  map[string]Entry
	mutex sync.RWMutex
	ttl   time.Duration
}

func New(ttl time.Duration) *Cache {
	return &Cache{
		data: make(map[string]Entry),
		ttl:  ttl,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// check if entry has expired
	if time.Since(entry.Timestamp) >= entry.TTL {
		// Remove expired entry
		c.mutex.RUnlock()
		c.mutex.Lock()
		delete(c.data, key)
		c.mutex.Unlock()
		c.mutex.RLock()
		return nil, false
	}

	return entry.Data, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = Entry{
		Data:      value,
		Timestamp: time.Now(),
		TTL:       c.ttl,
	}
}

func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = Entry{
		Data:      value,
		Timestamp: time.Now(),
		TTL:       ttl,
	}
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key := range c.data {
		delete(c.data, key)
	}
}

func (c *Cache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.data)
}

func (c *Cache) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}
