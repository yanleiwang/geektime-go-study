package registry_study

import "google.golang.org/grpc/resolver"

/*
实现最简单的resolver

grpc 服务发现原理:
根据grpc.Dial的 target 里的 scheme 找到对应的ResolverBuilder, 执行Build, 构建Resolver
然后
*/

func init() {
	// Register the SimpleResolverBuilder.
	// 通常写在init
	resolver.Register(&SimpleResolverBuilder{})
}

// 实现接口 resolver.Builder
type SimpleResolverBuilder struct{}

func (e *SimpleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ret := &SimpleResolver{
		cc: cc,
	}
	// 因为 clientConn 不会主动去调ResolveNow, 所以一开始 Resolver是没数据的, 需要我们手动去初始化数据
	ret.ResolveNow(resolver.ResolveNowOptions{})
	return ret, nil

}

// Scheme 返回一个固定的值
func (e *SimpleResolverBuilder) Scheme() string {
	return "simple"
}

// 实现接口 Resolver
type SimpleResolver struct {
	cc resolver.ClientConn
}

// ResolveNow 在该函数 中 获取注册中心数据 更新服务端地址
func (e *SimpleResolver) ResolveNow(options resolver.ResolveNowOptions) {
	// 只是示例, 所以固定写死IP和端口
	err := e.cc.UpdateState(resolver.State{
		Addresses: []resolver.Address{
			{
				Addr: "localhost:8081",
			},
		},
	})
	if err != nil {
		// ReportError 会重新调用ResolveNow
		e.cc.ReportError(err)
	}

}

func (e *SimpleResolver) Close() {
	// 释放资源
	// 因为SimpleResolver 没干啥, 所以这里什么也不做
}
