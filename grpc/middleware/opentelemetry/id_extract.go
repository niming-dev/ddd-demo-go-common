package grpcopentelemetry

import (
	"context"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

const (
	TagTraceId = "trace.traceid"
	TagSpanId  = "trace.spanid"
	TagSampled = "trace.sampled"
)

var (
	traceContext      = propagation.TraceContext{}
	traceParentHeader = traceContext.Fields()[0]
	traceStateHeader  = traceContext.Fields()[1]
)

// InjectTraceIdsToTags writes trace data to ctxtags.
func InjectTraceIdsToTags(span trace.Span, tags grpc_ctxtags.Tags) {
	spanCtx := span.SpanContext()
	if !spanCtx.IsValid() {
		return
	}
	tags.Set(TagTraceId, spanCtx.TraceID().String())
	tags.Set(TagSpanId, spanCtx.SpanID().String())
	tags.Set(TagSampled, spanCtx.IsSampled())
}

// InjectTraceIdsToOutCtx 把ctx中的trace_id流入到一个新的outgoingContext，方便grpc调用时传递trace数据
func InjectTraceIdsToOutCtx(ctx context.Context) context.Context {
	mc := MetadataCarrier{}
	traceContext.Inject(ctx, mc)
	if len(mc) == 0 {
		return ctx
	}

	return metadata.NewOutgoingContext(ctx, metadata.MD(mc))
}
