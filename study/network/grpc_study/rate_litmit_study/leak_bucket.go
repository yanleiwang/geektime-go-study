package rate_litmit_study

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// LeakBucket 漏桶
// 算法要点:
//
// + 请求过来先排队
// + 每隔一段时间，放过去一个请求
// + 请求排队直到通过，或者超时
//
// 令牌桶是每隔一段时间生成一个令牌,  并且令牌能够放进桶里存储起来.
//
// **而漏桶也是每隔一段时间生成一个令牌,   但是令牌是没有桶存储的,  或者说漏桶相当于容量为0的令牌桶.**   所以其实漏桶其实就是开一个定时器, 定时放一个请求过去即可.
type LeakBucket struct {
	producer *time.Ticker
}

func NewLeakBucket(duration time.Duration) *LeakBucket {
	return &LeakBucket{
		producer: time.NewTicker(duration),
	}
}

func (l *LeakBucket) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 拿到令牌
		select {
		// 超时
		case <-ctx.Done():
			err = ctx.Err()
			return
		// 拿到令牌
		case <-l.producer.C:
			resp, err = handler(ctx, req)
			return
		}
	}
}

func (l *LeakBucket) Close() error {
	l.producer.Stop()
	return nil
}
