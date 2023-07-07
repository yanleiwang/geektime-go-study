package tcp_study

import (
	"net"
)

type Client struct {
	network string
	addr    string
}

func NewClient(network string, addr string) *Client {
	return &Client{network: network, addr: addr}
}

func (c *Client) Send(data string) (string, error) {
	conn, err := net.Dial(c.network, c.addr)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = conn.Close()
	}()

	err = WriteMsg(conn, []byte(data))
	if err != nil {
		return "", err
	}
	msg, err := ReadMsg(conn)
	return string(msg), err
}
