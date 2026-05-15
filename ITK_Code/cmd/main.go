package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type item struct {
	value     interface{}
	expiresAt time.Time
}

type ObjectCache struct {
	mu      sync.RWMutex
	items   map[string]item
	ttl     time.Duration
	closed  chan struct{}
	bufPool sync.Pool
}

func NewObjectCache(ttl time.Duration) *ObjectCache {
	c := &ObjectCache{
		items:  make(map[string]item),
		ttl:    ttl,
		closed: make(chan struct{}),
	}

	c.bufPool.New = func() any {
		return new(bytes.Buffer)
	}

	go c.startGC()

	return c
}

func (c *ObjectCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *ObjectCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if time.Now().After(it.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}

	return it.value, true
}

func (c *ObjectCache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *ObjectCache) ToJSON() ([]byte, error) {
	buf := c.bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	c.mu.RLock()
	defer c.mu.RUnlock()

	enc := json.NewEncoder(buf)
	err := enc.Encode(c.items)

	data := make([]byte, buf.Len())
	copy(data, buf.Bytes())

	c.bufPool.Put(buf)

	return data, err
}

func (c *ObjectCache) startGC() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.closed:
			return
		}
	}
}

func (c *ObjectCache) cleanup() {
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if now.After(v.expiresAt) {
			delete(c.items, k)
		}
	}
}

func (c *ObjectCache) Close() {
	close(c.closed)
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
