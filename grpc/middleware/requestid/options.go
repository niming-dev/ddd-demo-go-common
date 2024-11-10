package requestid

import (
	"github.com/google/uuid"
)

type options struct {
	idBuilder IdBuilder
}

var (
	defaultOptions = &options{
		idBuilder: defaultIdBuilder,
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

// IdBuilder uuid生成器
type IdBuilder func() string

// WithIdBuilder 自定义uuid生成器
func WithIdBuilder(f IdBuilder) Option {
	return func(o *options) {
		o.idBuilder = f
	}
}

// defaultIdBuilder 默认的uuid生成器
func defaultIdBuilder() string {
	return uuid.New().String()
}
