package queue

import (
	"context"
	_ "golang.org/x/sync/semaphore"
	"sync"
)

type SliceQueue[T any] struct {
	data     []T
	head     int
	tail     int
	size     int
	capacity int
	mutex    sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond

	zero T
}

type SliceQueueV1[T any] struct {
	data     []T
	head     int
	tail     int
	size     int
	capacity int
	mutex    sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
	zero     T
}

// InV1 没有超时控制
func (s *SliceQueueV1[T]) InV1(ctx context.Context, val T) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for s.isFull() {
		s.notFull.Wait()
	}

	s.data[s.tail] = val
	s.tail++
	s.size++
	if s.tail == s.capacity {
		s.tail = 0
	}
	s.notEmpty.Signal()
	return nil
}

// OutV1 没有超时控制
func (s *SliceQueueV1[T]) OutV1(ctx context.Context) (T, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for s.isEmpty() {
		s.notEmpty.Wait()
		select {
		case <-ctx.Done():
			return s.zero, ctx.Err()
		default:
			continue
		}
	}

	ret := s.data[s.head]
	s.data[s.head] = s.zero
	s.head++
	s.size--
	if s.head == s.capacity {
		s.head = 0
	}
	s.notFull.Broadcast()
	return ret, nil
}

func (s *SliceQueueV1[T]) isFull() bool {
	return s.size == s.capacity
}

func (s *SliceQueueV1[T]) isEmpty() bool {
	return s.size == 0
}

func NewSliceQueue[T any](capacity int) *SliceQueueV1[T] {

	ret := &SliceQueueV1[T]{
		data:     make([]T, capacity),
		capacity: capacity,
	}

	ret.notEmpty = sync.NewCond(&ret.mutex)
	ret.notFull = sync.NewCond(&ret.mutex)

	return ret
}
