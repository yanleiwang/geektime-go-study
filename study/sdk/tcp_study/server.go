package tcp_study

import "net"

type Server struct {
	network string
	addr    string
}

func NewServer(network, addr string) *Server {
	return &Server{
		network: network,
		addr:    addr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen(s.network, s.addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if er := s.handleConn(conn); er != nil {
				_ = conn.Close()
			}
		}()

	}

}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		readMsg, err := ReadMsg(conn)
		if err != nil {
			return err
		}

		res, err := s.handleMsg(readMsg)
		if err != nil {
			return err
		}

		if err = WriteMsg(conn, res); err != nil {
			return err
		}
	}

}

func (s *Server) handleMsg(msg []byte) ([]byte, error) {
	res := make([]byte, len(msg)*2)
	copy(res[:len(msg)], msg)
	copy(res[len(msg):], msg)
	return res, nil

}
