package hw04lrucache

import "sync"

type Key string

type Cache interface {
    Set(key Key, value interface{}) bool
    Get(key Key) (interface{}, bool)
    Clear()
    Len() int
}

type cacheEntry struct {
    key   Key
    value interface{}
}

type lruCache struct {
    sync.RWMutex   
    capacity int
    queue    List
    items    map[Key]*ListItem
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
    if cache == nil {
        panic("cache is nil")
    }
    cache.Lock()
    defer cache.Unlock()
    if item, ok := cache.items[key]; ok {
        entry := item.Value.(cacheEntry)
        entry.value = value
        item.Value = entry
        cache.queue.MoveToFront(item)
        return true
    }
    entry := cacheEntry{key, value}
    listItem := cache.queue.PushFront(entry)
    cache.items[key] = listItem
    if cache.queue.Len() > cache.capacity {
        oldest := cache.queue.Back()
        oldestEntry := oldest.Value.(cacheEntry)
        delete(cache.items, oldestEntry.key)
        cache.queue.Remove(oldest)
    }
    return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
    if cache == nil {
        panic("cache is nil")
    }
    cache.Lock()
    defer cache.Unlock()
    if item, ok := cache.items[key]; ok {
        cache.queue.MoveToFront(item)
        entry := item.Value.(cacheEntry)
        return entry.value, true
    }
    return nil, false
}

func (cache *lruCache) Clear() {
    if cache == nil {
        panic("cache is nil")
    }
    cache.Lock()
    defer cache.Unlock()
   
    for key := range cache.items {
        delete(cache.items, key)
    }
}

func (cache *lruCache) Len() int {
    cache.RLock()
    defer cache.RUnlock()
    return cache.queue.Len()
}

func NewCache(capacity int) Cache {
    return &lruCache{
        capacity: capacity,
        queue:    NewList(),
        items:    make(map[Key]*ListItem, capacity),
    }
}
