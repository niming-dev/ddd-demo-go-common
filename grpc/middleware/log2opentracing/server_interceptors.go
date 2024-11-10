package log2opentracing

import (
	"bytes"
	"context"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/modern-go/reflect2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
)

var (
	// JsonPbMarshaller 用于序列化pb消息
	JsonPbMarshaller grpc_logging.JsonPbMarshaler = &jsonpb.Marshaler{}
)

// UnaryServerInterceptor 一元请求拦截器
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !o.shouldLog(ctx, info.FullMethod, info.Server) {
			return handler(ctx, req)
		}

		logProtoMessageAsJson(ctx, req, "grpc.request.content")
		resp, err := handler(ctx, req)
		logProtoMessageAsJson(ctx, resp, "grpc.response.content")
		return resp, err
	}
}

// StreamServerInterceptor 流式请求拦截器
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !o.shouldLog(stream.Context(), info.FullMethod, srv) {
			return handler(srv, stream)
		}

		newStream := &loggingServerStream{ServerStream: stream}
		return handler(srv, newStream)
	}
}

func logProtoMessageAsJson(ctx context.Context, pbMsg interface{}, key string) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}
	if reflect2.IsNil(pbMsg) {
		span.LogFields(log.String("event", key), log.String("message", "{}"))
		return
	}

	if p, ok := pbMsg.(proto.Message); ok {
		logMsg, err := marshalJSON(p)
		if err != nil {
			span.LogFields(log.String("event", "error"), log.Error(err))
			return
		}
		span.LogFields(log.String("event", key), log.String("message", string(logMsg)))
	}
}

func marshalJSON(pb proto.Message) ([]byte, error) {
	b := &bytes.Buffer{}
	if err := JsonPbMarshaller.Marshal(b, pb); err != nil {
		return nil, fmt.Errorf("jsonpb serializer failed: %v", err)
	}

	return b.Bytes(), nil
}

type loggingServerStream struct {
	grpc.ServerStream
}

func (l *loggingServerStream) SendMsg(m interface{}) error {
	err := l.ServerStream.SendMsg(m)
	if err == nil {
		logProtoMessageAsJson(l.ServerStream.Context(), m, "grpc.response.content")
	}
	return err
}

func (l *loggingServerStream) RecvMsg(m interface{}) error {
	err := l.ServerStream.RecvMsg(m)
	if err == nil {
		logProtoMessageAsJson(l.ServerStream.Context(), m, "grpc.request.content")
	}
	return err
}
