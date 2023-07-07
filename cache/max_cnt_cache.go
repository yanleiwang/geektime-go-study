package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	ErrOverCapacity = errors.New("cache: 超过容量限制")
)

type MaxCntCache struct {
	*BuildInMapCache
	cnt    int32
	maxCnt int32
}

func NewMaxCntCache(b *BuildInMapCache, maxCnt int32) *MaxCntCache {
	res := &MaxCntCache{
		BuildInMapCache: b,
		cnt:             0,
		maxCnt:          maxCnt,
	}
	origin := b.onEvicted
	res.onEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		if origin != nil {
			origin(key, val)
		}
	}

	return res
}

func (c *MaxCntCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.data[key]
	if !ok {
		if c.cnt+1 > c.maxCnt {
			return ErrOverCapacity
		}
		err := c.set(key, val, expiration)
		if err != nil {
			return err
		}
		c.cnt++
		return nil
	}

	return c.set(key, val, expiration)

}
