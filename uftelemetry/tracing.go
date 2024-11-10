package uftelemetry

import (
	"strings"
)

var (
	// DefaultTracingLevel 默认日志追踪级别，此包内的函数只记录级别高于默认级别的日志
	DefaultTracingLevel = InfoLevel
)

type Level int8

const (
	Unspecified Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

var (
	LevelName = map[Level]string{
		DebugLevel: "debug",
		InfoLevel:  "info",
		WarnLevel:  "warn",
		ErrorLevel: "error",
	}
	LevelValue = map[string]Level{
		"debug": DebugLevel,
		"info":  InfoLevel,
		"warn":  WarnLevel,
		"error": ErrorLevel,
	}
)

func (l Level) String() string {
	return LevelName[l]
}

// NewLevel .
func NewLevel(level string) Level {
	if v, ok := LevelValue[strings.ToLower(level)]; ok {
		return v
	}

	return DefaultTracingLevel
}
