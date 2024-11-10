package expression

import (
	"fmt"
	"io"
	"sync"
)

func SetDebugLevel(l int) {
	yyDebug = l
}

var currentExecutor *Executor
var globalMutex *sync.Mutex = &sync.Mutex{}

var (
	ErrSyntaxError    = fmt.Errorf("syntax error")
	ErrNoExpression   = fmt.Errorf("no expression to be execute")
	ErrMissArgument   = fmt.Errorf("miss required arguments")
	ErrUnsupport      = fmt.Errorf("unsupport operator")
	ErrInvalidConvert = fmt.Errorf("invalid convert")
	ErrReturn         = fmt.Errorf("return") // 该错误并不是错误，只是用来打断后续语句的执行
	ErrNoReturn       = fmt.Errorf("no return")
)

func Parse(input io.Reader) (exec *Executor, err error) {
	globalMutex.Lock()
	currentExecutor = &Executor{}
	defer func() {
		if r := recover(); nil != r {
			err = r.(error)
			globalMutex.Unlock()
		}
	}()
	yyParse(NewLexer(input))
	exec = currentExecutor
	globalMutex.Unlock()
	return exec, nil
}
