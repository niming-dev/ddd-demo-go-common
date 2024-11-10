// Package mdhelper metadata plus library
package mdhelper

import (
	"context"

	grpcopentelemetry "github.com/niming-dev/ddd-demo/go-common/grpc/middleware/opentelemetry"
	"google.golang.org/grpc/metadata"
)

// ExtractToOutCtx 从ctx中提取需要的信息并返回一个 outgoing context
func ExtractToOutCtx(ctx context.Context, opts ...ExtractorOption) context.Context {
	o := extractorEvaluateOpt(opts)

	// 尝试从ctx中获取 open telemetry 的追踪ID
	newCtx := grpcopentelemetry.InjectTraceIdsToOutCtx(ctx)
	// 检查如果 newCtx 中没有 outgoing metadata则写入一个空的
	rawMD, _, ok := metadata.FromOutgoingContextRaw(newCtx)
	if !ok {
		rawMD = metadata.MD{}
		newCtx = metadata.NewOutgoingContext(newCtx, rawMD)
	}

	// 把ctx当作 incoming context 提取key对应的values并写入到 outgoing context 中
	if inMD, ok := metadata.FromIncomingContext(ctx); ok {
		for _, key := range o.keys {
			values := inMD.Get(key)
			if len(values) > 0 {
				rawMD.Set(key, values...)
			}
		}
	}

	return newCtx
}
