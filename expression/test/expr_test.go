package test

import (
	"strings"
	"testing"

	"github.com/niming-dev/ddd-demo/go-common/expression"
)

type EXPR_SEQ struct {
	expr   string
	expect *expression.Data
	ctx    expression.ExecuteContext
}

type testContext struct {
}

func (c *testContext) Get(name string) (*expression.Data, error) {
	return expression.NewString("zhangsan"), nil
}

func (c *testContext) Call(name string, args []*expression.Data) (*expression.Data, error) {
	return expression.NewString("lisi"), nil
}

var testExprSeq []EXPR_SEQ = []EXPR_SEQ{
	{expr: `"3" > 2`, expect: expression.NewBool(true)},
	{expr: `"3" + 2`, expect: expression.NewInt(5)},
	{expr: `6 + 3 * 2`, expect: expression.NewInt(12)},
	{expr: `"3" + "2"`, expect: expression.NewString("32")},
	{expr: `$(abc())`, ctx: &testContext{}, expect: expression.NewString("lisi")},
	{expr: `${xxxxx/skldjf}`, ctx: &testContext{}, expect: expression.NewString("zhangsan")},
}

func Test_A(t *testing.T) {

	for i, v := range testExprSeq {
		exec, err := expression.Parse(strings.NewReader(v.expr))
		if nil != err {
			t.Fatalf("testExprSeq[%d] %v", i, err)
		}

		d, err := exec.Evaluate(v.ctx)
		if nil != err {
			t.Fatalf("testExprSeq[%d] %v", i, err)
		}

		if b, err := expression.NotEqual(d, v.expect); nil != err {
			t.Fatalf("testExprSeq[%d] %v", i, err)
		} else {
			if b.Bool() {
				t.Fatalf("testExprSeq[%d] %v", i, "result not match")
			}
		}
	}
}
