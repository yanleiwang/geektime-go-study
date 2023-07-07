package main

import (
	"context"
	"geektime-go-study/study/network/grpc_study/proto/gen"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	gen.UnimplementedGreeterServer
} //服务对象

// SayHello 实现服务的接口 在proto中定义的所有服务都是接口
func (s *server) SayHello(ctx context.Context, in *gen.HelloRequest) (*gen.HelloReply, error) {
	return &gen.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 创建grpc服务
	s := grpc.NewServer()

	// 注册 GreeterServer 到 grpc 服务上
	gen.RegisterGreeterServer(s, &server{})

	// 让grpc 监听 连接
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
