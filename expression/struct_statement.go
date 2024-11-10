package expression

import (
	"fmt"
	"log"
)

type StructStatement struct {
	fields map[string]Expression
}

func NewStructStatement(fields map[string]Expression) *StructStatement {
	return &StructStatement{
		fields: fields,
	}
}

func (stmt *StructStatement) Dump(deep int) string {
	ret := fmt.Sprintf("%*sHEAD: StructStatement, name: %v", deep*4, "", stmt.fields)
	return ret
}

func (stmt *StructStatement) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	temp := map[string]*Data{}
	for n, f := range stmt.fields {
		d, err := f.Evaluate(ctx)
		log.Printf("%v eval => %v, %v\n", n, d, err)
		if nil != err {
			return nil, err
		}
		temp[n] = d
	}
	return NewStruct(temp), nil
}
