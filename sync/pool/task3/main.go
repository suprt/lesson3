package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ObjectCache struct {
	mu       sync.RWMutex
	items    map[string]*cacheItem
	ttlNanos int64
	bufPool  sync.Pool
}

type cacheItem struct {
	Value      any
	LastAccess atomic.Int64
}

func (ci *cacheItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value      any       `json:"value"`
		LastAccess time.Time `json:"lastAccess"`
		//LastAccess int64 `json:"lastAccess"`
	}{
		Value:      ci.Value,
		LastAccess: time.Unix(0, ci.LastAccess.Load()),
		//Если требуется выводить время без форматирования
		//LastAccess: ci.LastAccess.Load()
	})
}

func NewObjectCache(ttl time.Duration) *ObjectCache {
	c := &ObjectCache{
		ttlNanos: ttl.Nanoseconds(),
		items:    make(map[string]*cacheItem),
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
		key        string
		lastAccess int64
	}

	for range ticker.C {
		expired := make([]expiredKey, 0)
		now := time.Now().UnixNano()
		c.mu.RLock()

		for k, v := range c.items {
			if now > v.LastAccess.Load()+c.ttlNanos {
				expired = append(expired, expiredKey{k, v.LastAccess.Load()})
			}
		}
		c.mu.RUnlock()

		if len(expired) == 0 {
			continue
		}

		c.mu.Lock()

		for _, v := range expired {
			item, ok := c.items[v.key]
			if ok && item.LastAccess.Load() == v.lastAccess {
				delete(c.items, v.key)
			}
		}
		c.mu.Unlock()

	}
}

func (c *ObjectCache) Set(key string, value any) {
	c.mu.Lock()
	now := time.Now().UnixNano()
	item := &cacheItem{Value: value}
	item.LastAccess.Store(now)
	c.items[key] = item
	c.mu.Unlock()
}

func (c *ObjectCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	now := time.Now().UnixNano()
	if !ok || (now > item.LastAccess.Load()+c.ttlNanos) {
		return nil, false
	}

	item.LastAccess.Store(now)

	return item.Value, true
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
	time.Sleep(time.Second)
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
