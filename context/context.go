package context

import (
	"context"

	"github.com/niming-dev/ddd-demo/go-common/log"
)

func WithLogField(parent context.Context, key string, value interface{}) context.Context {
	return log.WithFieldContext(parent, key, value)
}

func WithLogFields(parent context.Context, fields log.Fields) context.Context {
	return log.WithFieldsContext(parent, fields)
}
