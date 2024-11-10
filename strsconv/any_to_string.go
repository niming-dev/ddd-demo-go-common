package strsconv

import (
	json "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

type any2StringOptions struct {
	objectEncoder func(v interface{}) string
	nilValue      string
}

var defaultAny2StringOptions = &any2StringOptions{
	objectEncoder: defaultObjEncoder,
	nilValue:      "<nil>",
}

type Any2StringOption func(*any2StringOptions)

func evaluateOpt(opts []Any2StringOption) *any2StringOptions {
	optCopy := &any2StringOptions{}
	*optCopy = *defaultAny2StringOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// ObjectEncoder 对象转字符串编码器
type ObjectEncoder func(v interface{}) string

var defaultObjEncoder = func(v interface{}) string {
	str, err := json.MarshalToString(v)
	if err != nil {
		return ""
	}
	return str
}

// Any2String 把任意类型转为字符串
func Any2String(v interface{}, options ...Any2StringOption) string {
	o := evaluateOpt(options)
	switch v := v.(type) {
	case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64,
		*string, *bool, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64:
		return Simple2String(v)
	case error:
		return v.Error()
	case nil:
		return o.nilValue
	}

	// 对于接口包装过的nil使用反射判断
	if reflect2.IsNil(v) {
		return o.nilValue
	}

	return o.objectEncoder(v)
}
