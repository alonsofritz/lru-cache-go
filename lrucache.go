package main

import (
	"container/list"
	"fmt"
	"log"
	"sync"
)

type EvictCallback func(key string, value interface{})

type Options struct {
	logs          bool
	evictCallback EvictCallback
}

type CacheItem struct {
	key   string
	value interface{}
}

type LRUCache struct {
	maxEntries int
	keyMap     map[string]*list.Element
	order      *list.List
	options    Options
	sync.RWMutex
}

func NewLRUCache(maxEntries int, opts Options) *LRUCache {
	if maxEntries <= 0 {
		return nil
	}

	lru := &LRUCache{
		maxEntries: maxEntries,
		keyMap:     make(map[string]*list.Element),
		order:      list.New(),
		options:    opts,
	}

	return lru
}

func (l *LRUCache) Set(key string, value interface{}) bool {
	l.Lock()
	defer l.Unlock()

	if elem, ok := l.keyMap[key]; ok {
		item := elem.Value.(*CacheItem)
		item.value = value

		l.order.MoveToFront(elem)

		if l.options.logs {
			log.Printf("Elem %s updated to the front", key)
			l.printList()
		}

		return false
	}

	item := &CacheItem{
		key:   key,
		value: value,
	}

	elem := l.order.PushFront(item)

	l.keyMap[key] = elem

	if l.order.Len() > l.maxEntries {
		l.removeLastElement()
	}

	if l.options.logs {
		log.Printf("Elem %s added to the front", key)
		l.printList()
	}

	return true
}

func (l *LRUCache) Get(key string) (value interface{}) {
	l.Lock()
	defer l.Unlock()

	if elem, ok := l.keyMap[key]; ok {
		item := elem.Value.(*CacheItem)
		return item.value
	}

	return nil
}

func (l *LRUCache) printList() {
	for elem := l.order.Front(); elem != nil; elem = elem.Next() {
		item := elem.Value.(*CacheItem)
		fmt.Printf("key: %s, value: %v \n", item.key, item.value)
	}
}

func (l *LRUCache) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	l.order.Remove(e)

	item := e.Value.(*CacheItem)
	delete(l.keyMap, item.key)

	if l.options.evictCallback != nil {
		l.options.evictCallback(item.key, item.value)
	}
}

func (l *LRUCache) removeLastElement() {
	l.removeElement(l.order.Back())
}
