package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

type Server struct {
	services map[string]reflectionStub
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]reflectionStub, 16),
	}
}

func (s *Server) RegisterService(service Service) {
	s.services[service.Name()] = reflectionStub{
		s:     service,
		value: reflect.ValueOf(service),
	}
}

func (s *Server) Start(network, addr string) error {
	listener, err := net.Listen(network, addr)
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
		reqBs, err := ReadMsg(conn)
		if err != nil {
			return err
		}

		// 还原调用信息
		req := &Request{}
		err = json.Unmarshal(reqBs, req)
		if err != nil {
			return err
		}

		resp, err := s.Invoke(context.Background(), req)
		if err != nil {
			// 这个可能你的业务 error
			// 暂时不知道怎么回传 error，所以我们简单记录一下
			return err
		}

		if err = WriteMsg(conn, resp.Data); err != nil {
			return err
		}
	}

}

func (s *Server) Invoke(ctx context.Context, req *Request) (*Response, error) {
	service, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("你要调用的服务不存在")
	}
	resp, err := service.invoke(ctx, req.MethodName, req.Arg)
	if err != nil {
		return nil, err
	}
	return &Response{
		Data: resp,
	}, nil
}

type reflectionStub struct {
	s     Service
	value reflect.Value
}

func (s *reflectionStub) invoke(ctx context.Context, methodName string, data []byte) ([]byte, error) {
	// 反射找到方法，并且执行调用
	method := s.value.MethodByName(methodName)
	in := make([]reflect.Value, 2)
	// 暂时我们不知道怎么传这个 context，所以我们就直接写死
	in[0] = reflect.ValueOf(context.Background())
	inReq := reflect.New(method.Type().In(1).Elem())
	err := json.Unmarshal(data, inReq.Interface())
	if err != nil {
		return nil, err
	}
	in[1] = inReq
	results := method.Call(in)
	// results[0] 是返回值
	// results[1] 是error
	if results[1].Interface() != nil {
		return nil, results[1].Interface().(error)
	}
	return json.Marshal(results[0].Interface())
}
