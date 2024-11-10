package ufgateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/niming-dev/ddd-demo/go-common/uftelemetry"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

// 网关定义的metadata key
const (
	GWProjectIdKey       = "x-uf-project-id"
	GWProjectNameKey     = "x-uf-project-name"
	GWAccountIdentityKey = "x-uf-account-identity"
	GWAccountNameKey     = "x-uf-account-name"
	GWWorkflowTracerKey  = "x-uf-workflow-tracer"
)

var (
	// GWKeys 网关http head与grpc metadata key的映射关系
	GWKeys = map[string]string{
		GWProjectIdKey:             GWProjectIdKey,
		GWProjectNameKey:           GWProjectNameKey,
		GWAccountIdentityKey:       GWAccountIdentityKey,
		GWAccountNameKey:           GWAccountNameKey,
		GWWorkflowTracerKey:        GWWorkflowTracerKey,
		uftelemetry.HeaderLevelKey: uftelemetry.HeaderLevelKey,
		// opentelemetry keys:  traceparent, tracestate
		propagation.TraceContext{}.Fields()[0]: propagation.TraceContext{}.Fields()[0],
		propagation.TraceContext{}.Fields()[1]: propagation.TraceContext{}.Fields()[1],
	}

	gwValidator = validator.New()
)

// GRPCGWKeyMapping gRPC网关的http header与metadata映射转换
func GRPCGWKeyMapping(key string) (string, bool) {
	if v, ok := GWKeys[strings.ToLower(key)]; ok {
		return v, true
	}
	return key, false
}

// GatewayData 网关传递过来的metadata数据
type GatewayData struct {
	ProjectId       string `validate:"required"`
	ProjectName     string `validate:"required"`
	AccountIdentity string `validate:"required"`
	AccountName     string `validate:"required"`
}

// CheckAll 检查所有数据
func (g GatewayData) CheckAll() error {
	return gwValidator.Struct(g)
}

// ExtractDataFromGRPC 从GRPC请求的context中提取网关数据
func ExtractDataFromGRPC(ctx context.Context) (GatewayData, bool) {
	gwData := GatewayData{}

	inMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return gwData, false
	}

	if ss := inMD.Get(GWProjectNameKey); len(ss) > 0 {
		gwData.ProjectName = ss[0]
	}
	if ss := inMD.Get(GWProjectIdKey); len(ss) > 0 {
		gwData.ProjectId = ss[0]
	}
	if ss := inMD.Get(GWAccountIdentityKey); len(ss) > 0 {
		gwData.AccountIdentity = ss[0]
	}
	if ss := inMD.Get(GWAccountNameKey); len(ss) > 0 {
		gwData.AccountName = ss[0]
	}

	return gwData, true
}

// ExtractDataFromHttp 从http请求中提取网关数据
func ExtractDataFromHttp(r *http.Request) (GatewayData, bool) {
	gwData := GatewayData{
		ProjectName:     r.Header.Get(GWProjectNameKey),
		ProjectId:       r.Header.Get(GWProjectIdKey),
		AccountIdentity: r.Header.Get(GWAccountIdentityKey),
		AccountName:     r.Header.Get(GWAccountNameKey),
	}

	return gwData, true
}
