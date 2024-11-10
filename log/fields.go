package log

import (
	"context"
)

const (
	FieldKeyModule       = "module"
	FieldKeyTraceID      = "traceID"
	FieldKeySpanID       = "spanID"
	FieldKeyParentSpanID = "parent.spanID"
	FieldKeyElapsed      = "elapsed"
)

type Fields map[string]interface{}

func ExtractLogFields(ctx context.Context) Fields {
	fields := Fields{}
	if ctx != nil {
		if ctxFields, ok := ctx.Value(contextKeyLogFields).(Fields); ok {
			for k, v := range ctxFields {
				fields[k] = v
			}
		}
	}
	return fields
}
