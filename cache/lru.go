package cache

import (
	"container/list"
	"sync"
	"time"
)

// LRU缓存
type LRU struct {
	mutex       sync.RWMutex
	maxItemSize int
	defaultTtl  time.Duration
	cacheList   *list.List
	cache       map[interface{}]*list.Element
	hits, gets  atomicUint
}

type entry struct {
	key    string
	value  interface{}
	expiry time.Time
}

// 如果 maxItemSize = zero 则表示不限制
func NewLRU(maxItemSize int, ttl time.Duration) *LRU {
	return &LRU{
		maxItemSize: maxItemSize,
		cacheList:   list.New(),
		cache:       make(map[interface{}]*list.Element),
		defaultTtl:  ttl,
	}
}

// 缓存状态
func (c *LRU) Status() *Status {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return &Status{
		MaxItemSize: c.maxItemSize,
		CurrentSize: c.cacheList.Len(),
		Gets:        c.gets.get(),
		Hits:        c.hits.get(),
	}
}

func (c *LRU) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	needToDelete := false
	defer func() {
		c.mutex.RUnlock()
		// 必须在RLock解锁后再执行删除
		if needToDelete {
			c.Delete(key)
		}
	}()
	c.gets.add(1)
	if ele, hit := c.cache[key]; hit {
		ent := ele.Value.(*entry)
		if ent.expiry.IsZero() || ent.expiry.After(time.Now()) {
			c.hits.add(1)
			c.cacheList.MoveToFront(ele)
			return ele.Value.(*entry).value, true
		}
		needToDelete = true
	}
	return nil, false
}

func (c *LRU) SetEx(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.cacheList = list.New()
		c.defaultTtl = -1
	}
	var expiry time.Time
	if ttl > 0 {
		expiry = time.Now().Add(ttl)
	} else {
		expiry = time.Time{}
	}

	if ele, ok := c.cache[key]; ok {
		c.cacheList.MoveToFront(ele)
		ele.Value.(*entry).value = value
		ele.Value.(*entry).expiry = expiry
		return
	}

	ele := c.cacheList.PushFront(&entry{key: key, value: value, expiry: expiry})
	c.cache[key] = ele
	if c.maxItemSize != 0 && c.cacheList.Len() > c.maxItemSize {
		c.removeOldest()
	}
}

func (c *LRU) Set(key string, value interface{}) {
	c.SetEx(key, value, c.defaultTtl)
}

func (c *LRU) remove(ele *list.Element) {
	c.cacheList.Remove(ele)
	key := ele.Value.(*entry).key
	delete(c.cache, key)
}

func (c *LRU) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.remove(ele)
		return
	}
}

func (c *LRU) removeOldest() {
	if c.cache == nil {
		return
	}
	ele := c.cacheList.Back()
	if ele != nil {
		c.remove(ele)
	}
}
