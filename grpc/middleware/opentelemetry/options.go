package grpcopentelemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

var (
	defaultOptions = &options{
		filterOutFunc: nil,
		tracer:        nil,
	}
)

// FilterFunc allows users to provide a function that filters out certain methods from being traced.
//
// If it returns false, the given request will not be traced.
type FilterFunc func(ctx context.Context, fullMethodName string) bool

// UnaryRequestHandlerFunc is a custom request handler
type UnaryRequestHandlerFunc func(span trace.Span, req interface{})

// OpNameFunc is a func that allows custom operation names instead of the gRPC method.
type OpNameFunc func(method string) string

type options struct {
	filterOutFunc           FilterFunc
	tracer                  trace.Tracer
	unaryRequestHandlerFunc UnaryRequestHandlerFunc
	opNameFunc              OpNameFunc
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	if optCopy.tracer == nil {
		optCopy.tracer = trace.NewNoopTracerProvider().Tracer("noop")
	}
	return optCopy
}

type Option func(*options)

// WithFilterFunc customizes the function used for deciding whether a given call is traced or not.
func WithFilterFunc(f FilterFunc) Option {
	return func(o *options) {
		o.filterOutFunc = f
	}
}

// WithTracer sets a custom tracer to be used for this middleware, otherwise the opentracing.GlobalTracer is used.
func WithTracer(tracer trace.Tracer) Option {
	return func(o *options) {
		o.tracer = tracer
	}
}

// WithUnaryRequestHandlerFunc sets a custom handler for the request
func WithUnaryRequestHandlerFunc(f UnaryRequestHandlerFunc) Option {
	return func(o *options) {
		o.unaryRequestHandlerFunc = f
	}
}

// WithOpName customizes the trace Operation name
func WithOpName(f OpNameFunc) Option {
	return func(o *options) {
		o.opNameFunc = f
	}
}
