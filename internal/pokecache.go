package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	data map[string]cacheEntry
	mu sync.RWMutex
}

func NewCache(interval time.Duration) Cache {
	var new Cache
	new.data = make(map[string]cacheEntry)
	new.mu = sync.RWMutex{}
	new.reapLoop(interval)
	return new
}

func (c *Cache)Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{time.Now(), val}
	c.data[key] = entry
}

func (c *Cache)Get(key string) ([]byte, bool) {
	//if key is in cache return data and true
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.data[key]
	if ok {
		return entry.val, true
	}
	return nil, false
}

func (c *Cache)reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {

		for range ticker.C {
			c.mu.Lock()
			for key, entry := range c.data {
				if time.Since(entry.createdAt) > (interval) {
					delete(c.data, key)
					//fmt.Println("\nCache entry deleted")
					} 
				}
				//fmt.Println("\nCache refreshed")
				c.mu.Unlock()
			}
	}()



}