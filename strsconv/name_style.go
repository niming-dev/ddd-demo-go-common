package strsconv

import (
	"bytes"
	"strings"
)

const (
	// ASCIIMaxLower ASCII范围内最大的小写字母
	ASCIIMaxLower = '\u007A'
	// ASCIIMinLower ASCII范围内最小的小写字母
	ASCIIMinLower = '\u0061'
	// ASCIIMaxUpper ASCII范围内最大的大写字母
	ASCIIMaxUpper = '\u005A'
	// ASCIIMinUpper ASCII范围内最小的大写字母
	ASCIIMinUpper = '\u0041'
	// ASCIICaseDifferenceValue ASCII范围内大写与小写字母的差值
	ASCIICaseDifferenceValue = 32
)

// IsASCIIUpper 判断是否为ASCII范围内的大写字母
func IsASCIIUpper(r rune) bool {
	return ASCIIMinUpper <= r && r <= ASCIIMaxUpper
}

// IsASCIILower 判断是否为ASCII范围内的小写字母
func IsASCIILower(r rune) bool {
	return ASCIIMinLower <= r && r <= ASCIIMaxLower
}

// ASCII2Lower 把ASCII范围内的字母转为小写
func ASCII2Lower(r rune) rune {
	if IsASCIIUpper(r) {
		return r + ASCIICaseDifferenceValue
	}

	return r
}

// ASCII2Upper 把ASCII范围内的字母转为大写
func ASCII2Upper(r rune) rune {
	if IsASCIILower(r) {
		return r - ASCIICaseDifferenceValue
	}

	return r
}

// Camel2Snake 驼峰命名转蛇形
func Camel2Snake(name string) string {
	buffer := new(bytes.Buffer)
	for i, r := range name {
		if IsASCIIUpper(r) {
			if i != 0 {
				buffer.WriteByte('_')
			}
			buffer.WriteRune(ASCII2Lower(r))
		} else {
			buffer.WriteRune(r)
		}
	}

	return buffer.String()
}

// Pascal2Snake 帕斯卡命名转蛇形
func Pascal2Snake(name string) string {
	return Camel2Snake(name)
}

// Snake2Pascal 蛇形命名转帕斯卡
func Snake2Pascal(name string) string {
	buffer := new(bytes.Buffer)
	underline := false
	for i, r := range name {
		if i == 0 {
			buffer.WriteRune(ASCII2Upper(r))
			continue
		}

		if r == '_' {
			underline = true
			continue
		}

		if underline {
			r = ASCII2Upper(r)
			underline = false
		}

		buffer.WriteRune(r)
	}

	return buffer.String()
}

// Camel2Pascal 驼峰命名转帕斯卡
func Camel2Pascal(name string) string {
	if name == "" {
		return ""
	}

	nameRunes := []rune(name)
	buffer := new(bytes.Buffer)
	if IsASCIILower(nameRunes[0]) {
		buffer.WriteRune(ASCII2Upper(nameRunes[0]))
	}

	buffer.WriteString(string(nameRunes[1:]))
	return buffer.String()
}

// Snake2Camel 蛇形命名转驼峰
func Snake2Camel(name string) string {
	buffer := new(bytes.Buffer)
	underline := false
	for i, r := range name {
		if i == 0 {
			buffer.WriteRune(r)
			continue
		}

		if r == '_' {
			underline = true
			continue
		}

		if underline {
			r = ASCII2Upper(r)
			underline = false
		}

		buffer.WriteRune(r)
	}

	return buffer.String()
}

// Pascal2Camel 帕斯卡命名转驼峰
func Pascal2Camel(name string) string {
	l := len(name)
	if l == 0 {
		return ""
	}
	if l == 1 {
		return strings.ToLower(name[0:1])
	}

	return strings.ToLower(name[0:1]) + name[1:]
}

type NameStyle string

const (
	NameStyleUnspecified NameStyle = "NAME_STYLE_UNSPECIFIED"
	NameStylePascal      NameStyle = "PASCAL"
	NameStyleCamel       NameStyle = "CAMEL"
	NameStyleSnake       NameStyle = "SNAKE"
)

func (s NameStyle) IsValid() bool {
	_, ok := NameStyleValue[string(s)]
	return ok
}

func (s NameStyle) String() string {
	return string(s)
}

var NameStyleValue = map[string]NameStyle{
	"NAME_STYLE_UNSPECIFIED": NameStyleUnspecified,
	"PASCAL":                 NameStylePascal,
	"CAMEL":                  NameStyleCamel,
	"SNAKE":                  NameStyleSnake,
}

var convertFunc = map[string]func(name string) string{
	"CAMEL_SNAKE":  Camel2Snake,
	"PASCAL_SNAKE": Pascal2Snake,
	"SNAKE_PASCAL": Snake2Pascal,
	"CAMEL_PASCAL": Camel2Pascal,
	"SNAKE_CAMEL":  Snake2Camel,
	"PASCAL_CAMEL": Pascal2Camel,
}

// ConvertNameStyle 转换命名风格
func ConvertNameStyle(name string, src NameStyle, dest NameStyle) string {
	if f, ok := convertFunc[string(src)+"_"+string(dest)]; ok {
		return f(name)
	}
	return name
}

// NameStyleAnyToOne 任何风格名称转指定风格
func NameStyleAnyToOne(name string, dest NameStyle) string {
	if !dest.IsValid() {
		return name
	}

	return ConvertNameStyle(name, CheckNameStyle(name), dest)
}

// CheckNameStyle 检测命名风格
func CheckNameStyle(name string) NameStyle {
	if name == "" {
		return NameStyleUnspecified
	}

	if strings.ContainsRune(name, '_') {
		return NameStyleSnake
	}

	if IsASCIIUpper([]rune(name)[0]) {
		return NameStylePascal
	}

	return NameStyleCamel
}
