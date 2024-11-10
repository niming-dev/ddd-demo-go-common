package expression

import (
	"fmt"
	"log"
)

type CaseExpression struct {
	op   string
	expr Expression
}

func NewCaseExpression(op string, expr Expression) *CaseExpression {
	return &CaseExpression{op: op, expr: expr}
}

func (stmt *CaseExpression) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	log.Println("CaseExpression::Evaluate")
	val, err := ctx.Get(SWITCH_VALUE_NAME)
	if nil != err {
		return nil, err
	}
	log.Printf("val = %v\n", val.Dump(0))

	data, err = stmt.expr.Evaluate(ctx)
	if nil != err {
		return nil, err
	}

	log.Printf("data = %v\n", data.Dump(0))

	if len(stmt.op) > 0 {
		mathStmt := NewMathStatement([]Expression{val, data}, stmt.op)
		return mathStmt.Evaluate(ctx)
	} else {
		// 没有符号的直接判断是否相等
		return Equal(val, data)
	}
}

func (stmt *CaseExpression) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: CaseExpression, op: %v\n", deep*4, "", stmt.op)
	ret += stmt.expr.Dump(deep + 1)
	return ret
}
