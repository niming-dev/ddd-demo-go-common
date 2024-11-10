package expression

import (
	"fmt"
	"reflect"
)

type VariableStatement struct {
	name string
}

func NewVariableStatement(name string) *VariableStatement {
	return &VariableStatement{name: name}
}

func (stmt *VariableStatement) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	data, err = ctx.Get(stmt.name)

	typ := reflect.TypeOf(ctx)
	if typ.Implements(debuggerType) {
		ctx.(ExecuteContextDebugger).Log("VariableStatement::Evaluate ->",
			stmt.name, "==>", data, err)
	}
	return
}

func (stmt *VariableStatement) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: VariableStatement, name: %v", deep*4, "", stmt.name)
	return ret
}
