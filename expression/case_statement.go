package expression

import (
	"fmt"
	"log"
)

type CaseStatement struct {
	cond  Expression
	stmts []Expression
}

func NewCaseStatement(cond Expression, stmts []Expression) *CaseStatement {
	return &CaseStatement{cond: cond, stmts: stmts}
}

func (stmt *CaseStatement) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	log.Println(stmt.cond.Dump(0))
	cond, err := stmt.cond.Evaluate(ctx)
	if nil != err {
		return nil, err
	}
	log.Println(cond, err)
	if !cond.Bool() {
		return nil, nil
	}
	for _, s := range stmt.stmts {
		data, err = s.Evaluate(ctx)
		if nil != err {
			if err == ErrReturn {
				return data, err
			} else {
				return nil, err
			}
		}
	}
	return nil, ErrNoReturn
}

func (stmt *CaseStatement) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: CaseStatement\n", deep*4, "")
	ret += stmt.cond.Dump(deep+1) + "\n"
	for _, s := range stmt.stmts {
		ret += fmt.Sprintf("%v\n", s.Dump(deep+1))
	}
	ret = ret[:len(ret)-1]
	return ret
}
