package registry_study

import (
	"context"
	"geektime-go-study/study/network/grpc_study/registry_study/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"time"
)

type grpcResolverBuilder struct {
	r       registry.Registry
	timeout time.Duration
}

func NewResolverBuilder(r registry.Registry, timeout time.Duration) *grpcResolverBuilder {
	return &grpcResolverBuilder{
		r:       r,
		timeout: timeout,
	}
}

func (g *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ret := &grpcResolver{
		target:  target,
		cc:      cc,
		timeout: g.timeout,
		r:       g.r,
		close:   make(chan struct{}),
	}
	ret.resolve()
	ret.watch()
	return ret, nil

}

func (g *grpcResolverBuilder) Scheme() string {
	return "registry"
}

type grpcResolver struct {
	target  resolver.Target
	cc      resolver.ClientConn
	timeout time.Duration
	r       registry.Registry
	close   chan struct{}
}

func (g *grpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	// 重新获取一下所有服务
	g.resolve()
}

func (g *grpcResolver) resolve() {
	serviceName := g.target.Endpoint()
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	instances, err := g.r.ListServices(ctx, serviceName)
	cancel()
	if err != nil {
		g.cc.ReportError(err)
	}
	address := make([]resolver.Address, 0, len(instances))
	for _, ins := range instances {
		address = append(address, resolver.Address{
			Addr: ins.Address,
			Attributes: attributes.New("weight", ins.Weight).
				WithValue("group", ins.Group),
		})
	}
	err = g.cc.UpdateState(resolver.State{
		Addresses: address,
	})
	if err != nil {
		g.cc.ReportError(err)
	}
}

func (g *grpcResolver) watch() {
	go func() {
		serviceName := g.target.Endpoint()
		events := g.r.Subscribe(serviceName)
		for {
			select {
			case <-events:
				// 一种做法就是我们这边区别处理不同事件类型，然后更新数据
				// switch event.Type {
				//
				//}
				// 另外一种做法就是我们这里采用的，每次事件发生的时候，就直接刷新整个可用服务列表
				g.resolve()
			case <-g.close:
				return
			}
		}

	}()
}

func (g *grpcResolver) Close() {
	// 需要考虑要不要防止 Close() 多次调用
	close(g.close)
	// 或者
	//g.close <- struct{}{}

}
