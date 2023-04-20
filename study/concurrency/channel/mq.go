package channel

import (
	"errors"
	"sync"
)

/*
利用 channel 来实现一个基于内存的消息队列，并且有消费组的概念
*/

// 方案一：每一个消费者订阅的时候，创建一个子
// channel
type Broker struct {
	mutex sync.RWMutex
	chans []chan Msg
	//chans map[string][]chan Msg  // key是 topic
}

type Msg struct {
	//Topic string
	Content string
}

func (b *Broker) Send(m Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, c := range b.chans {
		select {
		case c <- m:
		default:
			return errors.New("消息队列已满")
		}
	}

	return nil
}

func (b *Broker) Subscribe(capacity int) (<-chan Msg, error) {
	c := make(chan Msg, capacity)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.chans = append(b.chans, c)
	return c, nil
}

func (b *Broker) Close() error {
	b.mutex.Lock()
	chans := b.chans
	b.chans = nil
	b.mutex.Unlock()
	for _, msgs := range chans {
		close(msgs)
	}
	return nil
}

// 方案二：轮询所有的消费者
type BrokerV2 struct {
	mutex     sync.RWMutex
	consumers []func(Msg)
}

func (b *BrokerV2) Send(m Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, c := range b.consumers {
		c(m)
	}
	return nil
}

func (b *BrokerV2) Subscribe(f func(Msg)) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.consumers = append(b.consumers, f)
	return nil
}
