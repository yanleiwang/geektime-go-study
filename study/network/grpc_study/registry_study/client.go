package registry_study

import (
	"context"
	"fmt"
	"geektime-go-study/study/network/grpc_study/registry_study/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"time"
)

type Client struct {
	rb resolver.Builder
}

type ClientOption func(*Client)

func NewClient(opts ...ClientOption) *Client {
	ret := &Client{}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func ClientWithRegistry(r registry.Registry, timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.rb = NewResolverBuilder(r, timeout)
	}
}

func (c *Client) Dial(ctx context.Context, serviceName string) (*grpc.ClientConn, error) {
	schema := ""
	if c.rb != nil {
		schema = c.rb.Scheme()
	}
	address := fmt.Sprintf("%s:///%s", schema, serviceName)
	opts := []grpc.DialOption{grpc.WithResolvers(c.rb), grpc.WithTransportCredentials(insecure.NewCredentials())}
	return grpc.DialContext(ctx, address, opts...)

}
