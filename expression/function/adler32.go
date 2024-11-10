package function

import (
	a "hash/adler32"

	"github.com/niming-dev/ddd-demo/go-common/expression"
)

type adler32 struct{}

func (adler32) Name() string { return "adler32" }

func (adler32) Call(ctx expression.ExecuteContext, args []*expression.Data) (*expression.Data, error) {
	if len(args) == 0 {
		return nil, expression.ErrMissArgument
	}
	if !args[0].IsString() {
		return nil, expression.ErrDataTypeNotMatch
	}

	return expression.NewInt(int64(a.Checksum([]byte(args[0].String())))), nil
}

func init() {
	Register(&adler32{})
}
