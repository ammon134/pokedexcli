package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	CacheMap map[string]cacheEntry
	mu       *sync.Mutex
}

type cacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(duration time.Duration) Cache {
	cache := Cache{
		CacheMap: make(map[string]cacheEntry),
		mu:       &sync.Mutex{},
	}
	go cache.reapLoop(duration)
	return cache
}

func (c Cache) Add(url string, resBody []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CacheMap[url] = cacheEntry{
		CreatedAt: time.Now(),
		Val:       resBody,
	}
	return nil
}

func (c Cache) Get(url string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if cEntry, ok := c.CacheMap[url]; ok {
		return cEntry.Val, nil
	}

	return []byte{}, nil
}

func (c Cache) reapLoop(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {
		for key, cEntry := range c.CacheMap {
			if time.Since(cEntry.CreatedAt) > duration {
				delete(c.CacheMap, key)
			}
		}
	}
}
