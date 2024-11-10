package expression

import (
	"fmt"
	"log"
)

type ReturnStatement struct {
	expr Expression
}

func NewReturnStatement(expr Expression) *ReturnStatement {
	return &ReturnStatement{expr: expr}
}

func (stmt *ReturnStatement) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	log.Println("ReturnStatement::Evaluate")
	data, err = stmt.expr.Evaluate(ctx)
	if nil != err && ErrReturn != err {
		return nil, err
	}
	return data, ErrReturn
}

func (stmt *ReturnStatement) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: ReturnStatement\n", deep*4, "")
	ret += stmt.expr.Dump(deep + 1)
	return ret
}
