package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"time"
)

var (
	ErrRPCNotSupportNil            = errors.New("micro: rpc_study 不支持nil")
	ErrRPCOnlySupportPointOfStruct = errors.New("micro: rpc_study 只支持指向结构体的一级指针")
)

func InitClientProxy(addr string, service Service) error {
	c := NewClient(addr)
	return setFuncField(service, c)
}

func setFuncField(service Service, p Proxy) error {
	if service == nil {
		return ErrRPCNotSupportNil
	}

	val := reflect.ValueOf(service)
	typ := val.Type()
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return ErrRPCOnlySupportPointOfStruct
	}

	val = val.Elem()
	typ = typ.Elem()
	numField := val.NumField()
	for i := 0; i < numField; i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		// 捕获本地调用, 而后调用set篡改了他, 改而发起rpc调用
		// args[0] 是context
		// args[1] 是req
		fn := func(args []reflect.Value) (results []reflect.Value) {

			retVal := reflect.New(fieldTyp.Type.Out(0).Elem())

			// args[0] 是 context
			ctx := args[0].Interface().(context.Context)
			// args[1] 是 req
			reqData, err := json.Marshal(args[1].Interface())
			if err != nil {
				return []reflect.Value{retVal, reflect.ValueOf(err)}
			}
			req := &Request{
				ServiceName: service.Name(),
				MethodName:  fieldTyp.Name,
				Arg:         reqData,
			}

			resp, err := p.Invoke(ctx, req)
			if err != nil {
				return []reflect.Value{retVal, reflect.ValueOf(err)}
			}
			err = json.Unmarshal(resp.Data, retVal.Interface())
			if err != nil {
				return []reflect.Value{retVal, reflect.ValueOf(err)}
			}

			return []reflect.Value{retVal, reflect.Zero(reflect.TypeOf(new(error)).Elem())}

		}
		fnVal := reflect.MakeFunc(fieldVal.Type(), fn)
		fieldVal.Set(fnVal)

	}
	return nil
}

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Invoke(ctx context.Context, req *Request) (*Response, error) {
	data, err := c.serialize(req)
	if err != nil {
		return nil, err
	}

	respData, err := c.Send(data)
	if err != nil {
		return nil, err
	}

	return &Response{
		Data: respData,
	}, nil

}

func (c *Client) serialize(req *Request) ([]byte, error) {
	return json.Marshal(req)
}

func (c *Client) Send(data []byte) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", c.addr, time.Second*3)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	err = WriteMsg(conn, data)
	if err != nil {
		return nil, err
	}

	return ReadMsg(conn)
}
