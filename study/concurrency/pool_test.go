package concurrency

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

/*
一般情况下，如果要考虑缓存资源，比如说创建好的对
象，那么可以使用 sync.Pool。
• sync.Pool 会先查看自己是否有资源，有则直接返回
• 没有则创建一个新的
• sync.Pool 会在 GC 的时候释放缓存的资源
*/

func TestPool(t *testing.T) {
	pool := sync.Pool{New: func() any {
		return new(bytes.Buffer)
	}}

	f := func() {
		b := pool.Get().(*bytes.Buffer)
		b.Reset()
		b.WriteString("hello")
		fmt.Println(b)
		pool.Put(b)
	}

	for i := 0; i < 10; i++ {
		f()
	}

}
