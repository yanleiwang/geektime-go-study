package queue

import (
	"context"
	"errors"
	"sync"
)

type node[T any] struct {
	val  T
	next *node[T]
	//prev *node[T]  实现 队列 所以不需要prev
}

// LinkedQueue 并发, 不阻塞, 加锁, 通过链表实现的队列
type LinkedQueue[T any] struct {
	head  *node[T] // 永远是 哨兵头结点
	tail  *node[T] // 指向最后一个元素, 当tail == head的时候 表示节点个数为空
	zero  T
	mutex *sync.Mutex
}

func NewLinkedQueue[T any]() *LinkedQueue[T] {
	// 给他一个 哨兵头节点
	dummy := &node[T]{}
	return &LinkedQueue[T]{
		head:  dummy,
		tail:  dummy,
		mutex: &sync.Mutex{},
	}
}

func (l *LinkedQueue[T]) In(ctx context.Context, val T) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	newOne := &node[T]{
		val: val,
	}

	l.tail.next = newOne
	l.tail = newOne

	return nil
}

func (l *LinkedQueue[T]) Out(ctx context.Context) (T, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.head == l.tail {
		return l.zero, errors.New("没有新的元素了")
	}

	ret := l.head.next.val
	l.head = l.head.next // 最精妙的地方, 把 下一个节点 当作新的dummy 节点
	return ret, nil
}
