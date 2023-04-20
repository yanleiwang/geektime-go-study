package cache

import (
	"context"
	"geektime-go-study/cache/internal"
	"sync"
	"time"
)

type item struct {
	val      any
	deadline time.Time
}

func (i *item) IsExpired(t time.Time) bool {
	if i.deadline.IsZero() || i.deadline.After(t) {
		return false
	}
	return true
}

type BuildInMapCache struct {
	data      map[string]*item
	mutex     sync.RWMutex
	close     chan struct{}
	maxCnt    int                       // 每次轮询过期key 最大次数
	onEvicted func(key string, val any) // 当删除key的时候 执行
}

func (c *BuildInMapCache) Get(ctx context.Context, key string) (any, error) {
	c.mutex.RLock()
	res, ok := c.data[key]
	c.mutex.RUnlock()
	if !ok {
		return nil, internal.NewErrKeyNotFound(key)
	}

	now := time.Now()
	if res.IsExpired(now) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		res, ok = c.data[key]
		if !ok {
			return nil, internal.NewErrKeyNotFound(key)
		}
		if res.IsExpired(now) {
			c.delete(key)
			return nil, internal.NewErrKeyNotFound(key)
		}
	}

	return res.val, nil
}

func (c *BuildInMapCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	var dl time.Time
	if expiration > 0 {
		dl = time.Now().Add(expiration)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = &item{
		val:      val,
		deadline: dl,
	}
	return nil
}

func (c *BuildInMapCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.delete(key)
	return nil
}

func (c *BuildInMapCache) Close() error {
	c.close <- struct{}{}
	return nil
}

func (c *BuildInMapCache) delete(k string) {
	itm, ok := c.data[k]
	if !ok {
		return
	}
	delete(c.data, k)
	c.onEvicted(k, itm.val)
}

type BuildInMapCacheOption func(cache *BuildInMapCache)

func BuildInMapCacheWithEvicted(fn func(key string, val any)) BuildInMapCacheOption {
	return func(cache *BuildInMapCache) {
		cache.onEvicted = fn
	}
}

func NewBuildInMapCache(interval time.Duration, opts ...BuildInMapCacheOption) *BuildInMapCache {
	res := &BuildInMapCache{
		data:   make(map[string]*item, 10),
		close:  make(chan struct{}),
		maxCnt: 10000,
		onEvicted: func(key string, val any) {

		},
	}

	for _, opt := range opts {
		opt(res)
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-res.close:
				return
			case t := <-ticker.C:
				res.mutex.Lock()
				cnt := 0
				for k, v := range res.data {
					if v.IsExpired(t) {
						res.delete(k)
					}
					cnt++
					if cnt > res.maxCnt {
						break
					}
				}
				res.mutex.Unlock()
			}
		}

	}()

	return res

}
