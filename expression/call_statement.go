package expression

import (
	"fmt"
	"reflect"
)

type CallStatement struct {
	args []Expression
	name string
}

func NewCallStatement(args []Expression, name string) *CallStatement {
	return &CallStatement{args: args, name: name}
}

func (stmt *CallStatement) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	var args []*Data

	for _, v := range stmt.args {
		d, err := v.Evaluate(ctx)
		if nil != err {
			goto final_return
		}
		args = append(args, d)
	}

	data, err = ctx.Call(stmt.name, args)

final_return:
	typ := reflect.TypeOf(ctx)
	if typ.Implements(debuggerType) {
		ctx.(ExecuteContextDebugger).Log("CallStatement::Evaluate ->",
			stmt.name, "(", args, ") ==>", data, err)
	}
	return data, err
}

func (stmt *CallStatement) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: CallStatement, name: %v, args count: %v\n", deep*4, "", stmt.name, len(stmt.args))
	for _, v := range stmt.args {
		ret += fmt.Sprintf("%v\n", v.Dump(deep+1))
	}
	// trim last \n
	ret = ret[:len(ret)-1]
	return ret
}
