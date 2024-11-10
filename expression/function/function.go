package function

import (
	"fmt"
	"sync"

	"github.com/niming-dev/ddd-demo/go-common/expression"
)

type Function interface {
	Name() string
	Call(expression.ExecuteContext, []*expression.Data) (*expression.Data, error)
}

var (
	funcs = map[string]Function{}
	m     sync.Mutex
)

func Register(f Function) {
	m.Lock()
	defer m.Unlock()
	_, ok := funcs[f.Name()]
	if ok {
		panic(fmt.Sprintf("function %s exists", f.Name()))
	}
	funcs[f.Name()] = f
}

func Get(name string) Function {
	m.Lock()
	defer m.Unlock()
	f, ok := funcs[name]
	if !ok {
		return nil
	} else {
		return f
	}
}

func Call(ctx expression.ExecuteContext, name string, args []*expression.Data) (*expression.Data, error) {
	f := Get(name)
	if nil == f {
		return nil, expression.ErrUnsupport
	}
	return f.Call(ctx, args)
}
