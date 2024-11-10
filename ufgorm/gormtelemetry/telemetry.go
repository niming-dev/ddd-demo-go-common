package gormtelemetry

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/niming-dev/ddd-demo/go-common/grpc/middleware/logging/zap/ctxzap"
	"github.com/niming-dev/ddd-demo/go-common/uftelemetry"
	"github.com/niming-dev/ddd-demo/go-common/ufzap"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Logger 实现了gorm日志接口
type Logger struct {
	zapLogger                 *zap.Logger
	telemetryLevel            uftelemetry.Level
	slowThreshold             time.Duration
	skipCallerLookup          bool
	ignoreRecordNotFoundError bool
}

// New 创建Logger实例
func New(opts ...Option) Logger {
	o := evaluateOpts(opts)

	return Logger{
		zapLogger:                 ufzap.ResetLevel(o.zapLogger, LevelToString(o.logLevel)),
		telemetryLevel:            uftelemetry.NewLevel(LevelToString(o.logLevel)),
		slowThreshold:             o.slowThreshold,
		skipCallerLookup:          o.skipCallerLookup,
		ignoreRecordNotFoundError: o.ignoreRecordNotFoundError,
	}
}

// SetAsDefault 设置为gorm的默认Logger
func (l Logger) SetAsDefault() {
	logger.Default = l
}

// LogMode .
func (l Logger) LogMode(level logger.LogLevel) logger.Interface {
	return Logger{
		zapLogger:                 l.zapLogger,
		telemetryLevel:            uftelemetry.NewLevel(LevelToString(level)),
		slowThreshold:             l.slowThreshold,
		skipCallerLookup:          l.skipCallerLookup,
		ignoreRecordNotFoundError: l.ignoreRecordNotFoundError,
	}
}

// Info .
func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	l.logger(ctx).With(ctxzap.TagsToFields(ctx)...).Info(fmt.Sprintf(str, args...))
	span := uftelemetry.SpanFromContext(ctx).StartChild("gorm", "trace")
	defer span.End()
	span.Info(fmt.Sprintf(str, args...))
}

// Warn .
func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	l.logger(ctx).With(ctxzap.TagsToFields(ctx)...).Warn(fmt.Sprintf(str, args...))
	span := uftelemetry.SpanFromContext(ctx).StartChild("gorm", "trace")
	defer span.End()
	span.Warn(fmt.Sprintf(str, args...))
}

// Error .
func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	l.logger(ctx).With(ctxzap.TagsToFields(ctx)...).Error(fmt.Sprintf(str, args...))
	span := uftelemetry.SpanFromContext(ctx).StartChild("gorm", "trace")
	defer span.End()
	span.Error(fmt.Errorf(str, args...))
}

// Trace .
func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	ll := l.logger(ctx).With(ctxzap.TagsToFields(ctx)...)
	span := uftelemetry.SpanFromContext(ctx).StartChild("ufgorm", "gorm.trace")
	defer span.End()
	span.SetLevel(l.telemetryLevel)
	span = uftelemetry.ResetLevelFromCtx(ctx, span)

	attrs := []attribute.KeyValue{
		attribute.String("elapsed", elapsed.String()),
		attribute.Int64("rows", rows),
		attribute.String("sql", sql),
	}
	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}

	// 报错时记录err级别日志
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFoundError {
			fields = append(fields, zap.Error(err))
		}
		ll.Error("gorm.trace", fields...)
		span.Error(err, attrs...)
		return
	}

	// 有慢查询时记录warn级别日志
	if l.slowThreshold != 0 && elapsed > l.slowThreshold {
		ll.Warn("gorm.trace", fields...)
		span.Warn("gorm.trace", attrs...)
		return
	}

	// 默认记录info级别日志
	ll.Info("gorm.trace", fields...)
	span.Tracing(uftelemetry.InfoLevel, "sql", attrs...)
}

var (
	gormPackage    = filepath.Join("gorm.io", "gorm")
	zapGormPackage = filepath.Join("go-common", "ufgorm")
)

func (l Logger) logger(ctx context.Context) *zap.Logger {
	if l.skipCallerLookup {
		return ufzap.ResetLevelFromCtx(ctx, l.zapLogger)
	}

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapGormPackage):
		default:
			return ufzap.ResetLevelFromCtx(ctx, l.zapLogger.WithOptions(zap.AddCallerSkip(i)))
		}
	}
	return ufzap.ResetLevelFromCtx(ctx, l.zapLogger)
}
