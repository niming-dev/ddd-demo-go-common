package log2opentracing

import (
	"context"
)

type options struct {
	shouldLog Decider
}

var (
	defaultServerOptions = &options{
		shouldLog: DefaultDeciderMethod,
	}
)

func evaluateServerOpt(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultServerOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option 选项配置
type Option func(*options)

// WithDecider 设置一个函数决定是否记录日志信息
func WithDecider(f Decider) Option {
	return func(o *options) {
		o.shouldLog = f
	}
}

// Decider function defines rules for suppressing any interceptor logs
type Decider func(ctx context.Context, fullMethodName string, servingObject interface{}) bool

// DefaultDeciderMethod is the default implementation of decider to see if you should log the call
// by default this if always true so all calls are logged
func DefaultDeciderMethod(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
	return true
}
