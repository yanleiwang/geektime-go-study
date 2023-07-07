package registry_study

import (
	"context"
	gen2 "geektime-go-study/study/network/grpc_study/proto/gen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"testing"
	"time"
)

func TestSimpleResolver(t *testing.T) {
	go func() {
		StartServer(t)
	}()

	time.Sleep(time.Second)

	//l, err := grpc.Dial("simple:///不管写啥都ok", grpc.WithInsecure())
	// simple:///不管写啥都ok   解析结果如下
	//  scheme: simple
	//  authority: 空  如果是 simple://123/不管写啥都ok, 那就是123
	// Endpoint: 不管写啥都ok
	conn, err := grpc.Dial("simple:///不管写啥都ok", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := gen2.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.GetById(ctx, &gen2.GetByIdReq{Id: 567})
	cancel()

	require.NoError(t, err)
	log.Println(resp)
}

func StartServer(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:8081")
	require.NoError(t, err)
	s := grpc.NewServer()
	us := &UserService{}
	gen2.RegisterUserServiceServer(s, us)

	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
