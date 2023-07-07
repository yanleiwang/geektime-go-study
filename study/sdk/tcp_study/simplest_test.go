package tcp_study

import (
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	go func() {
		err := Serve(":8082")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	err := Connect("localhost:8082")
	t.Log(err)
}
