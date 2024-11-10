package nsq

import (
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger .
type ZapLogger struct {
	*zap.Logger
}

// NewZapLogger .
func NewZapLogger(logger *zap.Logger) ZapLogger {
	return ZapLogger{logger}
}

// Output 输出日志, s的格式：%-4s %3d [%s/%s] %s
func (l ZapLogger) Output(calldepth int, s string) error {
	parts := strings.Split(s, " ")
	if len(parts) == 0 {
		return nil
	}

	level := nsqLevelToZap(parts[0])
	// 裁切掉LEVEL信息
	parts = parts[1:]

	var fields []zap.Field
	// 如果有ID提取ID信息
	if len(parts) >= 4 {
		nsqID, err := strconv.Atoi(parts[3])
		if err == nil {
			fields = append(fields, zap.Int("nsq.id", nsqID))
			parts = parts[4:]
		}
	}

	if len(parts) > 0 {
		tMsg := parts[0]
		if tMsg[0] == '[' {
			tMsg = strings.TrimFunc(tMsg, func(r rune) bool {
				return r == '[' || r == ']'
			})

			tParts := strings.Split(tMsg, "/")
			if len(tParts) == 2 {
				fields = append(fields, zap.String("nsq.topic", tParts[0]))
				fields = append(fields, zap.String("nsq.channel", tParts[1]))
				parts = parts[1:]
			}
		}
	}

	msg := strings.Join(parts, " ")
	logger := l.Logger
	if calldepth > 0 {
		logger = l.Logger.WithOptions(zap.AddCallerSkip(calldepth))
	}
	logger.Check(level, msg).Write(fields...)
	return nil
}

var nsqLevels = map[string]zapcore.Level{
	"INF": zapcore.InfoLevel,
	"WRN": zapcore.WarnLevel,
	"ERR": zapcore.ErrorLevel,
	"DBG": zapcore.DebugLevel,
}

func nsqLevelToZap(level string) zapcore.Level {
	if l, ok := nsqLevels[level]; ok {
		return l
	}
	return zapcore.InfoLevel
}
