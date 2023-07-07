package lru_study

import (
	"fmt"
	"github.com/hashicorp/golang-lru/v2"
	"testing"
)

func TestGolangLru(t *testing.T) {
	l, _ := lru.New[int, interface{}](128)
	for i := 0; i < 256; i++ {
		l.Add(i, nil)
	}
	if l.Len() != 128 {
		panic(fmt.Sprintf("bad len: %v", l.Len()))
	}
}
