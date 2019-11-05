package rpcc

import (
	"context"
	"time"

	"github.com/fidelfly/gox/pkg/ctxx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Connection struct {
	*grpc.ClientConn
}

func NewConn(address string, timeout time.Duration, authKey string) (*Connection, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(clientInterceptor(timeout, authKey)))
	if err != nil {
		return nil, err
	}
	return &Connection{conn}, nil
}

func clientInterceptor(defaultTimeout time.Duration, authKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if len(authKey) > 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authKey)
		}
		timeout := defaultTimeout
		if md := ctxx.GetMetadata(ctx); md != nil {
			if v, ok := md.Get("rpc.timeout"); ok {
				if newTimeout, ok := v.(time.Duration); ok {
					timeout = newTimeout
				}
			}
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

type SimpleAuthKey string

func (sa SimpleAuthKey) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if len(sa) > 0 {
		return map[string]string{
			"authorization": string(sa),
		}, nil
	}
	return make(map[string]string, 0), nil
}

func (sa SimpleAuthKey) RequireTransportSecurity() bool {
	return true
}

func WithTimeout(ctx context.Context, timeout time.Duration) context.Context {
	return ctxx.WithMetadata(ctx, map[interface{}]interface{}{
		"rpc.timeout": timeout,
	})
	//return context.WithValue(ctx, "rpc.timeout", timeout)
}
