package pagetoken

type options struct {
	maxLimit       int    // 最大的limit值
	defaultLimit   int    // 默认的limit值
	defaultOrderBy string // 默认的排序规则
}

var (
	defaultOptions = &options{
		defaultLimit: 500,
		maxLimit:     1000,
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

// WithMaxLimit 限制最大的limit值
func WithMaxLimit(limit int) Option {
	return func(o *options) {
		o.maxLimit = limit
	}
}

// WithDefaultLimit 默认的的limit值
func WithDefaultLimit(limit int) Option {
	return func(o *options) {
		o.defaultLimit = limit
	}
}

// WithDefaultOrderBy 设置默认的order_by
func WithDefaultOrderBy(orderBy string) Option {
	return func(o *options) {
		o.defaultOrderBy = orderBy
	}
}
