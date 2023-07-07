package main

import (
	"context"
	"fmt"
	gen2 "geektime-go-study/study/network/grpc_study/proto/gen"
	"geektime-go-study/study/network/grpc_study/registry_study"
	"geektime-go-study/study/network/grpc_study/registry_study/registry/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})

	if err != nil {
		panic(err)
	}
	r, err := etcd.NewRegistry(etcdClient)
	if err != nil {
		panic(err)
	}
	client := registry_study.NewClient(registry_study.ClientWithRegistry(r, time.Second*3))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// 设置初始化连接的 ctx
	conn, err := client.Dial(ctx, "user-service")
	cancel()
	if err != nil {
		panic(err)
	}

	userClient := gen2.NewUserServiceClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	resp, err := userClient.GetById(ctx, &gen2.GetByIdReq{
		Id: 12,
	})
	if err != nil {
		panic(err)
	}
	cancel()
	fmt.Println(resp)
}
