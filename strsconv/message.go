package strsconv

import (
	"bytes"
)

// Camel2Message 驼峰命名转为消息字符串
func Camel2Message(name string) string {
	buffer := new(bytes.Buffer)
	for i, r := range name {
		if IsASCIIUpper(r) {
			if i == 0 {
				buffer.WriteRune(r)
				continue
			}

			buffer.WriteByte(' ')
			buffer.WriteRune(ASCII2Lower(r))
			continue
		}

		buffer.WriteRune(r)
	}

	return buffer.String()
}
