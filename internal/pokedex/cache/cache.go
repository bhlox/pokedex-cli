package cache

import (
	"sync"
	"time"
)

type Cache struct {
	data map[string]CacheEntry
	mu     *sync.RWMutex
}

func (c *Cache) Add(key string, val []byte){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheEntry{createdAt: time.Now(),val: val}
}

func (c *Cache) Get(key string)([]byte,bool){
	c.mu.RLock()
	defer c.mu.RUnlock()
	data,ok := c.data[key]
	return data.val,ok
}

func (c *Cache) ReapLoop(interval time.Duration){
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.data {
			if time.Since(entry.createdAt) > interval {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

type CacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration)*Cache{
	c := &Cache{mu: &sync.RWMutex{},data: make(map[string]CacheEntry)}
	go c.ReapLoop(interval)
	return c
}