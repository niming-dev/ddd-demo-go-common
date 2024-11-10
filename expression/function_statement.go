package expression

import (
	"fmt"
	"log"
)

type FunctionStatement struct {
	stmts []Expression
}

func NewFunctionStatement(stmts []Expression) *FunctionStatement {
	return &FunctionStatement{stmts: stmts}
}

func (stmt *FunctionStatement) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	log.Println("count: ", len(stmt.stmts))
	for _, s := range stmt.stmts {
		data, err = s.Evaluate(ctx)
		log.Println(data, err)
		if nil != err {
			if err == ErrReturn {
				return data, nil
			} else {
				return nil, err
			}
		}
	}

	return nil, ErrNoReturn
}

func (stmt *FunctionStatement) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: FunctionStatement, statements count: %d\n", deep*4, "", len(stmt.stmts))
	for _, v := range stmt.stmts {
		ret += fmt.Sprintf("%v\n", v.Dump(deep+1))
	}
	// trim last \n
	ret = ret[:len(ret)-1]
	return ret
}
