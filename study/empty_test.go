package study

import (
	"context"
	"testing"
)

func TestEmpty(t *testing.T) {
	var a any
	a = 1
	if a == 1 {
		println("ok")
	}

	a = "123"
	if a == "123" {
		println("string ok")
	}

}

type Person struct {
	name string
	age  int
}

func TestRange(t *testing.T) {
	ctx := context.Background()
	<-ctx.Done()

}
