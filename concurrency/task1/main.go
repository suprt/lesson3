package main

import (
	"fmt"
	"sync"
)

type SafeCache struct {
	sync.RWMutex
	entries map[string]string
}

func NewSafeCache() *SafeCache {
	return &SafeCache{entries: make(map[string]string)}
}

func (c *SafeCache) Set(key, value string) {
	c.Lock()
	defer c.Unlock()
	c.entries[key] = value
}

func (c *SafeCache) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	value, ok := c.entries[key]
	return value, ok
}

func main() {
	cache := NewSafeCache()
	wg := sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			cache.Set(
				fmt.Sprintf("key%d", i),
				fmt.Sprintf("value%d", i),
			)
		}(i)
	}

	wg.Wait()

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("key%d", i)
			value, ok := cache.Get(key)

			fmt.Printf("%s: %q, exists=%v\n", key, value, ok)
		}(i)
	}

	wg.Wait()
}
