package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.items[key]; exists {
		c.queue.MoveToFront(element)
		return element.Value.(cacheItem).value, true
	}

	return nil, false
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	cItem := cacheItem{
		key:   key,
		value: value,
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	element, exists := c.items[key]

	if exists {
		element.Value = cItem
		c.queue.MoveToFront(element)
	} else {
		lItem := c.queue.PushFront(cItem)
		c.items[key] = lItem

		if c.queue.Len() > c.capacity {
			delete(c.items, c.queue.Back().Value.(cacheItem).key)
			c.queue.Remove(c.queue.Back())
		}
	}

	return exists
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) != 0 {
		c.queue = NewList()
		c.items = make(map[Key]*ListItem, c.capacity)
	}
}
