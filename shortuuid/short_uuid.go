package shortuuid

import (
	"github.com/google/uuid"

	"github.com/niming-dev/ddd-demo/go-common/intconv"
)

// New 通过UUID数据生成一个较短的ID字符串
func New() string {
	uuidBinary, _ := uuid.New().MarshalBinary()

	return intconv.BytesToBase62(uuidBinary)
}
