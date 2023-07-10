package rate_litmit_study

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"time"
)

// TokenBucket 令牌桶
// 算法要点:
// + 令牌每隔一段时间t生成,  放入容量为capacity的桶内
// + 每次请求从桶里拿走一个令牌
// + 拿到令牌的请求就会被处理
// + 没有拿到令牌的请求可以被:
//   - 直接被拒绝
//   - 阻塞直到拿到令牌或者超时
//   - 降级处理
type TokenBucket struct {
	tokens    chan struct{}
	closeChan chan struct{}
}

func NewTokenBucket(capacity int, duration time.Duration) *TokenBucket {
	tokens := make(chan struct{}, capacity)
	closeChan := make(chan struct{})

	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				// 可以放令牌
				case tokens <- struct{}{}:
				// 没人取令牌
				default:
				}
			case <-closeChan:
				return
			}
		}
	}()

	return &TokenBucket{
		tokens:    tokens,
		closeChan: closeChan,
	}
}

func (t *TokenBucket) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 拿到令牌
		select {
		// 超时
		case <-ctx.Done():
			err = ctx.Err()
			return
		// 拿到令牌
		case <-t.tokens:
			resp, err = handler(ctx, req)
		// 令牌桶关闭了, 可以采取两种策略之一:
		// 1. 返回err
		// 2. 认为没有限流措施, 直接处理
		case <-t.closeChan:
			err = errors.New("没有限流措施, 拒绝请求")
			//resp, err = handler(ctx, req)
		}
		return
	}
}

func (t *TokenBucket) Close() error {
	close(t.closeChan)
	return nil
}
