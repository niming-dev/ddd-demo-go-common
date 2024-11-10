package requestid

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryClientInterceptor 一元请求客户端拦截器，为ctx中放入uuid
func UnaryClientInterceptor(opts ...Option) grpc.UnaryClientInterceptor {
	o := evaluateOpt(opts)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx := newTagsForClientCtx(ctx, o)
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

// StreamClientInterceptor 流式请求客户端拦截器，为ctx中放入uuid
func StreamClientInterceptor(opts ...Option) grpc.StreamClientInterceptor {
	o := evaluateOpt(opts)
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		newCtx := newTagsForClientCtx(ctx, o)
		return streamer(newCtx, desc, cc, method, opts...)
	}
}

// UnaryServerInterceptor 一元请求服务端拦截器，ctx中不存在uuid的时候自动生成一个
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOpt(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx := newTagsForSrvCtx(ctx, o)
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor 流式请求服务端拦截器，ctx中不存在uuid的时候自动生成一个
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateOpt(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Short-circuit, don't do the expensive bit of allocating a wrappedStream.
		wrappedStream := grpc_middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = newTagsForSrvCtx(stream.Context(), o)
		return handler(srv, wrappedStream)
	}
}

// newTagsForSrvCtx 为server context生成一个带有uuid的tags
func newTagsForSrvCtx(ctx context.Context, o *options) context.Context {
	if t := grpc_ctxtags.Extract(ctx); t == grpc_ctxtags.NoopTags {
		newTags := grpc_ctxtags.NewTags()
		extractIdFromInCtxToTags(ctx, newTags, o)
		return grpc_ctxtags.SetInContext(ctx, newTags)
	} else {
		extractIdFromInCtxToTags(ctx, t, o)
	}

	return ctx
}

// extractIdFromInCtxToTags 从context中提取中request uuid并设置到tags中
func extractIdFromInCtxToTags(ctx context.Context, t grpc_ctxtags.Tags, o *options) {
	if uuid := extractIdFromInCtx(ctx); uuid != "" {
		t.Set(ContextKey, uuid)
		return
	}

	t.Set(ContextKey, o.idBuilder())
}

// newTagsForClientCtx 为client context生成一个带有uuid的tags
func newTagsForClientCtx(ctx context.Context, o *options) context.Context {
	// 尝试从原始请求的context中获取uuid
	uuid := extractIdFromInCtx(ctx)
	if uuid == "" {
		// 尝试从原始context中解析tags并取出tags中的uuid
		if t := grpc_ctxtags.Extract(ctx); t != grpc_ctxtags.NoopTags {
			uuid, _ = (t.Values()[ContextKey]).(string)
		}
	}
	if uuid == "" {
		uuid = o.idBuilder()
	}

	outCtx := ctx
	// 为outCtx设置uuid，当outCtx已有uuid时取出保存到uuid变量中
	if outId := extractIdFromOutCtx(outCtx); outId == "" {
		outCtx = metadata.AppendToOutgoingContext(outCtx, ContextKey, uuid)
	} else {
		uuid = outId
	}

	// 为outCtx设置tags
	t := grpc_ctxtags.Extract(outCtx)
	if t == grpc_ctxtags.NoopTags {
		newTags := grpc_ctxtags.NewTags().Set(ContextKey, uuid)
		outCtx = grpc_ctxtags.SetInContext(outCtx, newTags)
	} else {
		t.Set(ContextKey, uuid)
		outCtx = grpc_ctxtags.SetInContext(outCtx, t)
	}

	return outCtx
}
