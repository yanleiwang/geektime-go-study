package channel

import (
	"context"
)

/*
例子：利用 channel 来实现一个任务池。该任务
池允许开发者提交任务，并且设定最多多少个
goroutine 同时运行。
*/

type Task func()

type TaskPool struct {
	tasks chan Task
	close chan struct{}
}

// numG 是 goroutine 数量
// numT 是缓存的容量
func NewTaskPool(numG int, numT int) *TaskPool {
	res := &TaskPool{
		tasks: make(chan Task, numT),
		close: make(chan struct{}),
	}

	for i := 0; i < numG; i++ {
		go func() {
			for {
				select {
				case <-res.close:
					return
				case t := <-res.tasks:
					t()
				}
			}

		}()
	}
	return res
}

// 提交任务
func (p *TaskPool) Do(ctx context.Context, t Task) error {
	//p.tasks <- t
	//return nil
	select {
	case p.tasks <- t:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (p *TaskPool) Close() error {
	close(p.close)
	return nil
}
