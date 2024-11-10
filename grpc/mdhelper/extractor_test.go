package mdhelper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestExtractToOutCtx(t *testing.T) {
	table := []struct {
		name string
		kvs  metadata.MD
	}{
		{
			name: "three item",
			kvs: metadata.MD{
				DynamicLogLevelKey: {"debug"},
				WorkflowTracerOn:   {"true"},
				"other_key":        {"other_value"},
			},
		},
		{
			name: "one item",
			kvs: metadata.MD{
				"other_key": {"other_value"},
			},
		},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			keys := extractKeysFromMap(v.kvs)
			inCtx := metadata.NewIncomingContext(context.Background(), v.kvs)
			outCtx := ExtractToOutCtx(inCtx, WithExtractorKeys(keys))

			outMD, ok := metadata.FromOutgoingContext(outCtx)
			assert.True(t, ok)

			for key, value := range v.kvs {
				vs := outMD.Get(key)
				assert.Equal(t, len(vs), len(value))
				assert.Equal(t, vs, value)
			}
		})
	}
}

func extractKeysFromMap(kvs map[string][]string) []string {
	var keys []string
	for k := range kvs {
		keys = append(keys, k)
	}
	return keys
}
