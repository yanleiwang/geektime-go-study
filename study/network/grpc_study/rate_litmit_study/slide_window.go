package rate_litmit_study

import (
	"container/list"
	"context"
	"errors"
	"google.golang.org/grpc"
	"sync"
	"time"
)

// SlideWindowLimiter  滑动窗口
// 从当前时间开始，往前回溯一段时间，只能处理一定数量的请求。
// 滑动窗口的核心是：这个窗口永远是以当前时间戳为准，往前回溯。
// 例如当前时间点是 00:03:17 ，往前回溯一分钟，就是一个一分钟长度的窗口。
//
// 相较于固定窗口,  滑动窗口的效果更加平滑,  比如如果有大量的请求在固定窗口, 前后两个窗口之间进来,   对于服务器的压力是很大的.  而滑动窗口能够很好的处理这种情况.
//
// 滑动窗口实现要点:
//
// + 将执行请求的时间点  记录到队列里
// + 每次请求进来的时候,  从队尾淘汰不在当前窗口的记录
type SlideWindowLimiter struct {
	queue    *list.List // 请求时间 队列
	interval int64      // 窗口大小
	rate     int        // 当前窗口 允许的最大请求数
	mutex    sync.Mutex
}

func NewSlideWindowLimiter(interval time.Duration, rate int) *SlideWindowLimiter {
	return &SlideWindowLimiter{
		queue:    list.New(),
		interval: interval.Nanoseconds(),
		rate:     rate,
	}
}

func (t *SlideWindowLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		now := time.Now().UnixNano()
		boundary := now - t.interval

		// 快路径:  已经执行的请求数小于 最大请求数 直接允许执行
		t.mutex.Lock()
		length := t.queue.Len()
		if length < t.rate {
			resp, err = handler(ctx, req)
			// 记住了请求的时间戳
			t.queue.PushBack(now)
			t.mutex.Unlock()
			return
		}

		// 慢路径:  先尝试把所有不在窗口内的数据都删掉
		timestamp := t.queue.Front()
		// 这个循环把所有不在窗口内的数据都删掉了
		for timestamp != nil && timestamp.Value.(int64) < boundary {
			t.queue.Remove(timestamp)
			timestamp = t.queue.Front()
		}
		length = t.queue.Len()
		t.mutex.Unlock()
		if length >= t.rate {
			err = errors.New("到达瓶颈")
			return
		}
		resp, err = handler(ctx, req)
		// 记住了请求的时间戳
		t.queue.PushBack(now)
		return
	}
}
