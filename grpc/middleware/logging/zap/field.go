package grpc_zap

import (
	"go.uber.org/zap/zapcore"
)

// mergeFields 合并多个field数组
func mergeFields(fieldsArr ...[]zapcore.Field) (fields []zapcore.Field) {
	for _, fs := range fieldsArr {
		fields = append(fields, fs...)
	}
	return
}
