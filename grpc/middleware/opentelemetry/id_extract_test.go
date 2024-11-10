package grpcopentelemetry

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func TestInjectTraceIdsToOutCtx(t *testing.T) {
	scc := trace.SpanContextConfig{
		TraceID: trace.TraceID{1, 1},
		SpanID:  trace.SpanID{2, 2},
	}
	sc := trace.NewSpanContext(scc)
	ctx := trace.ContextWithSpanContext(context.Background(), sc)
	newCtx := InjectTraceIdsToOutCtx(ctx)
	md, ok := metadata.FromOutgoingContext(newCtx)
	assert.True(t, ok)
	traceParents := md.Get(traceParentHeader)
	assert.Equal(t, 1, len(traceParents))
	assert.Equal(t, fmt.Sprintf("00-%s-%s-00", hex.EncodeToString(scc.TraceID[:]), hex.EncodeToString(scc.SpanID[:])), traceParents[0])
}
