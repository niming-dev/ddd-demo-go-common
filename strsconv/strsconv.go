package strsconv

import (
	"strconv"
)

// Simple2String convert simple type to string
func Simple2String(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case *string:
		return *v
	case *bool:
		return strconv.FormatBool(*v)
	case *int:
		return strconv.FormatInt(int64(*v), 10)
	case *int8:
		return strconv.FormatInt(int64(*v), 10)
	case *int16:
		return strconv.FormatInt(int64(*v), 10)
	case *int32:
		return strconv.FormatInt(int64(*v), 10)
	case *int64:
		return strconv.FormatInt(*v, 10)
	case *uint:
		return strconv.FormatInt(int64(*v), 10)
	case *uint8:
		return strconv.FormatInt(int64(*v), 10)
	case *uint16:
		return strconv.FormatInt(int64(*v), 10)
	case *uint32:
		return strconv.FormatInt(int64(*v), 10)
	case *uint64:
		return strconv.FormatInt(int64(*v), 10)
	case *float32:
		return strconv.FormatFloat(float64(*v), 'f', -1, 64)
	case *float64:
		return strconv.FormatFloat(*v, 'f', -1, 64)
	}

	return ""
}
