package option

import (
	"context"
	"time"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpc_zap "github.com/niming-dev/ddd-demo/go-common/grpc/middleware/logging/zap"
	grpcopentelemetry "github.com/niming-dev/ddd-demo/go-common/grpc/middleware/opentelemetry"
	"github.com/niming-dev/ddd-demo/go-common/grpc/middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// 默认定义一些gRPC配置，方便使用
var (
	// DefaultKAEP gRPC长链接策略配置
	DefaultKAEP = keepalive.EnforcementPolicy{
		// If a client pings more than once every 5 seconds, terminate the connection
		MinTime: 5 * time.Second,
		// Allow pings even when there are no active streams
		PermitWithoutStream: true,
	}

	// DefaultKASP gRPC长链接策略配置
	DefaultKASP = keepalive.ServerParameters{
		// If a client is idle for 5 minute, send a GOAWAY
		MaxConnectionIdle: 5 * time.Minute,
		// If any connection is alive for more than 1 hour, send a GOAWAY
		MaxConnectionAge: 1 * time.Hour,
		// Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		MaxConnectionAgeGrace: 5 * time.Second,
		// Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Time: 5 * time.Second,
		// Wait 1 second for the ping ack before assuming the connection is dead
		Timeout: 1 * time.Second,
	}

	// DefaultKACP 客户端keepalive配置参数
	DefaultKACP = keepalive.ClientParameters{
		// send pings every 10 seconds if there is no activity
		Time: 10 * time.Second,
		// wait 1 second for ping ack before considering the connection dead
		Timeout: time.Second,
		// send pings even without active streams
		PermitWithoutStream: true,
	}

	// DefaultMaxRecvMsgSize 最大允许接收128MB的消息
	DefaultMaxRecvMsgSize = 1024 * 1024 * 128
	// DefaultMaxSendMsgSize 最大允许发送128MB的消息
	DefaultMaxSendMsgSize = 1024 * 1024 * 128

	// DefaultConnectionTimeout 默认建立链接超时时间
	DefaultConnectionTimeout = time.Second * 5

	// AlwaysLoggingDeciderServer 用于决定服务端是否需要记录请求或响应的payload, 默认所有都记录
	AlwaysLoggingDeciderServer = func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		return true
	}
	// AlwaysLoggingDeciderClient 用于决定客户端是否需要记录请求或响应的payload, 默认所有都记录
	AlwaysLoggingDeciderClient = func(ctx context.Context, fullMethodName string) bool {
		return true
	}
)

// GetServerOption 获取一组默认的gRPC服务端配置
func GetServerOption(options ...Option) []grpc.ServerOption {
	o := evaluateOpt(options)
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		// panic拦截器
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recovery.HandlerContext)),
		// 向ctx中设置一个空tags
		grpc_ctxtags.UnaryServerInterceptor(),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		// panic拦截器
		grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recovery.HandlerContext)),
		// 向ctx中设置一个空tags
		grpc_ctxtags.StreamServerInterceptor(),
	}

	if o.tracer != nil {
		// 把日志汇报进opentracing中
		streamInterceptors = append(streamInterceptors,
			grpcopentelemetry.StreamServerInterceptor(grpcopentelemetry.WithTracer(o.tracer)),
		)
		unaryInterceptors = append(unaryInterceptors,
			grpcopentelemetry.UnaryServerInterceptor(grpcopentelemetry.WithTracer(o.tracer)),
		)
	}

	if o.zapLogger != nil {
		unaryInterceptors = append(unaryInterceptors,
			// 日志记录，此处只记录状态与耗时
			grpc_zap.UnaryServerInterceptor(o.zapLogger),
			// 日志记录，记录请求与响应的payload
			grpc_zap.PayloadUnaryServerInterceptor(o.zapLogger, AlwaysLoggingDeciderServer),
		)

		streamInterceptors = append(streamInterceptors,
			// 日志记录，此处只记录状态与耗时
			grpc_zap.StreamServerInterceptor(o.zapLogger),
			// 日志记录，记录请求与响应的payload
			grpc_zap.PayloadStreamServerInterceptor(o.zapLogger, AlwaysLoggingDeciderServer),
		)
	}

	if o.logrusLogger != nil {
		logrusEntry := logrus.NewEntry(o.logrusLogger)
		unaryInterceptors = append(unaryInterceptors,
			// 日志记录，此处只记录状态与耗时
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			// 日志记录，记录请求与响应的payload
			grpc_logrus.PayloadUnaryServerInterceptor(logrusEntry, AlwaysLoggingDeciderServer),
		)

		streamInterceptors = append(streamInterceptors,
			// 日志记录，此处只记录状态与耗时
			grpc_logrus.StreamServerInterceptor(logrusEntry),
			// 日志记录，记录请求与响应的payload
			grpc_logrus.PayloadStreamServerInterceptor(logrusEntry, AlwaysLoggingDeciderServer),
		)
	}

	if o.useValidator {
		streamInterceptors = append(streamInterceptors, grpc_validator.StreamServerInterceptor())
		unaryInterceptors = append(unaryInterceptors, grpc_validator.UnaryServerInterceptor())
	}

	if len(o.promServerHistogramOptions) > 0 {
		grpc_prometheus.EnableHandlingTimeHistogram(o.promServerHistogramOptions...)
		streamInterceptors = append(streamInterceptors, grpc_prometheus.StreamServerInterceptor)
		unaryInterceptors = append(unaryInterceptors, grpc_prometheus.UnaryServerInterceptor)
	}

	if len(o.serverStreamInterceptors) > 0 {
		streamInterceptors = append(streamInterceptors, o.serverStreamInterceptors...)
	}

	if len(o.serverUnaryInterceptors) > 0 {
		unaryInterceptors = append(unaryInterceptors, o.serverUnaryInterceptors...)
	}

	return []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(DefaultKAEP),
		grpc.KeepaliveParams(DefaultKASP),
		grpc.MaxRecvMsgSize(DefaultMaxRecvMsgSize),
		grpc.MaxSendMsgSize(DefaultMaxSendMsgSize),
		grpc.ConnectionTimeout(DefaultConnectionTimeout),
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	}
}

