package queue

import (
	"context"
	"time"
)

type Queue[T any] interface {
	In(ctx context.Context, val T) error // 入队
	Out(ctx context.Context) (T, error)  // 出队
}

/*
设计一个并发队列, 要考虑的因素

+  是否阻塞?
  + 队列空, 阻塞拿数据的 .   队列满了, 阻塞放数据的
+ 阻塞是否有超时控制?
  + 很多业务能接受阻塞一段时间，但是不能接受一直阻塞；
  +  带超时控制可以防止资源泄露：
+ 有界/无界
  + 是否支持 自动扩容/ 缩容,  还是固定容量?
+ 底层数据结构是  链表 or 数组?
+ 有锁 or 无锁
+ 公平性原则:  是不是先到先得(尤其是阻塞的情况下)

*/

// Comparator 用于比较两个对象的大小 src < dst, 返回-1，src = dst, 返回0，src > dst, 返回1
// 不要返回任何其它值！
type Comparator[T any] func(src T, dst T) int

type Delayable interface {
	Delay() time.Duration
}
