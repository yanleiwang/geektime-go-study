package concurrency

import (
	"sync/atomic"
	"testing"
)

// 整型
func TestAtomicInt(t *testing.T) {
	i := int32(0)
	atomic.AddInt32(&i, 1)
	j := atomic.LoadInt32(&i)
	println(i, j)
}

// 任意类型
func TestAtomicAny(t *testing.T) {

	var val atomic.Value
	type Config struct {
		data string
	}

	c := Config{data: "Hello"}
	val.Store(c)
	get_c := val.Load().(Config)
	println(get_c.data)

}
