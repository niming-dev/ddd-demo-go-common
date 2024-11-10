package log

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

type myContext struct {
	context.Context
	Logger Logger
}

func (a myContext) GetLogger() Logger {
	return a.Logger
}

func setupContext() context.Context {
	ctxLogger := logrus.WithFields(logrus.Fields{
		"ctxField1": "ctxValue1",
	})

	ctx := context.Background()
	return &myContext{ctx, NewFromLogrusEntry(ctxLogger)}
}

func TestGlobal(t *testing.T) {
	ctx := context.Background()

	fields := Fields{
		"field1": "value1",
	}
	msg := "log.Debug"
	WithFields(ctx, fields).Info(msg)
}

func TestContext(t *testing.T) {
	ctx := setupContext()

	WithField(ctx, "field1", "value1").Info("logger.Info")
	Info(ctx, "logger.Info")
	WithField(ctx, "field1", "value1").Infof("logger.Info %s %v", "ctx", ctx)
	Infof(ctx, "logger.Info %s %v", "ctx", ctx)
}
