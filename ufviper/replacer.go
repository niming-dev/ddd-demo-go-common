package ufviper

import (
	"strings"
)

// LowCaseReplace 把字符串替换为小写
type LowCaseReplace struct{}

func (v LowCaseReplace) Replace(s string) string {
	return strings.ToLower(s)
}
