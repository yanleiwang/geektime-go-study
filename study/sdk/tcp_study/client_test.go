package tcp_study

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	go func() {
		srv := NewServer("tcp", ":8082")
		err := srv.Start()
		if err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(3 * time.Second)
	client := NewClient("tcp", "localhost:8082")
	for i := 0; i < 10; i++ {
		msg, err := client.Send("Hello")
		assert.NoError(t, err)
		assert.Equal(t, "HelloHello", msg)
	}

}
