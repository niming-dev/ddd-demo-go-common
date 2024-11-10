package test

import (
	"testing"

	"github.com/niming-dev/ddd-demo/go-common/expression"
)

type opFunc func(o1, o2 *expression.Data) (*expression.Data, error)
type DATA_SEQ struct {
	o1     *expression.Data
	o2     *expression.Data
	op     opFunc
	expect *expression.Data
}

var testDataSeq []DATA_SEQ = []DATA_SEQ{
	{expression.NewInt(5), expression.NewInt(3), expression.Add, expression.NewInt(8)},
	{expression.NewInt(3), expression.NewInt(5), expression.Add, expression.NewInt(8)},
	{expression.NewInt(5), expression.NewInt(3), expression.Sub, expression.NewInt(2)},
	{expression.NewInt(3), expression.NewInt(5), expression.Sub, expression.NewInt(-2)},
	{expression.NewInt(3), expression.NewInt(5), expression.Mul, expression.NewInt(15)},
	{expression.NewInt(15), expression.NewInt(3), expression.Div, expression.NewInt(5)},
	{expression.NewInt(15), expression.NewInt(16), expression.Div, expression.NewInt(0)},
	{expression.NewInt(15), expression.NewInt(0), expression.Div, nil},
	{expression.NewInt(15), expression.NewInt(3), expression.Mod, expression.NewInt(0)},
	{expression.NewInt(15), expression.NewInt(4), expression.Mod, expression.NewInt(3)},
	{expression.NewInt(15), expression.NewInt(16), expression.Mod, expression.NewInt(15)},
	{expression.NewInt(15), expression.NewInt(0), expression.Mod, nil},
	{expression.NewString("15"), expression.NewInt(3), expression.Add, expression.NewInt(18)},
	{expression.NewBool(true), expression.NewInt(3), expression.Add, expression.NewInt(4)},
	{expression.NewBool(false), expression.NewInt(3), expression.Add, expression.NewInt(3)},
	{expression.NewBool(false), expression.NewBool(true), expression.Add, nil},
	{expression.NewBool(false), expression.NewBool(true), expression.Mul, nil},
	{expression.NewString("15"), expression.NewString("nihao"), expression.Mul, nil},
	{expression.NewInt(15), expression.NewString("nihao"), expression.Mul, nil},
	{expression.NewInt(15), expression.NewString("15"), expression.Mul, expression.NewInt(225)},
	{expression.NewString("15"), expression.NewString("15"), expression.Add, expression.NewString("1515")},
	{expression.NewString("15"), expression.NewString("16"), expression.Sub, nil},
	{expression.NewString("115"), expression.NewString("16"), expression.Greate, expression.NewBool(false)},
	{expression.NewString("115"), expression.NewInt(16), expression.Greate, expression.NewBool(true)},
	{expression.NewString("115"), expression.NewString("16"), expression.Less, expression.NewBool(true)},
	{expression.NewString("115"), expression.NewInt(16), expression.Less, expression.NewBool(false)},
	{expression.NewString("16"), expression.NewInt(16), expression.GreateEqual, expression.NewBool(true)},
	{expression.NewString("16"), expression.NewInt(12), expression.GreateEqual, expression.NewBool(true)},
	{expression.NewString("16"), expression.NewInt(16), expression.LessEqual, expression.NewBool(true)},
	{expression.NewString("16"), expression.NewInt(12), expression.LessEqual, expression.NewBool(false)},
	{expression.NewString("16"), expression.NewInt(16), expression.Equal, expression.NewBool(true)},
	{expression.NewString("16"), expression.NewInt(12), expression.Equal, expression.NewBool(false)},
	{expression.NewString("16"), expression.NewInt(16), expression.NotEqual, expression.NewBool(false)},
	{expression.NewString("16"), expression.NewInt(12), expression.NotEqual, expression.NewBool(true)},
	{expression.NewBool(true), expression.NewInt(1), expression.Equal, expression.NewBool(true)},
	{expression.NewBool(true), expression.NewBool(true), expression.Equal, expression.NewBool(true)},
	{expression.NewBool(true), expression.NewBool(false), expression.Equal, expression.NewBool(false)},
	{expression.NewBool(true), expression.NewBool(false), expression.Greate, nil},
	{expression.NewString("sskldfj"), expression.NewBool(false), expression.Greate, nil},
	{expression.NewString("true"), expression.NewBool(false), expression.Greate, nil},
	{expression.NewString("true"), expression.NewBool(false), expression.Equal, expression.NewBool(false)},
	{expression.NewString("false"), expression.NewBool(false), expression.Equal, expression.NewBool(true)},
}

func Test_math(t *testing.T) {
	for i, v := range testDataSeq {
		ret, err := v.op(v.o1, v.o2)
		if nil != err {
			if v.expect == nil {
				// it's ok
			} else {
				t.Fatalf("testDataSeq[%d] failed, %v, %v %v\n", i, err, v.o1.Dump(0), v.o2.Dump(0))
			}
		} else {
			if nil == v.expect {
				t.Fatalf("testDataSeq[%d] expect got error", i)
			}
			neret, err := expression.Equal(ret, v.expect)
			if nil != err {
				t.Fatalf("testDataSeq[%d] Equal failed\n", i)
			}
			if !neret.Bool() {
				t.Fatalf("testDataSeq[%d] failed, not equal, got %v, expect %v\n", i, ret, v.expect)
			}
		}
	}
}
