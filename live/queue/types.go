package queue

import "context"

type Queue[T any] interface {
	In(ctx context.Context, val T) error
	Out(ctx context.Context) (T, error)
}

// Comparator 用于比较两个对象的大小 src < dst, 返回-1，src = dst, 返回0，src > dst, 返回1
// 不要返回任何其它值！
type Comparator[T any] func(src T, dst T) int
