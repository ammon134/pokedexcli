package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]cacheEntry
	mu       *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cacheMap: make(map[string]cacheEntry),
		mu:       &sync.Mutex{},
	}
	go cache.reapLoop(interval)
	return cache
}

func (c Cache) Add(url string, resBody []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheMap[url] = cacheEntry{
		createdAt: time.Now(),
		val:       resBody,
	}
}

func (c Cache) Get(url string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cEntry, ok := c.cacheMap[url]
	return cEntry.val, ok
}

func (c Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap(interval)
	}
}

func (c Cache) reap(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, cEntry := range c.cacheMap {
		if time.Since(cEntry.createdAt) > interval {
			delete(c.cacheMap, key)
		}
	}
}
