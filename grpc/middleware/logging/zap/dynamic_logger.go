package grpc_zap

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
)

const (
	HeaderLogLevelKey = "x-uf-log-level"
)

// resetLevel 根据context中的LogLevel生成一个新的Logger
func resetLevel(ctx context.Context, logger *zap.Logger) *zap.Logger {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return logger
	}

	vss := md.Get(HeaderLogLevelKey)
	if len(vss) == 0 || vss[0] == "" {
		return logger
	}

	newLevel := zapcore.InfoLevel
	if err := newLevel.Set(vss[0]); err != nil {
		// 新的level字符串不在zap.Core中
		return logger
	}

	return logger.WithOptions(zap.IncreaseLevel(newLevel))
}
