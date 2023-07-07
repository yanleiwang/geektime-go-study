package concurrency

import (
	"sync"
	"sync/atomic"
	"testing"
)

/*
WaitGroup 是用于同步多个 goroutine 之间工作
的。
• 要在开启 goroutine 之前先加1 -> wg.Add(1)
• 每一个小任务完成就减1 ->wg.Done()
• 调用 Wait 方法来等待所有子任务完成

常见场景是我们会把任务拆分给多个 goroutine 并
行完成。在完成之后需要合并这些任务的结果，或
者需要等到所有小任务都完成之后才能进入下一
步。
*/

func TestWaitGroup(t *testing.T) {
	wg := sync.WaitGroup{}
	var result int64 = 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(delta int) {
			defer wg.Done()
			atomic.AddInt64(&result, int64(delta))
		}(i)
	}
	wg.Wait()
	t.Log(result)
}
