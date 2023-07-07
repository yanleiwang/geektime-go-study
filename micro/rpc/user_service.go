package rpc

import (
	"context"
	"log"
)

type UserService struct {
	/*
		为什么用函数字段而不是方法? 因为 go的反射不支持 修改方法的实现,  所以只能曲线救国
		参数 context的作用:  1)传递元数据, 用以支持链路追踪,AB测试 等等. 2)超时控制
		参数 GetByIdReq 用来放 实际函数所需的参数.  为什么用结构体? 而不是 参数1, 参数2这种? 可以, 但是实现起来更复杂, 这里的实现只是为了 概览 rpc的实现.没必要
												为什么用结构体指针? 指针 不存在复制的问题(这里暂时还不太理解)
	*/
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

type GetByIdReq struct {
	Id int
}

func (u UserService) Name() string {
	return "user-service"
}

type GetByIdResp struct {
	Msg string
}

type UserServiceServer struct {
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	log.Println(req)
	return &GetByIdResp{
		Msg: "hello, world",
	}, nil
}

func (u *UserServiceServer) Name() string {
	return "user-service"
}
