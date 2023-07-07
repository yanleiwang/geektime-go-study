package round_robin

import (
	"context"
	"fmt"
	"geektime-go-study/study/network/grpc_study/proto/gen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)

func TestBalancer_e2e_Pick(t *testing.T) {
	go StartServer(t)
	// 等服务器起来
	time.Sleep(time.Second * 3)

	// 不在init里面注册 在这里 应该也是可以的
	//balancer.Register(base.NewBalancerBuilder(Name, &builder{}, base.Config{HealthCheck: true}))
	cc, err := grpc.Dial("localhost:8081", grpc.WithInsecure(),
		//grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "demo_round_robin"}`))
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, Name)))
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
