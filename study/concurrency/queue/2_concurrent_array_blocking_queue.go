package queue

import (
	"context"
	"sync"
)

// ConcurrentArrayBlockingQueue
// 并发队列 slice实现 有界 阻塞 不超时 加锁
// 通过条件变量来通知 可以入队/出队
// 跟 ConcurrentArrayQueue 相比 支持阻塞
// 缺点:  如果没有被通知, 入队/出队操作 会一直阻塞下去, 也就是不支持超时控制.
// 很多业务能接受阻塞一段时间，但是不能接受一直阻塞；
// 带超时控制可以防止资源泄露：
type ConcurrentArrayBlockingQueue[T any] struct {
	data     []T
	mutex    *sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	capacity int
	zero     T
}

func NewConcurrentArrayBlockingQueue[T any](capacity int) *ConcurrentArrayBlockingQueue[T] {
	mutex := &sync.Mutex{}

	return &ConcurrentArrayBlockingQueue[T]{
		capacity: capacity,
		mutex:    mutex,
		notFull:  sync.NewCond(mutex),
		notEmpty: sync.NewCond(mutex),
	}
}

func (c *ConcurrentArrayBlockingQueue[T]) In(ctx context.Context, val T) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for len(c.data) >= c.capacity {
		c.notFull.Wait()
	}
	c.data = append(c.data, val)
	c.notEmpty.Signal()
	return nil
}

func (c *ConcurrentArrayBlockingQueue[T]) Out(ctx context.Context) (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for len(c.data) == 0 {
		c.notEmpty.Wait()
	}
	ret := c.data[0]
	c.data = c.data[1:]
	c.notFull.Signal()
	return ret, nil
}
