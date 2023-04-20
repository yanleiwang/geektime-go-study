package concurrency

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCond(t *testing.T) {
	// 初始化
	// sync.Cond实例的初始化
	// 需要一个满足实现了sync.Locker接口的类型实例，
	// 通常我们使用sync.Mutex。
	// 条件变量需要这个互斥锁来同步临界区，保护用作条件的数据。
	cond := sync.NewCond(&sync.Mutex{})
	ready := false
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		defer wg.Done()
		// 等待 run == true
		// 用Lock 保护共享变量
		// wait 进入等待状态
		cond.L.Lock()
		for !ready {
			cond.Wait()
		}

		fmt.Println("run 2")
		cond.L.Unlock()
	}()

	time.Sleep(2 * time.Second)
	cond.L.Lock()
	ready = true
	fmt.Println("run 1")
	cond.Broadcast() // 广播 唤醒所有等待的goroutine 也可以唤醒一个  cond.Signal()
	cond.L.Unlock()

	wg.Wait()

}
