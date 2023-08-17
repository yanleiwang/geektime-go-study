package zookeeper

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 在conn上注册 事件回调
func TestWatcherV1(t *testing.T) {
	cb := func(event zk.Event) {
		fmt.Println("###########################")
		fmt.Println("path: ", event.Path)
		fmt.Println("type: ", event.Type.String())
		fmt.Println("state: ", event.State.String())
		fmt.Println("---------------------------")
	}

	conn, _, err := zk.Connect(hosts, time.Second*3, zk.WithEventCallback(cb))
	assert.NoError(t, err)
	defer conn.Close()

	// 开启监听, 监听只能监听一次
	// 所以只能监听到create ， 不能监听到 delete
	exist, _, _, err := conn.ExistsW("/wyl")
	assert.NoError(t, err)
	assert.False(t, exist)

	_, err = conn.Create("/wyl", []byte("123"), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	assert.NoError(t, err)

	err = conn.Delete("/wyl", -1)
	assert.NoError(t, err)
}

// 单个函数 监听
func TestWatcherV2(t *testing.T) {
	conn, _, err := zk.Connect(hosts, time.Second*3)
	assert.NoError(t, err)
	defer conn.Close()

	// 开启监听, 监听只能监听一次
	// 所以只能监听到create ， 不能监听到 delete
	exist, _, eventChan, err := conn.ExistsW("/wyl")
	assert.NoError(t, err)
	assert.False(t, exist)

	go func() {
		event := <-eventChan
		fmt.Println("path: ", event.Path)
		fmt.Println("type: ", event.Type.String())
		fmt.Println("state: ", event.State.String())
	}()

	_, err = conn.Create("/wyl", []byte("123"), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	assert.NoError(t, err)

	err = conn.Delete("/wyl", -1)
	assert.NoError(t, err)
}
