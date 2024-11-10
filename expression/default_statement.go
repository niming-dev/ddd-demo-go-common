package expression

import (
	"fmt"
	"log"
)

type DefaultStatement struct {
	stmts []Expression
}

func NewDefaultStatement(stmts []Expression) *DefaultStatement {
	return &DefaultStatement{stmts: stmts}
}

func (stmt *DefaultStatement) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	log.Println("DefaultStatement::Evaluate")
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

func (stmt *DefaultStatement) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: DefaultStatement\n", deep*4, "")
	for _, s := range stmt.stmts {
		ret += fmt.Sprintf("%v\n", s.Dump(deep+1))
	}
	ret = ret[:len(ret)-1]
	return ret
}
