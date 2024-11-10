package test

import (
	"strings"
	"testing"

	"github.com/niming-dev/ddd-demo/go-common/expression"
)

func Test_Parse(t *testing.T) {
	exec, err := expression.Parse(strings.NewReader(`$(fetch("GET", "http://www.baidu.com", 
		{a: "abcd", b : 3, c: ${getval}}))`))
	if nil != err {
		t.Fatal(err)
	}
	ctx := &testContext{}
	t.Log(exec.Evaluate(ctx))
}
