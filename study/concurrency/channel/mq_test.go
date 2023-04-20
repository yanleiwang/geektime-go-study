package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBroker(t *testing.T) {
	broker := &Broker{}
	go func() {
		for {
			err := broker.Send(Msg{
				Content: time.Now().String(),
			})
			if err != nil {
				t.Log(err)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}

	}()

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		name := fmt.Sprintf("消费者 %d", i)
		go func() {
			defer wg.Done()
			msgs, err := broker.Subscribe(10)
			if err != nil {
				t.Log(err)
				return
			}

			for msg := range msgs {
				fmt.Println(name, msg.Content)
			}
		}()
	}
	wg.Wait()
}
