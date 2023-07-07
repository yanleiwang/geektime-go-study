package example

import (
	"context"
	"geektime-go-study/study/network/grpc_study/proto/gen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	go StartServer(t)
	// 等服务器起来
	time.Sleep(time.Second * 3)

	// step1 创建 grpc 连接
	// 启用负载均衡 其实就是 在 Dial的时候传入ServiceConfig, 在其中指定loadBalance的name
	// 这里 使用的是 grpc 自带的 round_robin
	cc, err := grpc.Dial(":8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))
	require.NoError(t, err)
	defer cc.Close()

	// step2 创建 grpc 客户端
	c := gen.NewGreeterClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// step3 调用rpc方法
	resp, err := c.SayHello(ctx, &gen.HelloRequest{Name: "world"})
	require.NoError(t, err)
	t.Log(resp)

}

func StartServer(t *testing.T) {
	l, err := net.Listen("tcp", ":8081")
	require.NoError(t, err)
	// 创建grpc服务
	s := grpc.NewServer()

	// 注册 GreeterServer 到 grpc 服务上
	gen.RegisterGreeterServer(s, &server{})
	// 让grpc 监听 连接
	if err := s.Serve(l); err != nil {
		t.Log(err)
	}
}

type server struct {
	gen.UnimplementedGreeterServer
} //服务对象

// SayHello 实现服务的接口 在proto中定义的所有服务都是接口
func (s *server) SayHello(ctx context.Context, in *gen.HelloRequest) (*gen.HelloReply, error) {
	return &gen.HelloReply{Message: "Hello " + in.Name}, nil
}
