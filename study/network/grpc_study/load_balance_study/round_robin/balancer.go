package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
	"sync/atomic"
)

const Name = "demo_round_robin" // 跟grpc的round_robin 区别开来

// 注册 负载均衡器
// balancer.Register(b) 说 只能在init的时候调用, 但是我看应该在grpc.Dial之前也能用.
// Name用来区别不同负载均衡器
func init() {
	b := base.NewBalancerBuilder(Name, &builder{}, base.Config{
		HealthCheck: true,
	})
	balancer.Register(b)
}

// 实现接口 PickerBuilder
type builder struct {
}

func (b *builder) Build(info base.PickerBuildInfo) balancer.Picker {
	// 如果 可用节点为空, 就直接返回 永远返回err的 ErrPicker
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	subConns := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for k := range info.ReadySCs {
		subConns = append(subConns, k)
	}
	return &Picker{
		// 从随机下标开始 因为resolver 每次updateState的时候,
		// 都会重新调用Builder 创建新的Picker,
		//  所以为了防止如果节点频繁的上下线, 每次都是第一个节点被选中
		idx:      uint32(rand.Int31n(int32(len(info.ReadySCs)))),
		subConns: subConns,
	}

}

// 实现接口 Picker
type picker struct {
	idx      uint32
	subConns []balancer.SubConn
}

func (p *picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	// 不考虑 int溢出的问题,  因为溢出也就是从最大值回到了0
	// 造成的影响只是 某些节点少被选中了一轮
	idx := atomic.AddUint32(&p.idx, 1)
	sc := p.subConns[idx%uint32(len(p.subConns))]
	return balancer.PickResult{
		SubConn: sc,
	}, nil

}
