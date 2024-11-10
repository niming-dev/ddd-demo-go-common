package mdhelper

const (
	// DynamicLogLevelKey 动态日志级别key名称
	DynamicLogLevelKey = "x-uf-log-level"
	// WorkflowTracerOn 开启工作流服务的追踪功能，后续可能会调整为由动态日志级别控制
	WorkflowTracerOn = "x-uf-workflow-tracer"
)

var (
	// DefaultExtractInKeys 默认需要提取的key列表
	DefaultExtractInKeys = []string{DynamicLogLevelKey, WorkflowTracerOn}
)

type extractorOptions struct {
	keys []string
}

var (
	extractorDefaultOptions = &extractorOptions{
		keys: DefaultExtractInKeys,
	}
)

func extractorEvaluateOpt(opts []ExtractorOption) *extractorOptions {
	optCopy := &extractorOptions{}
	*optCopy = *extractorDefaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// ExtractorOption 选项配置
type ExtractorOption func(*extractorOptions)

// WithExtractorKeys 使用指定的keys列表
func WithExtractorKeys(keys []string) ExtractorOption {
	return func(options *extractorOptions) {
		options.keys = keys
	}
}
