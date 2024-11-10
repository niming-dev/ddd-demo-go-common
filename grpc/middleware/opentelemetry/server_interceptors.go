package grpcopentelemetry

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/niming-dev/ddd-demo/go-common/strsconv"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	grpcCodeKey = "grpc.code"
	grpcMsgKey  = "grpc.message"
)

var (
	grpcAttr = attribute.String("component", "gRPC")
	traceCtx = propagation.TraceContext{}
)

// UnaryServerInterceptor returns a new unary server interceptor for OpenTracing.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if o.filterOutFunc != nil && !o.filterOutFunc(ctx, info.FullMethod) {
			return handler(ctx, req)
		}
		opName := info.FullMethod
		if o.opNameFunc != nil {
			opName = o.opNameFunc(info.FullMethod)
		}
		newCtx, serverSpan := newServerSpanFromInbound(ctx, o.tracer, opName)
		if o.unaryRequestHandlerFunc != nil {
			o.unaryRequestHandlerFunc(serverSpan, req)
		}
		resp, err := handler(newCtx, req)
		finishServerSpan(newCtx, serverSpan, err)
		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for OpenTracing.
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateOptions(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if o.filterOutFunc != nil && !o.filterOutFunc(stream.Context(), info.FullMethod) {
			return handler(srv, stream)
		}
		opName := info.FullMethod
		if o.opNameFunc != nil {
			opName = o.opNameFunc(info.FullMethod)
		}
		newCtx, serverSpan := newServerSpanFromInbound(stream.Context(), o.tracer, opName)
		wrappedStream := grpc_middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = newCtx
		err := handler(srv, wrappedStream)
		finishServerSpan(newCtx, serverSpan, err)
		return err
	}
}

func newServerSpanFromInbound(ctx context.Context, tracer trace.Tracer, opName string) (context.Context, trace.Span) {
	md := metautils.ExtractIncoming(ctx)
	parentSpanContext := traceCtx.Extract(ctx, MetadataCarrier(md))
	newCtx, serverSpan := tracer.Start(parentSpanContext, opName, trace.WithAttributes(grpcAttr))
	InjectTraceIdsToTags(serverSpan, grpc_ctxtags.Extract(ctx))
	return newCtx, serverSpan
}

func finishServerSpan(ctx context.Context, serverSpan trace.Span, err error) {
	// Log context information
	tags := grpc_ctxtags.Extract(ctx)
	for k, v := range tags.Values() {
		// Don't tag errors, log them instead.
		if vErr, ok := v.(error); ok {
			serverSpan.RecordError(vErr)
		} else {
			serverSpan.SetAttributes(attribute.String(k, strsconv.Any2String(v)))
		}
	}
	if err != nil {
		if state, ok := status.FromError(err); ok {
			serverSpan.SetAttributes(attribute.String(grpcCodeKey, state.Code().String()))
			serverSpan.SetAttributes(attribute.String(grpcMsgKey, state.Message()))
		}
		serverSpan.RecordError(err)
	} else {
		serverSpan.SetAttributes(attribute.String(grpcCodeKey, "OK"))
	}
	serverSpan.End()
}
