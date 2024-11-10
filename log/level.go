package log

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type Level int64

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func toLogrusLevel(level Level) logrus.Level {
	switch level {
	case LevelDebug:
		return logrus.DebugLevel
	case LevelInfo:
		return logrus.InfoLevel
	case LevelWarn:
		return logrus.WarnLevel
	case LevelError:
		return logrus.ErrorLevel
	case LevelFatal:
		return logrus.FatalLevel
	case LevelPanic:
		return logrus.PanicLevel
	default:
		return logrus.TraceLevel
	}
}

func ParseLevel(level string) (Level, error) {
	switch strings.ToUpper(level) {
	case "TRACE":
		return LevelTrace, nil
	case "DEBUG":
		return LevelDebug, nil
	case "INFO":
		return LevelInfo, nil
	case "WARN":
		return LevelWarn, nil
	case "ERROR":
		return LevelError, nil
	case "FATAL":
		return LevelFatal, nil
	case "PANIC":
		return LevelPanic, nil
	}

	return LevelTrace, fmt.Errorf("not a valid logrus Level: %q", level)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (a *Level) UnmarshalText(text []byte) error {
	level, err := ParseLevel(string(text))
	if err != nil {
		return err
	}

	*a = level
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (a Level) MarshalText() ([]byte, error) {
	switch a {
	case LevelTrace:
		return []byte("TRACE"), nil
	case LevelDebug:
		return []byte("DEBUG"), nil
	case LevelInfo:
		return []byte("INFO"), nil
	case LevelWarn:
		return []byte("WARN"), nil
	case LevelError:
		return []byte("ERROR"), nil
	case LevelFatal:
		return []byte("FATAL"), nil
	case LevelPanic:
		return []byte("PANIC"), nil
	}

	return nil, fmt.Errorf("invalid level %d", a)
}
