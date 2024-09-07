package main

import (
	"sync"
	"time"
)

const defaultClearInterval = 10 * time.Second

type item[V any] struct {
	value      V
	expiration time.Time
}

func (i item[V]) hasExpired() bool {
	return time.Now().After(i.expiration)
}

type Cache[K comparable, V any] struct {
	items map[K]item[V]
	mu    sync.RWMutex
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	cache := &Cache[K, V]{
		items: make(map[K]item[V]),
	}

	go func() {
		for range time.Tick(defaultClearInterval) {
			cache.mu.Lock()
			for key, value := range cache.items {
				if value.hasExpired() {
					delete(cache.items, key)
				}
			}
			cache.mu.Unlock()
		}
	}()
	return cache
}

func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = item[V]{value: value, expiration: time.Now().Add(ttl)}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return item.value, false
	}
	if item.hasExpired() {
		delete(c.items, key)
		return item.value, false
	}

	return item.value, true
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache[K, V]) Pop(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists {
		return item.value, false
	}

	delete(c.items, key)
	if item.hasExpired() {
		return item.value, false
	}

	return item.value, true
}
