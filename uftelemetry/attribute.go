package uftelemetry

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"go.opentelemetry.io/otel/attribute"
)

// Any2Attr 将任意值传为attribute.KeyValue
func Any2Attr(key string, value interface{}) attribute.KeyValue {
	switch v := value.(type) {
	case string:
		return attribute.String(key, v)
	case *string:
		return attribute.String(key, *v)
	case bool:
		return attribute.Bool(key, v)
	case *bool:
		return attribute.Bool(key, *v)
	case int:
		return attribute.Int(key, v)
	case *int:
		return attribute.Int(key, *v)
	case int32:
		return attribute.Int(key, int(v))
	case *int32:
		return attribute.Int(key, int(*v))
	case int64:
		return attribute.Int64(key, v)
	case *int64:
		return attribute.Int64(key, *v)
	case uint32:
		return attribute.Int(key, int(v))
	case *uint32:
		return attribute.Int(key, int(*v))
	case uint64:
		return attribute.Int64(key, int64(v))
	case *uint64:
		return attribute.Int64(key, int64(*v))
	case float32:
		return attribute.Float64(key, float64(v))
	case *float32:
		return attribute.Float64(key, float64(*v))
	case float64:
		return attribute.Float64(key, v)
	case *float64:
		return attribute.Float64(key, *v)
	case time.Duration:
		return attribute.String(key, strconv.Itoa(int(v))+"ns")
	case time.Time:
		return attribute.String(key, v.String())
	case error:
		return attribute.String(key, v.Error())
	case fmt.Stringer:
		return attribute.String(key, v.String())
	case attribute.MergeIterator:
		return v.Label()
	case json.Marshaler:
		bs, _ := v.MarshalJSON()
		return attribute.String(key, string(bs))
	case nil:
		return attribute.String(key, "<nil>")
	}

	// 对于接口包装过的nil使用反射判断
	if reflect2.IsNil(value) {
		return attribute.String(key, "<nil>")
	}

	// 其它格式转为json字符串
	str, _ := jsoniter.MarshalToString(value)
	return attribute.String(key, str)
}
