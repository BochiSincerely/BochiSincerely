package geecache

import (
	"GoCache/gocache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil { //	若为nil，则再创建实例。 --延迟初始化
		c.lru = lru.New(c.cacheBytes, nil)

	}
	c.lru.Add(key, value)
} //对象的延迟初始化就意味着对象的创建会延迟到第一次使用该对象，
//用于提高性能，减少程序的内存要求

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
