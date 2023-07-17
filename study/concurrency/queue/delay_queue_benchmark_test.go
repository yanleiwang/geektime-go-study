package queue

import (
	"context"
	"sync"
	"testing"
	"time"
)

type DelayQueue[T Delayable] interface {
	Enqueue(ctx context.Context, t T) error
	Dequeue(ctx context.Context) (T, error)
}

func BenchmarkDelayQueueEnqueue(b *testing.B) {
	f := func(b *testing.B, q DelayQueue[delayElem]) {
		e := delayElem{
			//10秒后过期
			deadline: time.Now().Add(time.Second * 10),
			val:      3,
		}
		ctx := context.Background()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				q.Enqueue(ctx, e)
			}
		})
	}

	b.Run("入队 sema实现", func(b *testing.B) {
		q := NewDelayQueueUseSema[delayElem](b.N)
		f(b, q)
	})

	b.Run("入队 ekit实现", func(b *testing.B) {
		q := NewDelayQueueEkit[delayElem](b.N)
		f(b, q)
	})
}

func BenchmarkDelayQueueDequeue(b *testing.B) {

	f := func(b *testing.B, q DelayQueue[delayElem]) {
		e := delayElem{
			//1秒后过期
			deadline: time.Now().Add(time.Millisecond * 10),
			val:      3,
		}
		ctx := context.Background()
		for i := 0; i < b.N; i++ {
			q.Enqueue(ctx, e)
		}

		// 全部过期
		time.Sleep(time.Millisecond * 100)

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				q.Dequeue(ctx)
			}
		})
	}

	b.Run("出队 sema实现", func(b *testing.B) {
		q := NewDelayQueueUseSema[delayElem](b.N)
		f(b, q)
	})

	b.Run("出队 ekit实现", func(b *testing.B) {
		q := NewDelayQueueEkit[delayElem](b.N)
		f(b, q)
	})

}

func BenchmarkDelayQueue(b *testing.B) {
	f := func(b *testing.B, q DelayQueue[delayElem]) {
		var wg sync.WaitGroup
		wg.Add(20)
		//10个 goroutine 一直入队
		for i := 0; i < 10; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < b.N; j++ {
					e := delayElem{
						//1秒后过期
						deadline: time.Now().Add(time.Millisecond * 1),
						val:      3,
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
					q.Enqueue(ctx, e)
					cancel()
				}
			}()
		}
		//10个 goroutine 一直出队
		for i := 0; i < 10; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < b.N; j++ {
					ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
					q.Dequeue(ctx)
					cancel()
				}
			}()
		}
		wg.Wait()
	}

	b.Run("同时出入队 sema实现", func(b *testing.B) {
		q := NewDelayQueueUseSema[delayElem](1000)
		f(b, q)
	})

	b.Run("同时出入队 ekit实现", func(b *testing.B) {
		q := NewDelayQueueEkit[delayElem](1000)
		f(b, q)
	})

}
