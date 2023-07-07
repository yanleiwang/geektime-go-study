package registry_study

import (
	"context"
	"geektime-go-study/study/network/grpc_study/registry_study/registry"
	"google.golang.org/grpc"
	"net"
	"time"
)

// Server 封装grpc Server 和 服务注册
// 当调用Start方法的时候 会去注册中心 注册服务, 并且启动serve
type Server struct {
	name     string
	register registry.Registry
	si       registry.ServiceInstance
	listener net.Listener
	// 单个操作的超时时间，一般用于和注册中心打交道
	timeout time.Duration

	*grpc.Server
}

type Option func(s *Server)

func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithRegister(register registry.Registry) Option {
	return func(s *Server) {
		s.register = register
	}
}

func NewServer(name string, opts ...Option) *Server {
	s := &Server{
		name:    name,
		timeout: 1 * time.Second,
		Server:  grpc.NewServer(),
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = l
	if s.register != nil {
		s.si = registry.ServiceInstance{
			Name:    s.name,
			Address: l.Addr().String(), // 如果是部署在容器内, 要替换成容器内的服务名
		}
		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
		err = s.register.Register(ctx, s.si)
		cancel()
		if err != nil {
			return err
		}
	}
	return s.Server.Serve(l)
}

func (s *Server) Close() error {

	if s.register != nil {
		err := s.register.Close()
		if err != nil {
			return err
		}
	}

	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			return err
		}
	}
	return nil

}
