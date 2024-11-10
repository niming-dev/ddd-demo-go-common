package gorm

import (
	"context"
	"errors"
	"time"

	"github.com/niming-dev/ddd-demo/go-common/log"
	"gorm.io/gorm/logger"
)

type gormLogger struct {
	slowThreshold time.Duration
}

func (a *gormLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	// NOT USE THIS, noop
	return a
}
func (a *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Infof(ctx, msg, data)
}
func (a *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warnf(ctx, msg, data)
}
func (a *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Errorf(ctx, msg, data)
}
func (a *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()
	l := log.WithFields(ctx, log.Fields{
		log.FieldKeyModule:  "GORM",
		log.FieldKeyElapsed: elapsed.String(),
		"rowsAffected":      rowsAffected,
	})
	if err != nil && !errors.Is(err, logger.ErrRecordNotFound) {
		l.Errorf("%s %v", sql, err)
	} else if err == nil && elapsed > a.slowThreshold && a.slowThreshold != 0 {
		l.Warnf("%s", sql)
	} else {
		l.Infof("%s", sql)
	}
}
