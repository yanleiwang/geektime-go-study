package connection_pool

import (
	"fmt"
	"github.com/silenceper/pool"
	"net"
	"os"
	"testing"
	"time"
)

func server() error {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			for {
				conn.Write([]byte("Hello, World!"))
				time.Sleep(10 * time.Second)
			}
		}()

	}

}

func TestPool(t *testing.T) {
	go func() {
		err := server()
		if err != nil {
			t.Log(err)
			os.Exit(1)
		}
	}()
	time.Sleep(time.Second * 3)
	//factory 创建连接的方法
	factory := func() (interface{}, error) { return net.Dial("tcp", "127.0.0.1:8082") }

	//close 关闭连接的方法
	close := func(v interface{}) error { return v.(net.Conn).Close() }

	//ping 检测连接的方法
	//ping := func(v interface{}) error { return nil }

	//创建一个连接池: 初始化5, 最大空闲连接是20,  最大并发连接30
	poolConfig := &pool.Config{
		InitialCap: 5,  // 初始连接数
		MaxIdle:    20, // 最大空闲连接数
		MaxCap:     30, // 最大并发连接数
		Factory:    factory,
		Close:      close,
		//Ping:       ping,
		// 连接最大空闲时间, 超过该时间的连接 将会关闭, 可避免空闲时连接EOF, 自动失效的问题
		IdleTimeout: 15 * time.Second,
	}
	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		fmt.Println("err=", err)
	}

	// 从连接池里获取一个连接
	v, err := p.Get()

	//do something
	//conn=v.(net.Conn)

	//用完了 放回连接
	p.Put(v)

	//Release all connections in the connection pool, when resources need to be destroyed
	p.Release()

	// 获取当前连接池 连接数
	current := p.Len()
	t.Log(current)
}