// GetDialOption 获取一组默认的gRPC客户端配置
func GetDialOption(options ...Option) []grpc.DialOption {
	o := evaluateOpt(options)
	var unaryInterceptors []grpc.UnaryClientInterceptor
	var streamInterceptors []grpc.StreamClientInterceptor

	if o.tracer != nil {
		// 把日志汇报进opentracing中
		streamInterceptors = append(streamInterceptors,
			grpcopentelemetry.StreamClientInterceptor(grpcopentelemetry.WithTracer(o.tracer)),
		)
		unaryInterceptors = append(unaryInterceptors,
			grpcopentelemetry.UnaryClientInterceptor(grpcopentelemetry.WithTracer(o.tracer)),
		)
	}

	if o.zapLogger != nil {
		unaryInterceptors = append(unaryInterceptors,
			// 日志记录，此处只记录状态与耗时
			grpc_zap.UnaryClientInterceptor(o.zapLogger),
			// 日志记录，记录请求与响应的payload
			grpc_zap.PayloadUnaryClientInterceptor(o.zapLogger, AlwaysLoggingDeciderClient),
		)

		streamInterceptors = append(streamInterceptors,
			// 日志记录，此处只记录状态与耗时
			grpc_zap.StreamClientInterceptor(o.zapLogger),
			// 日志记录，记录请求与响应的payload
			grpc_zap.PayloadStreamClientInterceptor(o.zapLogger, AlwaysLoggingDeciderClient),
		)
	}

	if o.logrusLogger != nil {
		logrusEntry := logrus.NewEntry(o.logrusLogger)
		unaryInterceptors = append(unaryInterceptors,
			// 日志记录，此处只记录状态与耗时
			grpc_logrus.UnaryClientInterceptor(logrusEntry),
			// 日志记录，记录请求与响应的payload
			grpc_logrus.PayloadUnaryClientInterceptor(logrusEntry, AlwaysLoggingDeciderClient),
		)

		streamInterceptors = append(streamInterceptors,
			// 日志记录，此处只记录状态与耗时
			grpc_logrus.StreamClientInterceptor(logrusEntry),
			// 日志记录，记录请求与响应的payload
			grpc_logrus.PayloadStreamClientInterceptor(logrusEntry, AlwaysLoggingDeciderClient),
		)
	}

	if o.useValidator {
		unaryInterceptors = append(unaryInterceptors, grpc_validator.UnaryClientInterceptor())
	}

	if len(o.promClientHistogramOptions) > 0 {
		grpc_prometheus.EnableHandlingTimeHistogram(o.promClientHistogramOptions...)
		unaryInterceptors = append(unaryInterceptors, grpc_prometheus.UnaryClientInterceptor)
		streamInterceptors = append(streamInterceptors, grpc_prometheus.StreamClientInterceptor)
	}

	if len(o.clientStreamInterceptors) > 0 {
		streamInterceptors = append(streamInterceptors, o.clientStreamInterceptors...)
	}

	if len(o.clientUnaryInterceptors) > 0 {
		unaryInterceptors = append(unaryInterceptors, o.clientUnaryInterceptors...)
	}

	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(DefaultKACP),
		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
		grpc.WithChainStreamInterceptor(streamInterceptors...),
	}
}
