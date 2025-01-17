package grpc_zap

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
)

// AddFields adds zap fields to the logger.
// Deprecated: should use the ctxzap.AddFields instead
func AddFields(ctx context.Context, fields ...zapcore.Field) {
	ctxzap.AddFields(ctx, fields...)
}

// Extract takes the call-scoped Logger from grpc_zap middleware.
// Deprecated: should use the ctxzap.Extract instead
func Extract(ctx context.Context) *zap.Logger {
	return ctxzap.Extract(ctx)
}

// ExtractInMetadataToField 从context中提取in metadata并转为field数组
func ExtractInMetadataToField(ctx context.Context) []zapcore.Field {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return []zapcore.Field{zap.Any("grpc.metadata", md)}
}

// ExtractOutMetadataToField 从context中提取out metadata并转为field数组
func ExtractOutMetadataToField(ctx context.Context) []zapcore.Field {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return nil
	}
	return []zapcore.Field{zap.Any("grpc.metadata", md)}
}
