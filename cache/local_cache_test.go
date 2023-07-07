package cache

import (
	"context"
	"geektime-go-study/cache/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBuildInMapCache_Get(t *testing.T) {
	testCases := []struct {
		name     string
		key      string
		wantVal  any
		wantErr  error
		newCache func() *BuildInMapCache
	}{
		{
			name:    "Not Exist",
			key:     "Not Exist",
			wantErr: internal.NewErrKeyNotFound("Not Exist"),
			newCache: func() *BuildInMapCache {
				return NewBuildInMapCache(time.Second)
			},
		},

		{
			name:    "normal",
			key:     "normal",
			wantVal: 10,
			newCache: func() *BuildInMapCache {
				res := NewBuildInMapCache(time.Second)
				err := res.Set(context.Background(), "normal", 10, 3*time.Second)
				require.NoError(t, err)
				return res
			},
		},
		{
			name: "expired",
			key:  "expired key",
			newCache: func() *BuildInMapCache {
				res := NewBuildInMapCache(10 * time.Second)
				err := res.Set(context.Background(), "expired key", 1, time.Second)
				require.NoError(t, err)
				time.Sleep(time.Second * 2)
				return res
			},
			wantErr: internal.NewErrKeyNotFound("expired key"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := tc.newCache()
			get, err := cache.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, get)
		})
	}

}

func TestBuildInMapCache_Loop(t *testing.T) {
	cnt := 0
	c := NewBuildInMapCache(time.Second, BuildInMapCacheWithEvicted(func(key string, val any) {
		cnt++
	}))
	err := c.Set(context.Background(), "key1", 123, time.Second)
	require.NoError(t, err)
	time.Sleep(time.Second * 3)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.data["key1"]
	require.False(t, ok)
	require.Equal(t, 1, cnt)
}
