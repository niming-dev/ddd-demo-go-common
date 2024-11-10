package storage

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"time"
)

func scan(v interface{}, dest interface{}) (err error) {
	// 收集panic信息
	defer func() {
		if r := recover(); nil != r {
			err = fmt.Errorf("error:[%v], stack:[%s]", r, string(debug.Stack()))
		}
	}()

	// 先使用类型断言方式赋值，类型不在匹配列表时使用反射赋值
	switch dest := dest.(type) {
	case nil:
		return DestCantBeNil
	case *string:
		if v, ok := v.(string); ok {
			*dest = v
			return nil
		}
	case *[]byte:
		if v, ok := v.([]byte); ok {
			*dest = v
			return nil
		}
	case *int:
		if v, ok := v.(int); ok {
			*dest = v
			return nil
		}
	case *int8:
		if v, ok := v.(int8); ok {
			*dest = v
			return nil
		}
	case *int16:
		if v, ok := v.(int16); ok {
			*dest = v
			return nil
		}
	case *int32:
		if v, ok := v.(int32); ok {
			*dest = v
			return nil
		}
	case *int64:
		if v, ok := v.(int64); ok {
			*dest = v
			return nil
		}
	case *uint:
		if v, ok := v.(uint); ok {
			*dest = v
			return nil
		}
	case *uint8:
		if v, ok := v.(uint8); ok {
			*dest = v
			return nil
		}
	case *uint16:
		if v, ok := v.(uint16); ok {
			*dest = v
			return nil
		}
	case *uint32:
		if v, ok := v.(uint32); ok {
			*dest = v
			return nil
		}
	case *uint64:
		if v, ok := v.(uint64); ok {
			*dest = v
			return nil
		}
	case *float32:
		if v, ok := v.(float32); ok {
			*dest = v
			return nil
		}
	case *float64:
		if v, ok := v.(float64); ok {
			*dest = v
			return nil
		}
	case *bool:
		if v, ok := v.(bool); ok {
			*dest = v
			return nil
		}
	case *[]string:
		if v, ok := v.([]string); ok {
			*dest = v
			return nil
		}
	case *[]int:
		if v, ok := v.([]int); ok {
			*dest = v
			return nil
		}
	case *[]bool:
		if v, ok := v.([]bool); ok {
			*dest = v
			return nil
		}
	case *map[string]string:
		if v, ok := v.(map[string]string); ok {
			*dest = v
			return nil
		}
	case *map[string]int:
		if v, ok := v.(map[string]int); ok {
			*dest = v
			return nil
		}
	case *map[string]bool:
		if v, ok := v.(map[string]bool); ok {
			*dest = v
			return nil
		}
	case *map[string]interface{}:
		if v, ok := v.(map[string]interface{}); ok {
			*dest = v
			return nil
		}
	case *map[string]struct{}:
		if v, ok := v.(map[string]struct{}); ok {
			*dest = v
			return nil
		}
	default:
		vDest := reflect.ValueOf(dest)
		if vDest.Kind() != reflect.Ptr {
			return DestMustBePtr
		}
		vDest = vDest.Elem()
		if !vDest.CanSet() {
			return DestUnsetable
		}

		// value为nil时把dest设置为对应类型的零值
		if v == nil {
			v := reflect.ValueOf(dest)
			v.Elem().Set(reflect.Zero(v.Elem().Type()))
			return nil
		}

		vValue := reflect.ValueOf(v)
		switch vValue.Kind() {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
			// value为nil时把dest设置为对应类型的零值
			if vValue.IsNil() {
				v := reflect.ValueOf(dest)
				v.Elem().Set(reflect.Zero(v.Elem().Type()))
				return nil
			}
		}

		// 当待设置的值为指针但目标对应不是指针时，取值的Elem
		if vValue.Kind() == reflect.Ptr && vDest.Kind() != reflect.Ptr {
			vValue = vValue.Elem()
		}

		if vDest.Type() != vValue.Type() {
			return TypeMismatch
		}

		vDest.Set(vValue)
		return nil
	}

	return TypeMismatch
}

func duration2Str(d time.Duration) string {
	return strconv.FormatInt(int64(d/time.Millisecond), 10) + "ms"
}
