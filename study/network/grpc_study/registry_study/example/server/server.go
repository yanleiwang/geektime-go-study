package main

import (
	"fmt"
	"geektime-go-study/study/network/grpc_study/proto/gen"
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

	s := registry_study.NewServer("user-service", registry_study.WithRegister(r), registry_study.WithTimeout(time.Second*3))

	us := &registry_study.UserService{}
	// 我们将 UserService 什么样才算是初始化好的问题交给用户自己解决
	// 用户必须要在确认好 UserService 已经完全准备好之后才能启动并且注册
	gen.RegisterUserServiceServer(s, us) // 在grpc 服务端 注册userService
	fmt.Println("启动服务器")
	if err = s.Start(":8081"); err != nil {
		fmt.Println(err)
	}
	s.Close()

}
