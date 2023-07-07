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
}
