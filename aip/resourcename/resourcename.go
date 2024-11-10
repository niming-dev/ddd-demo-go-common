package resourcename

import (
	"errors"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type Scanner interface {
	Scan(kvs map[string]string)
}

type MSScanner map[string]string

func (s MSScanner) Scan(kvs map[string]string) {
	for k, v := range kvs {
		s[k] = v
	}
}

// Scan 把resource对应的值写入到指定指针变量中
func Scan(name string, pattern string, dst any) error {
	nameArr := strings.Split(name, "/")
	patternArr := strings.Split(pattern, "/")

	if len(nameArr) != len(patternArr) {
		return errors.New("invalid pattern")
	}

	result := map[string]string{}
	for i, v := range patternArr {
		if vLen := len(v); vLen > 2 && v[0] == '{' && v[vLen-1] == '}' {
			n := v[1 : len(v)-1]
			result[n] = nameArr[i]
		}
	}

	switch dstV := dst.(type) {
	case map[string]string:
		for k := range dstV {
			delete(dstV, k)
		}
		for k, v := range result {
			dstV[k] = v
		}
		return nil
	case *map[string]string:
		*dstV = result
		return nil
	case Scanner:
		dstV.Scan(result)
		return nil
	default:
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			ZeroFields:           true,
			WeaklyTypedInput:     true,
			Result:               dst,
			TagName:              "resourcename",
			IgnoreUntaggedFields: false,
			MatchName:            nil,
		})
		if err != nil {
			return err
		}
		return decoder.Decode(result)
	}
}
