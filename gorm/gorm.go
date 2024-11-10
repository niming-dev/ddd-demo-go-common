package gorm

import (
	"time"

	"gorm.io/gorm"
)

var (
	defaultConnMaxIdleTime = 10 * time.Minute
	defaultConnMaxLifetime = 2 * time.Hour
	defaultMaxOpenConns    = 100
	defaultMaxIdleConns    = 10
	defaultLogger          *gormLogger
	defaultOption          *loggerOption
)

func init() {
	defaultLogger = &gormLogger{
		slowThreshold: 200 * time.Millisecond,
	}
	defaultOption = &loggerOption{
		Override: false,
		Logger:   defaultLogger,
	}
}

func Open(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
	opts = append(opts, defaultOption)
	db, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxIdleTime(defaultConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(defaultConnMaxLifetime)
	sqlDB.SetMaxOpenConns(defaultMaxOpenConns)
	sqlDB.SetMaxIdleConns(defaultMaxIdleConns)

	return db, nil
}
