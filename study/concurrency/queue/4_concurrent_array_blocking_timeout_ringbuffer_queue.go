package queue

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

// ConcurrentArrayBlockingTimeoutRingBufferQueue
//
//	并发队列 ring buffer实现 有界 阻塞 超时 加锁
//
// 底层数据结构是 ring buffer,
// 相较于 ConcurrentArrayBlockingTimeoutQueue 没有slice扩容的问题.
// 但是数组必须一开始 分配好 所有的空间,  当容量比较大, 但实际使用量没这么大的时候, 就会比较浪费
// 解决方案: 可以 底层用 链表来存储
type ConcurrentArrayBlockingTimeoutRingBufferQueue[T any] struct {
	data []T
	// 队头元素下标
	head int
	// 队尾元素下标
	tail int
	// 包含多少个元素
	count int

	inCnt  *semaphore.Weighted // 可以入队的数量
	outCnt *semaphore.Weighted // 可以出队的数量
	zero   T
	mutex  sync.Mutex
}

func NewConcurrentArrayBlockingTimeoutRingBufferQueue[T any](capacity int) (*ConcurrentArrayBlockingTimeoutRingBufferQueue[T], error) {

	ret := &ConcurrentArrayBlockingTimeoutRingBufferQueue[T]{
		data:   make([]T, capacity), // 因为是 ring buffer, 所以一开始就要分配好数组
		inCnt:  semaphore.NewWeighted(int64(capacity)),
		outCnt: semaphore.NewWeighted(int64(capacity)),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	err := ret.outCnt.Acquire(ctx, int64(capacity))

	return ret, err
}

func (c *ConcurrentArrayBlockingTimeoutRingBufferQueue[T]) In(ctx context.Context, val T) error {
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

	c.data[c.tail] = val
	c.tail++
	c.count++

	// c.tail 已经是最后一个了，重置下标
	if c.tail == cap(c.data) {
		c.tail = 0
	}

	// 往出队的sema放入一个元素，出队的goroutine可以拿到并出队
	c.outCnt.Release(1)
	return nil

}

func (c *ConcurrentArrayBlockingTimeoutRingBufferQueue[T]) Out(ctx context.Context) (T, error) {
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

	ret := c.data[c.head]
	// 为了释放内存，GC
	c.data[c.head] = c.zero

	c.head++
	c.count--
	if c.head == cap(c.data) {
		c.head = 0
	}

	c.inCnt.Release(1)
	return ret, nil

}

func (c *ConcurrentArrayBlockingTimeoutRingBufferQueue[T]) Len() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.count
}
