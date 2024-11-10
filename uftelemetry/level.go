package uftelemetry

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	HeaderLevelKey = "x-uf-log-level"
)

// ResetLevelFromCtx 根据context中的LogLevel生成一个新的Logger
func ResetLevelFromCtx(ctx context.Context, span *Span) *Span {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return span
	}

	vss := md.Get(HeaderLevelKey)
	if len(vss) == 0 || vss[0] == "" {
		return span
	}

	if v, ok := LevelValue[vss[0]]; ok {
		span.level = v
	}
	return span
}
