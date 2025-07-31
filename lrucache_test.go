package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSetOutRange(t *testing.T) {
	opts := Options{
		evictCallback: nil,
	}

	c := NewLRUCache(3, opts)

	c.Set("a", 1)
	testOrder(t, c, []string{"a"})

	c.Set("b", 2)
	testOrder(t, c, []string{"b", "a"})

	c.Set("c", 3)
	testOrder(t, c, []string{"c", "b", "a"})

	c.Set("a", 4)
	testOrder(t, c, []string{"a", "c", "b"})

	c.Set("d", 5)
	c.Set("e", 6)

	testOrder(t, c, []string{"e", "d", "a"})
}

func TestGet(t *testing.T) {
	opts := Options{
		evictCallback: nil,
	}
	c := NewLRUCache(3, opts)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Set("a", 4)

	val := c.Get("a")
	if val != 4 {
		t.Errorf("Wrong value: %v expected: %v", 1, val)
	}
}

func TestEvictCallback(t *testing.T) {
	var evictedKey string

	opts := Options{
		evictCallback: func(key string, value interface{}) {
			evictedKey = key
		},
	}
	c := NewLRUCache(3, opts)

	c.Set("a", 8)
	c.Set("b", 16)
	c.Set("c", 32)

	c.Set("d", 48)

	if evictedKey != "a" {
		t.Errorf("evicted key should be 'a' and not %s", evictedKey)
	}
}

func testOrder(t *testing.T, lru *LRUCache, want []string) {
	i := 0
	for elem := lru.order.Front(); elem != nil; elem = elem.Next() {
		item := elem.Value.(*CacheItem)

		if want[i] != item.key {
			t.Errorf("Invalid order of key %s", item.key)
		}

		i++
	}
}

func TestConcurrency(t *testing.T) {
	opts := Options{
		evictCallback: nil,
	}
	c := NewLRUCache(100, opts)

	const numGoroutines = 10
	const operationsPerGoroutine = 100

	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := id*1000 + j
				c.Set(key, value)
			}
		}(i)
	}
	wg.Wait()

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				c.Get(key)
			}
		}(i)
	}
	wg.Wait()

	wg.Add(numGoroutines * 2)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("write-key-%d-%d", id, j)
				value := id*2000 + j
				c.Set(key, value)
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j%50)
				c.Get(key)
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrencyWithEviction(t *testing.T) {
	var evictedCount int64
	var mu sync.Mutex

	opts := Options{
		evictCallback: func(key string, value interface{}) {
			mu.Lock()
			evictedCount++
			mu.Unlock()
		},
	}

	c := NewLRUCache(10, opts)

	const numGoroutines = 5
	const operationsPerGoroutine = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("evict-key-%d-%d", id, j)
				c.Set(key, j)
			}
		}(i)
	}

	wg.Wait()

	mu.Lock()
	finalEvictedCount := evictedCount
	mu.Unlock()

	if finalEvictedCount == 0 {
		t.Errorf("should be some evictions, got %d", finalEvictedCount)
	}
}
