package cache

import (
	"context"
	"log"
	"time"
)

type WriteThroughCache struct {
	Cache
	StoreFunc func(ctx context.Context, key string, val any) error
}

// 同步 同步写db 同步写cache
func (c *WriteThroughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := c.StoreFunc(ctx, key, val)
	if err != nil {
		return err
	}
	return c.Cache.Set(ctx, key, val, expiration)
}

// 半异步  同步写db 异步写cache
func (c *WriteThroughCache) SetV2(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := c.StoreFunc(ctx, key, val)
	go func() {
		er := c.Cache.Set(ctx, key, val, expiration)
		if er != nil {
			log.Fatalln(er)
		}
	}()
	return err
}

// 完全异步  只具备理论意义，实际上几乎不会用
func (c *WriteThroughCache) SetV3(ctx context.Context, key string, val any, expiration time.Duration) error {
	go func() {
		err := c.StoreFunc(ctx, key, val)
		if err != nil {
			log.Fatalln(err)
		}
		if err = c.Cache.Set(ctx, key, val, expiration); err != nil {
			log.Fatalln(err)
		}
	}()
	return nil
}
