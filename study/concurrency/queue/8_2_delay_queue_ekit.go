// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
该实现 主要是用 close(channel)  来起到通知所有入队/出队的 goroutine的作用
*/

// DelayQueueEkit 延时队列
// 每次出队的元素必然都是已经到期的元素，即 Delay() 返回的值小于等于 0
// 延时队列本身对时间的精确度并不是很高，其时间精确度主要取决于 time.Timer
// 所以如果你需要极度精确的延时队列，那么这个结构并不太适合你。
// 但是如果你能够容忍至多在毫秒级的误差，那么这个结构还是可以使用的
type DelayQueueEkit[T Delayable] struct {
	q             PriorityQueue[T]
	mutex         *sync.Mutex
	dequeueSignal *cond
	enqueueSignal *cond
}

func NewDelayQueueEkit[T Delayable](c int) *DelayQueueEkit[T] {
	m := &sync.Mutex{}
	res := &DelayQueueEkit[T]{
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
		mutex:         m,
		dequeueSignal: newCond(),
		enqueueSignal: newCond(),
	}
	return res
}

func (d *DelayQueueEkit[T]) Enqueue(ctx context.Context, t T) error {
	for {
		select {
		// 先检测 ctx 有没有过期
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		d.mutex.Lock()
		err := d.q.Enqueue(t)
		switch err {
		case nil:
			// 入队成功, 直接通知所有等待出队的goroutine 重新等待
			// 也可以对其进行优化 也就是 只有在新入队的元素 超时时间小于 之前的队首元素才 通知
			d.enqueueSignal.broadcast()
			d.mutex.Unlock()
			return nil
		case ErrOutOfCapacity:
			signal := d.dequeueSignal.signalCh()
			d.mutex.Unlock()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-signal:
			}
		default:
			d.mutex.Unlock()
			return fmt.Errorf("ekit: 延时队列入队的时候遇到未知错误 %w，请上报", err)
		}
	}
}

func (d *DelayQueueEkit[T]) Dequeue(ctx context.Context) (T, error) {
	var timer *time.Timer
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()
	for {
		select {
		// 先检测 ctx 有没有过期
		case <-ctx.Done():
			var t T
			return t, ctx.Err()
		default:
		}
		d.mutex.Lock()
		val, err := d.q.Peek()
		switch err {
		case nil:
			delay := val.Delay()
			if delay <= 0 {
				val, err = d.q.Dequeue()
				d.dequeueSignal.broadcast()
				d.mutex.Unlock()
				// 理论上来说这里 err 不可能不为 nil
				return val, err
			}
			signal := d.enqueueSignal.signalCh()
			d.mutex.Unlock()
			if timer == nil {
				timer = time.NewTimer(delay)
			} else {
				timer.Stop()
				timer.Reset(delay)
			}
			select {
			case <-ctx.Done():
				var t T
				return t, ctx.Err()
			case <-timer.C:
				// 到了时间, 直接continue
				continue
			case <-signal:
				// 表示有新的元素入队了,  直接continue
				continue

			}
		case ErrEmptyQueue:
			signal := d.enqueueSignal.signalCh()
			d.mutex.Unlock()
			select {
			case <-ctx.Done():
				var t T
				return t, ctx.Err()
			case <-signal:
			}
		default:
			d.mutex.Unlock()
			var t T
			return t, fmt.Errorf("ekit: 延时队列出队的时候遇到未知错误 %w，请上报", err)
		}
	}
}

type cond struct {
	signal chan struct{}
}

func newCond() *cond {
	return &cond{
		signal: make(chan struct{}),
	}
}

// broadcast 唤醒等待者
// 如果没有人等待，那么什么也不会发生
// 必须加锁之后才能调用这个方法
func (c *cond) broadcast() {
	signal := make(chan struct{})
	old := c.signal
	c.signal = signal
	//c.l.Unlock()
	close(old)
}

// signalCh 返回一个 channel，用于监听广播信号
// 必须在锁范围内使用
func (c *cond) signalCh() <-chan struct{} {
	res := c.signal
	//c.l.Unlock()
	return res
}
