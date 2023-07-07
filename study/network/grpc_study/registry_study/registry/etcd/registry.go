package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"geektime-go-study/study/network/grpc_study/registry_study/registry"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sync"
)

var typesMap = map[mvccpb.Event_EventType]registry.EventType{
	mvccpb.PUT:    registry.EventTypeAdd,
	mvccpb.DELETE: registry.EventTypeDelete,
}

type Registry struct {
	sess        *concurrency.Session
	client      *clientv3.Client
	mutex       sync.RWMutex
	watchCancel []func()
}

func NewRegistry(c *clientv3.Client) (*Registry, error) {
	// 使用 etcd的租约session 帮我们来进行 自动续约
	// 这里没有设置ttl, 所以默认是60s
	sess, err := concurrency.NewSession(c)
	if err != nil {
		return nil, err
	}
	return &Registry{
		sess:   sess,
		client: c,
	}, nil
}

func (r *Registry) Register(ctx context.Context, instance registry.ServiceInstance) error {
	val, err := json.Marshal(instance)
	if err != nil {
		return err
	}
	_, err = r.client.Put(ctx, r.instanceKey(instance), string(val), clientv3.WithLease(r.sess.Lease()))
	return err
}

func (r *Registry) UnRegister(ctx context.Context, instance registry.ServiceInstance) error {
	_, err := r.client.Delete(ctx, r.instanceKey(instance))
	return err
}

func (r *Registry) ListServices(ctx context.Context, name string) ([]registry.ServiceInstance, error) {
	resp, err := r.client.Get(ctx, r.serviceKey(name), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	res := make([]registry.ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var si registry.ServiceInstance
		err = json.Unmarshal(kv.Value, &si)
		if err != nil {
			return nil, err
		}
		res = append(res, si)
	}
	return res, nil
}

func (r *Registry) Subscribe(name string) <-chan registry.Event {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = clientv3.WithRequireLeader(ctx)
	r.mutex.Lock()
	r.watchCancel = append(r.watchCancel, cancel)
	r.mutex.Unlock()
	ch := r.client.Watch(ctx, r.serviceKey(name), clientv3.WithPrefix())
	// 这里有没有 Buffer 都无所谓
	res := make(chan registry.Event)
	// 监听变更
	go func() {
		for {
			select {
			case resp := <-ch:
				if resp.Err() != nil {
					continue
				}
				if resp.Canceled {
					return
				}

				res <- registry.Event{}
				// 更精细的控制, 但其实没什么必要
				//for _, event := range resp.Events {
				//	res <- registry.Event{
				//		Type: typesMap[event.Type],
				//	}
				//}
			case <-ctx.Done():
				return
			}
		}

	}()
	return res
}

func (r *Registry) Close() error {
	r.mutex.Lock()
	watchCancel := r.watchCancel
	r.watchCancel = nil
	r.mutex.Unlock()
	for _, cancel := range watchCancel {
		cancel()
	}
	// 因为 client 是外面传进来的，所以我们这里不能关掉它。它可能被其它的人使用着
	return r.sess.Close()
}

func (r *Registry) instanceKey(instance registry.ServiceInstance) string {
	return fmt.Sprintf("/micro/%s/%s", instance.Name, instance.Address)
}

func (r *Registry) serviceKey(name string) string {
	return fmt.Sprintf("/micro/%s", name)
}
