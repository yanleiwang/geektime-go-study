package queue

import (
	"context"
	"errors"
	"sync"
)

// ConcurrentArrayQueue
// 并发队列 slice实现 有界 非阻塞 非超时 加锁
// 最简单的实现, 跟普通的队列(非并发) 实现 相比其实就是多了一个 锁
type ConcurrentArrayQueue[T any] struct {
	data     []T
	mutex    sync.Mutex
	capacity int
	empty    T
}

func NewConcurrentArrayQueue[T any](capacity int) *ConcurrentArrayQueue[T] {
	return &ConcurrentArrayQueue[T]{
		capacity: capacity,
	}
}

func (c *ConcurrentArrayQueue[T]) In(ctx context.Context, val T) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.data) >= c.capacity {
		return errors.New("队列已满")
	}
	c.data = append(c.data, val)
	return nil
}

func (c *ConcurrentArrayQueue[T]) Out(ctx context.Context) (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.data) == 0 {
		return c.empty, errors.New("空的队列")
	}
	ret := c.data[0]
	c.data = c.data[1:]
	return ret, nil
}
