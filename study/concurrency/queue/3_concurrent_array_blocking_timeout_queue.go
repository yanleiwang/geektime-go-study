package queue

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

// ConcurrentArrayBlockingTimeoutQueue
// 并发队列 slice实现 有界 阻塞 超时 加锁
// 使用信号量实现, 超时控制 (信号量 实际就相当于资源的数量)
// 解决了 ConcurrentArrayBlockingQueue 不支持超时控制的问题.
//
// 缺点: 入队操作是 slice的append操作, 所以实际底层数组是在不断增长
// 也就是可能引起频繁的扩容
type ConcurrentArrayBlockingTimeoutQueue[T any] struct {
	data     []T
	capacity int
	inCnt    *semaphore.Weighted // 可以入队的数量
	outCnt   *semaphore.Weighted // 可以出队的数量
	zero     T
	mutex    sync.Mutex
}

func NewConcurrentArrayBlockingTimeoutQueue[T any](capacity int) (*ConcurrentArrayBlockingTimeoutQueue[T], error) {

	ret := &ConcurrentArrayBlockingTimeoutQueue[T]{
		capacity: capacity,
		inCnt:    semaphore.NewWeighted(int64(capacity)),
		outCnt:   semaphore.NewWeighted(int64(capacity)),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	err := ret.outCnt.Acquire(ctx, int64(capacity))

	return ret, err
}

func (c *ConcurrentArrayBlockingTimeoutQueue[T]) In(ctx context.Context, val T) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	// 获取可以入队的位置
	err := c.inCnt.Acquire(ctx, 1)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 获取锁 可能耗费比较长时间, 所以再做一次超时检测
	if ctx.Err() != nil {
		// 超时应该主动归还信号量，避免容量泄露
		c.inCnt.Release(1)
		return ctx.Err()
	}
	c.data = append(c.data, val)

	// 往出队的sema放入一个元素，出队的goroutine可以拿到并出队
	c.outCnt.Release(1)
	return nil

}

func (c *ConcurrentArrayBlockingTimeoutQueue[T]) Out(ctx context.Context) (T, error) {
	if ctx.Err() != nil {
		return c.zero, ctx.Err()
	}

	// 获取可以出队的位置
	err := c.outCnt.Acquire(ctx, 1)
	if err != nil {
		return c.zero, err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 获取锁 可能耗费比较长时间, 所以再做一次超时检测
	if ctx.Err() != nil {
		c.outCnt.Release(1)
		return c.zero, ctx.Err()
	}

	ret := c.data[0]
	c.data = c.data[1:]
	c.inCnt.Release(1)
	return ret, nil

}
