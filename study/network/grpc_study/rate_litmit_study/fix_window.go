package rate_litmit_study

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
	"time"
)

// FixWindowLimiter 固定窗口
// 算法要点:
//
// + 一个时间段内(窗口), 只能执行n个请求
// + 窗口是固定的,  一个窗口结束了, 才能够启动新的窗口,  不会滑动
//
// 实际实现的时候,  没有必要开一个定时器去刷新窗口,  而是实现一种类似的效果:
//
// + 记录当前窗口开始的时间t1 和窗口大小 duration
// + 请求进来的时候,  假如当前时间t2  >  t1+ duration. 表明需要开始新的窗口
type FixWindowLimiter struct {
	timestamp int64 // 窗口起始时间
	duration  int64 // 窗口大小
	rate      int64 // 在这个窗口内允许通过的最大请求数
	cnt       int64
	mutex     sync.Mutex
}

func NewFixWindowLimiter(duration time.Duration, rate int64) *FixWindowLimiter {
	return &FixWindowLimiter{
		timestamp: time.Now().UnixNano(),
		duration:  duration.Nanoseconds(),
		rate:      rate,
		cnt:       0,
		mutex:     sync.Mutex{},
	}
}

// BuildServerInterceptor_Atomic 原子操作版
func (f *FixWindowLimiter) BuildServerInterceptor_Atomic() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		cur := time.Now().UnixNano()
		timestamp := atomic.LoadInt64(&f.timestamp)
		cnt := atomic.LoadInt64(&f.cnt)
		// 需求开启新的窗口
		if timestamp+f.duration < cur {
			if atomic.CompareAndSwapInt64(&f.timestamp, timestamp, cur) {
				// 可以用store 而不必cas,  因为同时只会有一个go程进入
				atomic.StoreInt64(&f.cnt, 0)
				//atomic.CompareAndSwapInt64(&f.cnt, cnt, 0)
			}
		}

		cnt = atomic.AddInt64(&f.cnt, 1)
		if cnt > f.rate {
			err = errors.New("触发瓶颈了")
			return
		}
		resp, err = handler(ctx, req)
		return

	}

}

// BuildServerInterceptor_Mutex  加锁版
func (f *FixWindowLimiter) BuildServerInterceptor_Mutex() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		f.mutex.Lock()

		// 需要开启新的窗口
		cur := time.Now().UnixNano()
		if f.timestamp+f.duration < cur {
			f.cnt = 0
			f.timestamp = cur
		}
		// 请求数满了
		if f.cnt >= f.rate {
			f.mutex.Unlock()
			err = errors.New("触发瓶颈了")
			return
		}
		f.cnt++
		f.mutex.Unlock()
		resp, err = handler(ctx, req)
		return

	}

}
