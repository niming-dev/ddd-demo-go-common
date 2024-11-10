package uftelemetry

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// SpanFromContext 从ctx中提取span，如果ctx中无span则创建一个noopSpan，并封装到Span实例
func SpanFromContext(ctx context.Context) *Span {
	return ResetLevelFromCtx(ctx, NewSpan(ctx, trace.SpanFromContext(ctx)))
}

// StartChildSpanFromTracer 从trace.Tracer + context中创建一个Span
func StartChildSpanFromTracer(ctx context.Context, tracer trace.Tracer, operationName string, opts ...trace.SpanStartOption) *Span {
	newCtx, span := tracer.Start(ctx, operationName, opts...)
	return ResetLevelFromCtx(ctx, NewSpan(newCtx, span))
}

// StartChildSpan 从ctx中提取span，如果ctx中无span则创建一个noopSpan，并创建一个子span
func StartChildSpan(ctx context.Context, libraryName, spanName string) *Span {
	s := SpanFromContext(ctx)
	// 如果operationName为空，则取其父spanName
	if spanName == "" {
		spanName = extractSpanName(s.Span)
	}
	return s.StartChild(libraryName, spanName)
}

// NewSpan .
func NewSpan(ctx context.Context, span trace.Span) *Span {
	return &Span{Context: ctx, Span: span, level: DefaultTracingLevel}
}

// Span 实现opentracing.Span，封装一些方便操作的功能
type Span struct {
	// 当前span的context
	context.Context
	// 原始span
	trace.Span
	// 追踪日志级别
	level Level
	// Error信息中是否包含调用栈信息
	tracingStack bool
}

// SetLevel 设置记录级别
func (s *Span) SetLevel(l Level) *Span {
	s.level = l
	return s
}

// SetTracingStack Error信息中是否包含调用栈信息
func (s *Span) SetTracingStack(on bool) *Span {
	s.tracingStack = on
	return s
}

// Tracing 记录日志追踪信息
func (s *Span) Tracing(l Level, msg string, attrs ...attribute.KeyValue) *Span {
	if s.level != Unspecified {
		if l < s.level {
			return s
		}
	} else if l < DefaultTracingLevel {
		return s
	}

	// 设置error属性为true
	if l == ErrorLevel {
		attrs = append(attrs, attribute.Bool("error", true))
	}

	// 把level转为event字段
	if v := l.String(); v != "" {
		attrs = append(attrs, attribute.String("trace.level", strings.ToLower(v)))
	}

	s.AddEvent(msg, trace.WithAttributes(attrs...))
	return s
}

// Debug .
func (s *Span) Debug(msg string, attrs ...attribute.KeyValue) *Span {
	s.Tracing(DebugLevel, msg, attrs...)
	return s
}

// Info .
func (s *Span) Info(msg string, attrs ...attribute.KeyValue) *Span {
	s.Tracing(InfoLevel, msg, attrs...)
	return s
}

// Error .
func (s *Span) Error(err error, attrs ...attribute.KeyValue) *Span {
	if err == nil {
		return s
	}

	if s.tracingStack {
		s.RecordError(err, trace.WithStackTrace(true))
	} else {
		s.RecordError(err)
	}
	s.Tracing(ErrorLevel, "error", attrs...)
	return s
}

// CheckError 记录err信息，如果err=nil则什么都不做
func (s *Span) CheckError(err error, attrs ...attribute.KeyValue) *Span {
	if err == nil {
		return s
	}
	return s.Error(err, attrs...)
}

// Warn .
func (s *Span) Warn(msg string, attrs ...attribute.KeyValue) *Span {
	s.Tracing(WarnLevel, msg, attrs...)
	return s
}

// StartChild 从当前Span创建一个子Span
func (s *Span) StartChild(libraryName, spanName string) *Span {
	if spanName == "" {
		spanName = extractSpanName(s.Span)
	}
	newCtx, span := s.TracerProvider().Tracer(libraryName).Start(s.Context, spanName)
	return &Span{Context: newCtx, Span: span}
}

func extractSpanName(span trace.Span) string {
	if rSpan, ok := span.(tracesdk.ReadOnlySpan); ok {
		rSpan.Name()
	}
	return ""
}
