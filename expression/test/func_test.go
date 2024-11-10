package test

import (
	"hash/adler32"
	"strings"
	"testing"

	"github.com/niming-dev/ddd-demo/go-common/expression"
	"github.com/niming-dev/ddd-demo/go-common/expression/function"
)

type funcContext struct {
}

func (c *funcContext) Get(name string) (*expression.Data, error) {
	return expression.NewString("zhangsan"), nil
}

func (c *funcContext) Call(name string, args []*expression.Data) (*expression.Data, error) {
	return function.Call(c, name, args)
}

func Test_func(t *testing.T) {
	exec, err := expression.Parse(strings.NewReader(`$(adler32("abcd"))`))
	if nil != err {
		t.Fatal(err)
	}
	i1 := int64(adler32.Checksum([]byte("abcd")))
	d1, err := exec.Evaluate(&funcContext{})
	if nil != err {
		t.Fatal(err)
	}
	if !d1.IsInt() {
		t.Fatal("expect got int")
	}
	if d1.Int() != i1 {
		t.Fatalf("expect got %v", i1)
	}

}
