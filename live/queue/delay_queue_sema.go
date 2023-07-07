package queue

import (
	"context"
	"golang.org/x/sync/semaphore"
	"math"
	"sync"
	"time"
)

type DelayQueueUseSema[T Delayable] struct {
	q           PriorityQueue[T]
	mutex       *sync.Mutex
	dequeueSema *semaphore.Weighted //  可出队元素个数,  表示超时元素个数
	enqueueSema *semaphore.Weighted //  可入队元素个数

	zero T
}

func (d *DelayQueueUseSema[T]) Enqueue(ctx context.Context, val T) error {

	// 为了 过测试用例, 测试用例 context 本身已经过期了
	if ctx.Err() != nil {
		return ctx.Err()
	}

	err := d.enqueueSema.Acquire(ctx, 1)
	if err != nil {
		return err
	}

	d.mutex.Lock()

	err = d.q.Enqueue(val)
	d.mutex.Unlock()

	if err != nil {
		d.enqueueSema.Release(1)
		return err
	}

	// 通知 有元素超时
	_ = time.AfterFunc(val.Delay(), func() {
		d.dequeueSema.Release(1)
	})

	return nil
}

func (d *DelayQueueUseSema[T]) Dequeue(ctx context.Context) (T, error) {
	if ctx.Err() != nil {
		return d.zero, ctx.Err()
	}
	err := d.dequeueSema.Acquire(ctx, 1)
	if err != nil {
		return d.zero, err
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()
	if ctx.Err() != nil {
		d.dequeueSema.Release(1)
		return d.zero, ctx.Err()
	}

	first, err := d.q.Dequeue()
	if err != nil || first.Delay() > 0 {
		panic("意外错误, 信号量Release了, 但没拿到超时元素")
	}
	// 出队了一个元素, 所以可以入队一个元素
	d.enqueueSema.Release(1)
	return first, nil
}

func NewDelayQueueUseSema[T Delayable](c int) *DelayQueueUseSema[T] {
	size := int64(c)
	if c <= 0 {
		size = math.MaxInt64
	}
	dequeueSema := semaphore.NewWeighted(size)
	enqueueSema := semaphore.NewWeighted(size)
	dequeueSema.Acquire(context.Background(), size)

	ret := &DelayQueueUseSema[T]{
		q: *NewPriorityQueue[T](c, func(src T, dst T) int {
			srcDelay := src.Delay()
			dstDelay := dst.Delay()
			if srcDelay > dstDelay {
				return 1
			}
			if srcDelay == dstDelay {
				return 0
			}
			return -1
		}),
		mutex:       &sync.Mutex{},
		enqueueSema: enqueueSema,
		dequeueSema: dequeueSema,
	}

	return ret
}

//type Delayable interface {
//	// Delay 实时计算
//	Delay() time.Duration
//}
