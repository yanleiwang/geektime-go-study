package queue

import (
	"context"
	"errors"
	"sync/atomic"
	"unsafe"
)

type atomicNode[T any] struct {
	val  T
	next unsafe.Pointer // *atomicNode[T]
	//prev *node[T]  实现 队列 所以不需要prev
}

// LinkedQueueAtomic 用 原子操作(cas) 实现的 无锁并发队列, 基于链表
// 相较于 有锁版的,  cas操作性能更好.
// 但是在极高并发的情况下, 因为cas操作不会阻塞当前goroutine, 而是一直占着cpu, 所以其性能反而会下降
type LinkedQueueAtomic[T any] struct {
	head unsafe.Pointer // *atomicNode[T]
	tail unsafe.Pointer // *atomicNode[T]
	zero T
}

func (l *LinkedQueueAtomic[T]) In(ctx context.Context, val T) error {
	newOne := atomicNode[T]{
		val: val,
	}
	newPtr := unsafe.Pointer(&newOne)
	for {
		// 首先原子的获取 l.tail 的地址 以及 tailNext的地址
		tailPtr := atomic.LoadPointer(&l.tail)
		tail := (*atomicNode[T])(tailPtr)
		tailNext := atomic.LoadPointer(&tail.next)

		// 说明 有人已经执行完 第一个cas操作, 但可能还没执行第二个cas操作, 所以重新循环
		// 这里必须 要加这个判断
		if tailNext != nil {
			continue
		}

		// 首先 执行 l.tail.next = newPtr
		// 再执行 l.tail = newPtr
		// 必须按这个步骤!, 否则会有并发问题, 比如:
		// 当前没有节点
		// g1: 执行完入队第一个cas
		// g2: 执行出队操作, 此时会因为 headNextPtr == nil 而panic掉
		if atomic.CompareAndSwapPointer(&tail.next, tailNext, newPtr) {
			// 因为同时只会有一个 goroutine进入这里, 所以 肯定会成功?
			atomic.CompareAndSwapPointer(&l.tail, tailPtr, newPtr)
			return nil
		}

	}

}

func (l *LinkedQueueAtomic[T]) Out(ctx context.Context) (T, error) {
	for {
		headPtr := atomic.LoadPointer(&l.head)
		head := (*atomicNode[T])(headPtr)
		tailPtr := atomic.LoadPointer(&l.tail)
		tail := (*atomicNode[T])(tailPtr)
		if head == tail {
			// 不需要做更多检测，在当下这一刻，我们就认为没有元素，即便这时候正好有人入队
			// 但是并不妨碍我们在它彻底入队完成——即所有的指针都调整好——之前，
			// 认为其实还是没有元素
			return l.zero, errors.New("没有新的元素了")
		}
		headNextPtr := atomic.LoadPointer(&head.next)
		if atomic.CompareAndSwapPointer(&l.head, headPtr, headNextPtr) {
			headNext := (*node[T])(headNextPtr)
			return headNext.val, nil
		}
	}
}
