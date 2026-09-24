package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type ObjectCache struct {
	mu      sync.RWMutex
	items   map[string]cacheItem
	ttl     time.Duration
	bufPool sync.Pool
}

type cacheItem struct {
	Value     any
	ExpiresAt time.Time
}

func NewObjectCache(ttl time.Duration) *ObjectCache {
	c := &ObjectCache{
		ttl:   ttl,
		items: make(map[string]cacheItem),
	}
	c.bufPool.New = func() any {
		return new(bytes.Buffer)
	}

	go c.cleanup()
	return c
}

func (c *ObjectCache) cleanup() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	type expiredKey struct {
		key string
		exp time.Time
	}

	for range ticker.C {
		expired := make([]expiredKey, 0)
		now := time.Now()
		c.mu.RLock()

		for k, v := range c.items {
			if now.After(v.ExpiresAt) {
				expired = append(expired, expiredKey{k, v.ExpiresAt})
			}
		}
		c.mu.RUnlock()

		if len(expired) == 0 {
			continue
		}

		c.mu.Lock()

		for _, v := range expired {
			item, ok := c.items[v.key]
			if ok && item.ExpiresAt == v.exp {
				delete(c.items, v.key)
			}
		}
		c.mu.Unlock()

	}
}

func (c *ObjectCache) Set(key string, value any) {
	c.mu.Lock()
	c.items[key] = cacheItem{value, time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *ObjectCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]

	if !ok || time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	item.ExpiresAt = time.Now().Add(c.ttl)
	c.items[key] = item

	return item.Value, ok
}

func (c *ObjectCache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *ObjectCache) ToJSON() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	buf := c.bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		c.bufPool.Put(buf)
	}()

	if err := json.NewEncoder(buf).Encode(c.items); err != nil {
		return nil, err
	}
	result := append([]byte(nil), buf.Bytes()...)

	return result, nil
}

func main() {
	cache := NewObjectCache(5 * time.Second)

	// Добавляем данные в кэш
	cache.Set("user:1", map[string]string{"name": "Alice", "role": "admin"})
	cache.Set("user:2", map[string]string{"name": "Bob", "role": "user"})

	// Получаем объект
	if user, found := cache.Get("user:1"); found {
		fmt.Println("Найден:", user)
	}

	// Выводим JSON
	jsonData, _ := cache.ToJSON()
	fmt.Println("Кэш в JSON:", string(jsonData))

	// Ждём истечения TTL и проверяем снова
	time.Sleep(6 * time.Second)
	_, found := cache.Get("user:1")
	fmt.Println("После TTL, user:1 найден?", found)
}
