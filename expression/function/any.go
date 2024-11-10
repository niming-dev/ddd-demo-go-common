package function

import "github.com/niming-dev/ddd-demo/go-common/expression"

type any struct{}

func (any) Name() string { return "any" }

func (any) Call(ctx expression.ExecuteContext, args []*expression.Data) (*expression.Data, error) {
	return expression.NewString("*"), nil
}

func init() {
	Register(&any{})
}
