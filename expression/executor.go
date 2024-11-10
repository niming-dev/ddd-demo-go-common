package expression

import (
	"fmt"
	"reflect"
)

type ExecuteContext interface {
	Call(name string, args []*Data) (*Data, error)
	Get(name string) (*Data, error)
}

type ExecuteContextDebugger interface {
	Log(args ...interface{})
}

var debuggerType = reflect.TypeOf((*ExecuteContextDebugger)(nil)).Elem()

type Executor struct {
	stack []Expression
}

// *************************************************
// 编译期间调用的函数 v
func (exec *Executor) Push(elem Expression) {
	exec.stack = append(exec.stack, elem)
}

func (exec *Executor) Pop() (Expression, bool) {
	if nil == exec.stack || len(exec.stack) == 0 {
		return nil, false
	}
	ret := exec.stack[len(exec.stack)-1]
	exec.stack = exec.stack[:len(exec.stack)-1]
	return ret, true
}

// 编译期间调用的函数 ^
// *************************************************

func (exec *Executor) Dump(deep int) string {
	ret := ""
	ret = fmt.Sprintf("%*sHEAD: Executor, stack count: %v\n", 4*deep, "", len(exec.stack))
	for _, v := range exec.stack {
		ret += fmt.Sprintf("%v\n", v.Dump(deep+1))
	}
	return ret
}

func (exec *Executor) Evaluate(ctx ExecuteContext) (data *Data, err error) {
	// 目前只可能有一个语句
	if len(exec.stack) < 1 {
		data, err = nil, ErrNoExpression
	} else {
		data, err = exec.stack[0].Evaluate(ctx)
	}
	if nil != ctx {
		typ := reflect.TypeOf(ctx)
		if typ.Implements(debuggerType) {
			ctx.(ExecuteContextDebugger).Log("Executor::Evaluate ->", data, err)
		}
	}
	return
}
