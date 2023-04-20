package concurrency

import (
	"sync"
	"testing"
)

var once sync.Once

func printOne(id int) {
	once.Do(func() {
		println("only once")
	})
	println(id)

}

func TestOnce(t *testing.T) {
	for i := 0; i < 10; i++ {
		printOne(i)
	}

}
