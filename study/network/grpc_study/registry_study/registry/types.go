package registry

import (
	"context"
	"io"
)

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeAdd
	EventTypeDelete
)

type Event struct {
	Type     EventType
	Instance ServiceInstance
}

type Registry interface {
	// 服务注册
	Register(ctx context.Context, instance ServiceInstance) error
	UnRegister(ctx context.Context, instance ServiceInstance) error
	// 服务发现
	ListServices(ctx context.Context, name string) ([]ServiceInstance, error)
	Subscribe(name string) <-chan Event
	// close
	io.Closer
}

type ServiceInstance struct {
	Name    string
	Address string
	// 这边你可以任意加字段，完全取决于你的服务治理需要什么字段
	Weight uint32
	// 可以考虑再加一个分组字段
	Group string

	// 也可以用这个
	//Attributes map[string]string
}
