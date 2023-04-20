package channel

import (
	"fmt"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	// 向一个已经关闭的channel 发送数据会panic
	c := make(chan struct{})
	close(c)
	c <- struct{}{}

	// 向一个nil channel 收发数据 会阻塞, 引起goroutine泄露
	var c1 chan struct{}
	c1 <- struct{}{}

}

func TestChannel_2(t *testing.T) {
	/*
		该代码会panic,  说明当 close channel 之前, 有goroutine阻塞在写,  然后再close, 那就会panic
		可以看代码 runtime/chan.go/chansend  最后几行
	*/

	c := make(chan int, 1)

	c <- 1
	go func() {
		c <- 2 // 会panic

	}()
	// 为了 阻塞在 c <- 2
	time.Sleep(time.Second)

	close(c)

	val, ok := <-c
	fmt.Println(val, ok) // 1 true

	val, ok = <-c
	fmt.Println(val, ok) // 0 false
	// 说明 c <- 2 没有成功

	time.Sleep(time.Second) // 保证 c <- 2 能执行

}
