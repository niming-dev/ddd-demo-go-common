package expression

import (
	"fmt"
	"reflect"
)

type MathStatement struct {
	args []Expression
	op   string
}

func NewMathStatement(args []Expression, op string) *MathStatement {
	return &MathStatement{args: args, op: op}
}

func (stmt *MathStatement) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: MathStatement, op: %v, args count: %v\n", deep*4, "", stmt.op, len(stmt.args))
	for _, v := range stmt.args {
		ret += fmt.Sprintf("%v\n", v.Dump(deep+1))
	}
	ret = ret[:len(ret)-1]
	return ret
}

func (stmt *MathStatement) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	var data0, data1 *Data
	if stmt.op == "neg" {
		if len(stmt.args) != 1 {
			return nil, ErrMissArgument
		}

		data0, err = data0.Evaluate(ctx)
		if nil == err {
			data, err = data.Neg()
			goto final_return
		}
	}

	if len(stmt.args) != 2 {
		err = ErrMissArgument
		goto final_return
	}
	data0, err = stmt.args[0].Evaluate(ctx)
	if nil != err {
		goto final_return
	}
	data1, err = stmt.args[1].Evaluate(ctx)
	if nil != err {
		goto final_return
	}
	switch stmt.op {
	case "+":
		data, err = Add(data0, data1)
	case "-":
		data, err = Sub(data0, data1)
	case "*":
		data, err = Mul(data0, data1)
	case "/":
		data, err = Div(data0, data1)
	case "%":
		data, err = Mod(data0, data1)
	case ">":
		data, err = Greate(data0, data1)
	case "<":
		data, err = Less(data0, data1)
	case ">=":
		data, err = GreateEqual(data0, data1)
	case "<=":
		data, err = LessEqual(data0, data1)
	case "==":
		data, err = Equal(data0, data1)
	case "!=":
		data, err = NotEqual(data0, data1)
	default:
		data, err = nil, ErrUnsupport
	}

final_return:
	if nil != ctx {
		typ := reflect.TypeOf(ctx)
		if typ.Implements(debuggerType) {
			ctx.(ExecuteContextDebugger).Log("MathStatement::Evaluate ->",
				data0, stmt.op, data1, "==>", data, err)
		}
	}
	return data, err
}
