package queue

import (
	"context"
	"golang.org/x/sync/semaphore"
	"math"
	"sync"
	"time"
)

/*
DelayQueue，即延时队列。延时队列的运作机制是：

按照元素的预期过期时间来进行排序，过期时间短的在前面；
当从队列中获取元素的时候，如果队列为空，或者元素还没有到期，那么调用者会被阻塞；
入队的时候，如果队列已经满了，那么调用者会被阻塞，直到有人取走元素，或者阻塞超时；
因此延时队列的使用场景主要就是那些依赖于时间的场景，例如本地缓存，定时任务等。

使用延时队列需要两个步骤：

实现 Delayable 接口
创建一个延时队列


所以实际上 延时队列的底层数据结构 就是根据 过期时间进行排序的 优先级队列

*/

// DelayQueueUseSema  实现一:  根据 信号量进行超时控制的 延时队列
// 本质上就是 用信号量来表示可入队元素个数和  可出队(也就是超时元素)个数
// 当一个新的元素入队的时候, 开一个定时器, 在超时的时候 可出队元素+1
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
	// 拿到锁后, 可能超时了, 所以要再一次检测
	if ctx.Err() != nil {
		d.enqueueSema.Release(1)
		return ctx.Err()
	}
	// 这时候 入队 一般不会有err
	_ = d.q.Enqueue(val)
	d.mutex.Unlock()

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
