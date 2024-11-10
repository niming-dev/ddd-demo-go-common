package gormtelemetry

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type options struct {
	zapLogger                 *zap.Logger
	logLevel                  logger.LogLevel
	slowThreshold             time.Duration
	skipCallerLookup          bool
	ignoreRecordNotFoundError bool
}

var (
	defaultOptions = &options{
		zapLogger:                 zap.NewNop(),
		logLevel:                  logger.Info,
		slowThreshold:             100 * time.Millisecond,
		skipCallerLookup:          false,
		ignoreRecordNotFoundError: false,
	}
)

func evaluateOpts(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option 选项配置
type Option func(*options)

// WithZap 设置日志级别
func WithZap(l *zap.Logger) Option {
	return func(o *options) {
		if l == nil {
			return
		}
		o.zapLogger = l
	}
}

// WithLogLevel 设置日志级别
func WithLogLevel(l logger.LogLevel) Option {
	return func(o *options) {
		o.logLevel = l
	}
}

// WithSlowThreshold 慢查询阈值
func WithSlowThreshold(t time.Duration) Option {
	return func(o *options) {
		o.slowThreshold = t
	}
}

// WithSkipCallerLookup 跳过调用者查找
func WithSkipCallerLookup(b bool) Option {
	return func(o *options) {
		o.skipCallerLookup = b
	}
}

// WithIgnoreRecordNotFoundError 忽略未找到记录错误
func WithIgnoreRecordNotFoundError(b bool) Option {
	return func(o *options) {
		o.ignoreRecordNotFoundError = b
	}
}
