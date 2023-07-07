package cache

import (
	"context"
	"errors"
	"fmt"
	"geektime-go-study/cache/internal"
	"log"
	"time"
)

var (
	ErrFailedToRefreshCache = errors.New("cache: 刷新缓存失败")
)

type ReadThroughCache struct {
	Cache
	LoadFunc   func(ctx context.Context, key string) (any, error)
	Expiration time.Duration
}

// 同步  同步读db 同步写缓存
func (c *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := c.Cache.Get(ctx, key)
	if errors.Is(err, internal.ErrKeyNotFound) {
		val, err = c.LoadFunc(ctx, key)
		if err == nil {
			er := c.Cache.Set(ctx, key, val, c.Expiration)
			if er != nil {
				return val, fmt.Errorf("%w, 原因：%s", ErrFailedToRefreshCache, er.Error())
			}
		}
	}
	return val, err
}

// 全异步 异步读db, 异步写缓存
func (c *ReadThroughCache) GetV1(ctx context.Context, key string) (any, error) {
	val, err := c.Cache.Get(ctx, key)
	if errors.Is(err, internal.ErrKeyNotFound) {
		go func() {
			val, err = c.LoadFunc(ctx, key)
			if err == nil {
				er := c.Cache.Set(ctx, key, val, c.Expiration)
				if er != nil {
					log.Fatalln(er)
				}
			}
		}()
	}
	return val, err
}

func (c *ReadThroughCache) GetV2(ctx context.Context, key string) (any, error) {
	val, err := c.Cache.Get(ctx, key)
	if errors.Is(err, internal.ErrKeyNotFound) {
		//半异步 同步读db, 异步写缓存
		val, err = c.LoadFunc(ctx, key)
		if err == nil {
			go func() {
				er := c.Cache.Set(ctx, key, val, c.Expiration)
				if er != nil {
					log.Fatalln(er)
				}
			}()
		}
	}
	return val, err
}
