package grpc

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/keepalive"

	"github.com/niming-dev/ddd-demo/go-common/log"
)

var (
	defaultDialOpts       []grpc.DialOption
	defaultRequestTimeout time.Duration
)

func init() {
	defaultDialOpts = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(false),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1.0 * time.Second,
				Multiplier: 1.0,
				Jitter:     0.2,
				MaxDelay:   3 * time.Second,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			requestUnaryClientInterceptor,
			loggingUnaryClientInterceptor,
		)),
	}
	defaultRequestTimeout = 20 * time.Second
}

type ClientConn struct {
	conn *grpc.ClientConn
}

func (a *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	return a.conn.Invoke(ctx, method, args, reply, opts...)
}

func (a *ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return a.conn.NewStream(ctx, desc, method, opts...)
}

func requestUnaryClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultRequestTimeout)
		defer cancel()
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

func loggingUnaryClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	logger := log.WithFields(ctx, log.Fields{
		log.FieldKeyModule:  "GRPC",
		"method":            method,
		"request":           req,
		log.FieldKeyElapsed: time.Since(start).String(),
	})
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info("success")
	}
	return err
}

func DialContext(ctx context.Context, dsn string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) {
	opts = append(defaultDialOpts, opts...)
	conn, err := grpc.DialContext(ctx, dsn, opts...)
	if err != nil {
		return conn, err
	}
	return &ClientConn{conn}, nil
}
