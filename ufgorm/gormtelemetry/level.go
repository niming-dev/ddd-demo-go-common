package gormtelemetry

import (
	"gorm.io/gorm/logger"
)

var (
	LevelName = map[logger.LogLevel]string{
		logger.Silent: "silent",
		logger.Error:  "error",
		logger.Warn:   "warn",
		logger.Info:   "info",
	}
	LevelValue = map[string]logger.LogLevel{
		"silent": logger.Silent,
		"error":  logger.Error,
		"warn":   logger.Warn,
		"info":   logger.Info,
	}
)

// LevelToString 日志级别转字符串
func LevelToString(level logger.LogLevel) string {
	name, ok := LevelName[level]
	if ok {
		return name
	}
	return "info"
}

// NewLevel 字符串转日志级别
func NewLevel(level string) logger.LogLevel {
	value, ok := LevelValue[level]
	if ok {
		return value
	}
	return logger.Info
}
