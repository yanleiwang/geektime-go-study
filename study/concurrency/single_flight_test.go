package concurrency

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"testing"
	"time"
)

/*
singleflight 提供了函数重复调用 抑制机制.  保证了只有一个goroutine会执行key相当的函数

主要API
Do  同步调用 会等待调用的返回
DoChan  异步调用  立即返回channel 多用于需要处理超时的情况
Forget  让singleflight 忘记这个key, 不然后面所有这个key的调用都是返回之前的结果

*/

func TestSingleFlight(t *testing.T) {
	var g singleflight.Group
	var wg sync.WaitGroup

	fmt.Println("Do 测试----------------")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer g.Forget("key1")
			v, err, shared := g.Do("key1", func() (interface{}, error) {
				time.Sleep(time.Second)
				fmt.Printf("[%d] goroutine executed\n", i)
				return 1, nil
			})
			fmt.Printf("[%d] goroutine, val : %v, err: %v, shared: %v \n", i, v, err, shared)
		}(i)
	}
	wg.Wait()

	fmt.Println("DoChan 测试----------------")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer g.Forget("key1")
			c := g.DoChan("key1", func() (interface{}, error) {
				time.Sleep(2 * time.Second)
				fmt.Printf("[%d] goroutine executed\n", i)
				return 1, nil
			})

			select {
			case res := <-c:
				fmt.Println("未超时", res)
			case <-time.After(time.Second):
				fmt.Println("超时了")
			}

		}(i)
	}
	wg.Wait()

}
