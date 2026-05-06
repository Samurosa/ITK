package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Storage struct {
	mu     sync.RWMutex
	caches map[string]Cache
}

type Cache struct {
	val   interface{}
	tlive time.Time
}

type User struct {
	Name string
}

type Item struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func NewCache() *Storage {
	return &Storage{
		caches: make(map[string]Cache),
	}
}

func (c *Storage) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cache := Cache{val: value, tlive: time.Now().Add(ttl)}
	c.caches[key] = cache
}

func (c *Storage) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	cache, ok := c.caches[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if time.Now().After(cache.tlive) {
		c.mu.Lock()
		delete(c.caches, key)
		c.mu.Unlock()
		return nil, false
	}
	return cache.val, true
}

func (c *Storage) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.caches, key)
}

func (c *Storage) Exists(key string) bool {
	_, ok := c.Get(key)
	return ok
}

func (c *Storage) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.caches = make(map[string]Cache)
}

func (c *Storage) ToJSON() ([]byte, error) {
	res := make([]Item, 0, len(c.caches))

	c.mu.RLock()
	for key, cache := range c.caches {
		res = append(res, Item{Key: key, Value: cache.val})
	}
	c.mu.RUnlock()

	return json.Marshal(res)
}

func GetAs[T any](c *Storage, key string) (T, error) {
	var zero T
	value, ok := c.Get(key)
	if ok {
		c.mu.Lock()
		defer c.mu.Unlock()

		value, ok := value.(T)
		if !ok {
			return zero, fmt.Errorf("неправильный формат")
		}

		return value, nil
	}
	return zero, fmt.Errorf("неправильный формат")
}

func main() {
	cache := NewCache()

	cache.Set("user", User{Name: "Alice"}, time.Hour) // Хранится 1 час
	cache.Set("temp_data", 42, time.Minute)           // Хранится 1 минуту

	jsonData, _ := cache.ToJSON()
	fmt.Println(string(jsonData)) // {"temp_data":42,"user":{"Name":"Alice"}}

	cache.Clear()
	fmt.Println("Exists (user):", cache.Exists("user")) //false
}
