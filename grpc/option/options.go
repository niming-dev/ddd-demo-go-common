package option

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type options struct {
	zapLogger                  *zap.Logger
	logrusLogger               *logrus.Logger
	tracer                     trace.Tracer
	useValidator               bool
	promServerHistogramOptions []grpc_prometheus.HistogramOption
	promClientHistogramOptions []grpc_prometheus.HistogramOption

	serverStreamInterceptors []grpc.StreamServerInterceptor
	serverUnaryInterceptors  []grpc.UnaryServerInterceptor
	clientStreamInterceptors []grpc.StreamClientInterceptor
	clientUnaryInterceptors  []grpc.UnaryClientInterceptor
}

var (
	defaultOptions = &options{
		zapLogger:    nil,
		logrusLogger: nil,
	}
)

func evaluateOpt(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option 选项配置
type Option func(*options)

// WithZapLogger 自定义日志记录器
func WithZapLogger(l *zap.Logger) Option {
	return func(o *options) {
		o.zapLogger = l
	}
}

// WithLogrusLogger 自定义日志记录器
func WithLogrusLogger(l *logrus.Logger) Option {
	return func(o *options) {
		o.logrusLogger = l
	}
}

// WithStreamServerInterceptor 服务端流式拦截器
func WithStreamServerInterceptor(i ...grpc.StreamServerInterceptor) Option {
	return func(o *options) {
		o.serverStreamInterceptors = append(o.serverStreamInterceptors, i...)
	}
}

// WithUnaryServerInterceptor 服务端一元拦截器
func WithUnaryServerInterceptor(i ...grpc.UnaryServerInterceptor) Option {
	return func(o *options) {
		o.serverUnaryInterceptors = append(o.serverUnaryInterceptors, i...)
	}
}

// WithStreamClientInterceptor 客户端流式拦截器
func WithStreamClientInterceptor(i ...grpc.StreamClientInterceptor) Option {
	return func(o *options) {
		o.clientStreamInterceptors = append(o.clientStreamInterceptors, i...)
	}
}

// WithUnaryClientInterceptor 客户端一元拦截器
func WithUnaryClientInterceptor(i ...grpc.UnaryClientInterceptor) Option {
	return func(o *options) {
		o.clientUnaryInterceptors = append(o.clientUnaryInterceptors, i...)
	}
}

// WithTracer 添加Tracer
func WithTracer(t trace.Tracer) Option {
	return func(o *options) {
		o.tracer = t
	}
}

// WithValidator 使用grpc validator
func WithValidator() Option {
	return func(o *options) {
		o.useValidator = true
	}
}

// WithServerPrometheus .
func WithServerPrometheus(ho ...grpc_prometheus.HistogramOption) Option {
	return func(o *options) {
		o.promServerHistogramOptions = ho
	}
}

// WithClientPrometheus .
func WithClientPrometheus(ho ...grpc_prometheus.HistogramOption) Option {
	return func(o *options) {
		o.promServerHistogramOptions = ho
	}
}
