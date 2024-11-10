package expression

import (
	"fmt"
	"log"
)

const (
	SWITCH_VALUE_NAME = "switchValue_for_internal_user"
)

type SwitchStatement struct {
	val   Expression
	stmts []Expression
}

type SwitchContext struct {
	originCtx ExecuteContext
	val       *Data
}

func (ctx *SwitchContext) Call(name string, args []*Data) (*Data, error) {
	return ctx.originCtx.Call(name, args)
}

func (ctx *SwitchContext) Get(name string) (*Data, error) {
	if name == SWITCH_VALUE_NAME {
		return ctx.val, nil
	}
	return ctx.originCtx.Get(name)
}

func NewSwitchStatement(val Expression, stmts []Expression) *SwitchStatement {
	return &SwitchStatement{
		val:   val,
		stmts: stmts,
	}
}

func (stmt *SwitchStatement) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: SwitchStatement, statements count: %d\n", deep*4, "", len(stmt.stmts))
	ret += fmt.Sprintf("%v\n", stmt.val.Dump(deep+1))
	for _, v := range stmt.stmts {
		ret += fmt.Sprintf("%v\n", v.Dump(deep+1))
	}
	// trim last \n
	ret = ret[:len(ret)-1]
	return ret
}

func (stmt *SwitchStatement) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	log.Println("SwitchStatement::Evaluate")
	val, err := stmt.val.Evaluate(ctx)
	if nil != err {
		return nil, err
	}

	switchCtx := &SwitchContext{originCtx: ctx, val: val}
	for _, s := range stmt.stmts {
		data, err = s.Evaluate(switchCtx)
		if nil != err {
			if err == ErrReturn {
				return data, err
			} else {
				return nil, err
			}
		}
	}
	return nil, nil
}
