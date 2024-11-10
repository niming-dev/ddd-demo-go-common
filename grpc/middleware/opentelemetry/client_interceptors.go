package grpcopentelemetry

import (
	"context"
	"io"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryClientInterceptor returns a new unary client interceptor for OpenTracing.
func UnaryClientInterceptor(opts ...Option) grpc.UnaryClientInterceptor {
	o := evaluateOptions(opts)
	return func(parentCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if o.filterOutFunc != nil && !o.filterOutFunc(parentCtx, method) {
			return invoker(parentCtx, method, req, reply, cc, opts...)
		}
		newCtx, clientSpan := newClientSpanFromContext(parentCtx, o.tracer, method)
		if o.unaryRequestHandlerFunc != nil {
			o.unaryRequestHandlerFunc(clientSpan, req)
		}
		err := invoker(newCtx, method, req, reply, cc, opts...)
		finishClientSpan(clientSpan, err)
		return err
	}
}

// StreamClientInterceptor returns a new streaming client interceptor for OpenTracing.
func StreamClientInterceptor(opts ...Option) grpc.StreamClientInterceptor {
	o := evaluateOptions(opts)
	return func(parentCtx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if o.filterOutFunc != nil && !o.filterOutFunc(parentCtx, method) {
			return streamer(parentCtx, desc, cc, method, opts...)
		}
		newCtx, clientSpan := newClientSpanFromContext(parentCtx, o.tracer, method)
		clientStream, err := streamer(newCtx, desc, cc, method, opts...)
		if err != nil {
			finishClientSpan(clientSpan, err)
			return nil, err
		}
		return &tracedClientStream{ClientStream: clientStream, clientSpan: clientSpan}, nil
	}
}

// type serverStreamingRetryingStream is the implementation of grpc.ClientStream that acts as a
// proxy to the underlying call. If any of the RecvMsg() calls fail, it will try to reestablish
// a new ClientStream according to the retry policy.
type tracedClientStream struct {
	grpc.ClientStream
	mu              sync.Mutex
	alreadyFinished bool
	clientSpan      trace.Span
}

func (s *tracedClientStream) Header() (metadata.MD, error) {
	h, err := s.ClientStream.Header()
	if err != nil {
		s.finishClientSpan(err)
	}
	return h, err
}

func (s *tracedClientStream) SendMsg(m interface{}) error {
	err := s.ClientStream.SendMsg(m)
	if err != nil {
		s.finishClientSpan(err)
	}
	return err
}

func (s *tracedClientStream) CloseSend() error {
	err := s.ClientStream.CloseSend()
	s.finishClientSpan(err)
	return err
}

func (s *tracedClientStream) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	if err != nil {
		s.finishClientSpan(err)
	}
	return err
}

func (s *tracedClientStream) finishClientSpan(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.alreadyFinished {
		finishClientSpan(s.clientSpan, err)
		s.alreadyFinished = true
	}
}

// ClientAddContextAttributes returns a context with specified opentracing tags, which
// are used by UnaryClientInterceptor/StreamClientInterceptor when creating a
// new span.
func ClientAddContextAttributes(ctx context.Context, tags attributes.Attributes) context.Context {
	return context.WithValue(ctx, clientSpanTagKey{}, tags)
}

type clientSpanTagKey struct{}

func newClientSpanFromContext(ctx context.Context, tracer trace.Tracer, fullMethodName string) (context.Context, trace.Span) {
	var parentSpanCtx trace.SpanContext
	parent := trace.SpanFromContext(ctx)
	parentSpanCtx = parent.SpanContext()

	opts := []trace.SpanStartOption{
		trace.WithAttributes(grpcAttr),
	}
	if attrX := ctx.Value(clientSpanTagKey{}); attrX != nil {
		if opt, ok := attrX.(trace.SpanStartOption); ok {
			opts = append(opts, opt)
		}
	}

	// Make sure we add this to the metadata of the call, so it gets propagated:
	md := metautils.ExtractOutgoing(ctx).Clone()
	traceCtx.Inject(ctx, MetadataCarrier(md))
	ctxWithMetadata := md.ToOutgoing(ctx)
	newCtx, clientSpan := tracer.Start(
		trace.ContextWithRemoteSpanContext(ctxWithMetadata, parentSpanCtx),
		fullMethodName,
		opts...,
	)

	return newCtx, clientSpan
}

func finishClientSpan(clientSpan trace.Span, err error) {
	if err != nil && err != io.EOF {
		if state, ok := status.FromError(err); ok {
			clientSpan.SetAttributes(attribute.String(grpcCodeKey, state.Code().String()))
			clientSpan.SetAttributes(attribute.String(grpcMsgKey, state.Message()))
		}
		clientSpan.RecordError(err)
		clientSpan.SetStatus(codes.Error, err.Error())
	} else {
		clientSpan.SetAttributes(attribute.String(grpcCodeKey, "OK"))
		clientSpan.SetStatus(codes.Ok, "")
	}
	clientSpan.End()
}
