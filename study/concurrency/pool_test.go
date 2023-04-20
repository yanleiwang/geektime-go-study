package concurrency

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

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
