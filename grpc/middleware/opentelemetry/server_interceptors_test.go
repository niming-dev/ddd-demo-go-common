package grpcopentelemetry

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc/peer"
)

func Test_newServerSpanFromInbound(t *testing.T) {
	tp, err := _tracerProvider("http://localhost:14268/api/traces")
	assert.NoError(t, err)
	tr := tp.Tracer("component-main")

	ctx := _newTagsForCtx(context.Background())
	newServerSpanFromInbound(ctx, tr, "main")
}

func _tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("test"),
			attribute.String("environment", "env"),
			attribute.Int64("ID", 1),
		)),
	)
	return tp, nil
}

func _newTagsForCtx(ctx context.Context) context.Context {
	t := grpc_ctxtags.NewTags()
	if p, ok := peer.FromContext(ctx); ok {
		t.Set("peer.address", p.Addr.String())
	}
	return grpc_ctxtags.SetInContext(ctx, t)
}
