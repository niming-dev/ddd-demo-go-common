package log

import (
	"context"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

type contextKey int64

const (
	contextKeyLogFields contextKey = iota
)

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	SetLevel(level Level)
	SetOutput(w io.Writer)
}

var defaultLogger Logger

func init() {
	logger := logrus.New()
	logger.SetFormatter(defaultFormatter())
	logger.SetLevel(logrus.InfoLevel)
	defaultLogger = NewFromLogrus(logger)
}

func NewLogger() Logger {
	logger := logrus.New()
	logger.SetFormatter(defaultFormatter())
	logger.SetLevel(logrus.InfoLevel)
	return NewFromLogrus(logger)
}

func WithField(ctx context.Context, key string, value interface{}) Logger {
	fields := ExtractLogFields(ctx)
	fields[key] = value
	return defaultLogger.WithFields(fields)
}

func WithFields(ctx context.Context, fields Fields) Logger {
	ctxFields := ExtractLogFields(ctx)
	for k, v := range fields {
		ctxFields[k] = v
	}
	return defaultLogger.WithFields(ctxFields)
}

func Debug(ctx context.Context, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Debug(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Debugf(format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Infof(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Warn(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Warnf(format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Errorf(format, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Fatal(args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Fatalf(format, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Panic(args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.WithFields(ExtractLogFields(ctx)).Panicf(format, args...)
}

func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

func WithFieldContext(ctx context.Context, key string, value interface{}) context.Context {
	ctxFields := ExtractLogFields(ctx)
	ctxFields[key] = value
	return context.WithValue(ctx, contextKeyLogFields, ctxFields)
}

func WithFieldsContext(ctx context.Context, fields Fields) context.Context {
	ctxFields := ExtractLogFields(ctx)
	for k, v := range fields {
		ctxFields[k] = v
	}
	return context.WithValue(ctx, contextKeyLogFields, ctxFields)
}

func GetTimestampFormat() string {
	return defaultTimestampFormat
}

func defaultFormatter() logrus.Formatter {
	return &TextFormatter{
		TimestampFormat: time.RFC3339,
	}
}
