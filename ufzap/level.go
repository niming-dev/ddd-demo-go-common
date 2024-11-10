package ufzap

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
)

const (
	// DefaultLevel 指定Zap默认的日志级别
	DefaultLevel = zap.InfoLevel

	// HeaderLevelKey 动态日志级别key名称
	HeaderLevelKey = "x-uf-log-level"
)

var (
	// LevelValue 日志级别字符串与Zap日志级别对应关系
	LevelValue = map[string]zapcore.Level{
		"debug":  zap.DebugLevel,
		"info":   zap.InfoLevel,
		"warn":   zap.WarnLevel,
		"error":  zap.ErrorLevel,
		"dpanic": zap.DPanicLevel,
		"panic":  zap.PanicLevel,
		"fatal":  zap.FatalLevel,
	}
	LevelName = map[zapcore.Level]string{
		zap.DebugLevel:  "debug",
		zap.InfoLevel:   "info",
		zap.WarnLevel:   "warn",
		zap.ErrorLevel:  "error",
		zap.DPanicLevel: "dpanic",
		zap.PanicLevel:  "panic",
		zap.FatalLevel:  "fatal",
	}
)

// NewLevel 转换字符串日志级别到zapcore.Level.
func NewLevel(lvl string) zapcore.Level {
	if l, ok := LevelValue[lvl]; ok {
		return l
	}

	return DefaultLevel
}

// NewAtomicLevel 转换字符串日志级别到zap.AtomicLevel.
func NewAtomicLevel(lvl string) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(NewLevel(lvl))
}

// ResetLevelFromCtx 根据context中的LogLevel生成一个新的Logger
func ResetLevelFromCtx(ctx context.Context, zapLogger *zap.Logger) *zap.Logger {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return zapLogger
	}

	vss := md.Get(HeaderLevelKey)
	if len(vss) == 0 || vss[0] == "" {
		return zapLogger
	}

	newLevel := DefaultLevel
	if err := newLevel.Set(vss[0]); err != nil {
		// 如果新的level字符串不在zap.Core中
		return zapLogger
	}

	return zapLogger.WithOptions(zap.IncreaseLevel(newLevel))
}

func ResetLevel(logger *zap.Logger, level string) *zap.Logger {
	newLevel := DefaultLevel
	if err := newLevel.Set(level); err != nil {
		// 如果新的level字符串不在zap.Core中
		return logger
	}

	return logger.WithOptions(zap.IncreaseLevel(newLevel))
}
