/*
最简单的 网络通信流程
*/
package tcp_study

import (
	"fmt"
	"io"
	"net"
	"time"
)

func Serve(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			handleConn(conn)
		}()

	}

}

/*
处理连接基本上就是在一个 for 循环内：
• 先读数据：读数据要根据上层协议来决定怎
么读。例如，简单的 RPC 协议一般是分成
两段读，先读头部，根据头部得知 Body 有
多长，再把剩下的数据读出来。
• 处理数据
• 回写响应：即便处理数据出错，也要返回一
个错误给客户端，不然客户端不知道你处理
出错了
*/
func handleConn(conn net.Conn) {
	for {
		bs := make([]byte, 5)
		_, err := conn.Read(bs)

		// 读写的时候, 可能遇到了错误,
		// 可以区分错误是 可以处理的 或者不可处理的
		//但是 建议出错了 就直接关闭, 这样代码比较简单
		if err != nil {
			_ = conn.Close()
			return
		}

		//// 表示连接关闭
		//if err == io.EOF || err == net.ErrClosed || err == io.ErrUnexpectedEOF {
		//	// 一般关闭的错误懒得管
		//	// 也可以把关闭错误输出到日志
		//	_ = conn.Close()
		//	return
		//}
		//// 其他错误是可以忽略的错误
		//if err != nil {
		//	continue
		//}

		res, err := handleMsg(bs)
		if err != nil {
			res = []byte("handleMsg err") //即便处理数据出错，也要返回一 个错误给客户端，不然客户端不知道你处理 出错了
		}
		_, err = conn.Write(res)
		if err == io.EOF || err == net.ErrClosed || err == io.ErrUnexpectedEOF {
			_ = conn.Close()
			return
		}
	}
}

func handleMsg(data []byte) ([]byte, error) {
	res := make([]byte, 2*len(data))
	copy(res[:len(data)], data)
	copy(res[len(data):], data)
	return res, nil
}

func Connect(addr string) error {

	conn, err := net.DialTimeout("tcp", addr, time.Second*3)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		_, err := conn.Write([]byte("hello"))
		if err != nil {
			return err
		}

		bs := make([]byte, 128)
		_, err = conn.Read(bs)
		if err != nil {
			return err
		}
		fmt.Println(string(bs))
	}
	return nil
}
