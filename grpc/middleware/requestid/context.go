package requestid

import (
	"context"

	"google.golang.org/grpc/metadata"
)

var (
	// ContextKey RequestId在Client/Server中的key名称
	ContextKey = "grpc.request.id"
)

// extractIdFromInCtx 从IncomingContext中获取uuid
func extractIdFromInCtx(ctx context.Context) string {
	inMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if uuids := inMD.Get(ContextKey); len(uuids) > 0 {
		return uuids[0]
	}

	return ""
}

// extractIdFromOutCtx 从OutgoingContext中获取uuid
func extractIdFromOutCtx(ctx context.Context) string {
	inMD, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	if uuids := inMD.Get(ContextKey); len(uuids) > 0 {
		return uuids[0]
	}

	return ""
}
