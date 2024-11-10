package uferror

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"log"
)

type UFError interface {
	WithMessage(v ...interface{}) UFError
	WithRequestId(requestId string) UFError
	// 缺省情况下只会产生一级调用的调用栈，并且只skip = 1
	WithDebugInfo(debugInfo string) UFError
	// 缺省情况下只会产生一级调用的调用栈，并且只skip = 1
	WithDebugInfof(format string, args ...interface{}) UFError
	// 指定调用栈的产生方法
	WithCallstack(skip int, deep int) UFError
	WithSuggestHttpCode(httpCode int) UFError
	Code() int
	Message() string
	RequestId() string
	DebugInfo() string
	SuggestHttpCode() int
	Callstack() []string
	ToError() error
	ToDebugError() error
	// 根据code判断两个错误是否是相同类型，忽略消息内容和请求id等
	IsKindOf(uferr UFError) bool
	// 在调用处重新设置CallStack
	// 并在基础错误之后增加args...的错误信息
	New(args ...interface{}) UFError
	Newf(fmt string, args ...interface{}) UFError
	error
	// String和GoString方法会返回全部的信息，默认fmt格式化会使用Error() string
	// 如需打印详细信息请显示调用String()或者GoString()或者 fmt.Printf("%#v", uferr)
	fmt.Stringer
	fmt.GoStringer
}

type UFErrorStruct struct {
	message         string
	code            int
	requestId       string
	suggestHttpCode int
	debugInfo       string
	callstack       []string
}

func New(code int, msg string) UFError {
	return &UFErrorStruct{
		code:    code,
		message: msg,
	}
}

func Newf(code int, format string, args ...interface{}) UFError {
	return &UFErrorStruct{
		code:    code,
		message: fmt.Sprintf(format, args...),
	}
}

func (e *UFErrorStruct) WithMessage(v ...interface{}) UFError {
	e.message = fmt.Sprint(v...)
	return e
}

func (e *UFErrorStruct) WithRequestId(requestId string) UFError {
	e.requestId = requestId
	return e
}

// 缺省情况下只会产生一级调用的调用栈，并且只skip = 1
func (e *UFErrorStruct) WithDebugInfo(debugInfo string) UFError {
	e.debugInfo = debugInfo
	e.WithCallstack(3, 20)
	return e
}

func (e *UFErrorStruct) WithDebugInfof(format string, args ...interface{}) UFError {
	e.debugInfo = fmt.Sprintf(format, args...)
	e.WithCallstack(3, 20)
	return e
}

func (e *UFErrorStruct) WithCallstack(skip int, deep int) UFError {
	var pc []uintptr = make([]uintptr, deep)
	n := runtime.Callers(skip, pc)

	e.callstack = []string{}
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()

		e.callstack = append(e.callstack, fmt.Sprintf("%s:%d", frame.File, frame.Line))
		if !more {
			break
		}
	}
	return e
}

func (e *UFErrorStruct) WithSuggestHttpCode(httpCode int) UFError {
	e.suggestHttpCode = httpCode
	return e
}

func (e *UFErrorStruct) Code() int {
	return e.code
}

func (e *UFErrorStruct) Message() string {
	return e.message
}

func (e *UFErrorStruct) RequestId() string {
	return e.requestId
}

func (e *UFErrorStruct) DebugInfo() string {
	return e.debugInfo
}

func (e *UFErrorStruct) SuggestHttpCode() int {
	return e.suggestHttpCode
}

func (e *UFErrorStruct) Callstack() []string {
	return e.callstack
}

func str_or_empty(key string, value interface{}) string {
	switch vt := value.(type) {
	case int:
		if value.(int) != 0 {
			return fmt.Sprintf(", %s: %d", key, value.(int))
		}
	case string:
		if len(value.(string)) > 0 {
			return fmt.Sprintf(", %s: %s", key, value.(string))
		}
	case []string:
		if len(value.([]string)) > 0 {
			jsonstr, _ := json.Marshal(value.([]string))
			return fmt.Sprintf(", %s: %s", key, string(jsonstr))
		}
	default:
		log.Printf("Unhandled type %v, return empty string", vt)
		return ""
	}
	return ""
}

func (e *UFErrorStruct) ToError() error {
	return errors.New(e.Error())
}

func (e UFErrorStruct) GoString() string {
	return e.String()
}

func (e UFErrorStruct) String() string {
	return fmt.Sprintf("code: %d, message: %s%s%s%s%s", e.code, e.message,
		str_or_empty("requestId", e.requestId),
		str_or_empty("suggestHttpCode", e.suggestHttpCode),
		str_or_empty("debugInfo", e.debugInfo),
		str_or_empty("callstack", e.callstack),
	)
}
func (e *UFErrorStruct) ToDebugError() error {
	log.Println("223")
	return errors.New(e.String())
}

func (e *UFErrorStruct) Error() string {
	// return e.String()
	return fmt.Sprintf("code: %d, message: %s%s", e.code, e.message, str_or_empty("requestId", e.requestId))
}

// 根据code判断两个错误是否是相同类型，忽略消息内容和请求id等
func (e UFErrorStruct) IsKindOf(uferr UFError) bool {
	return e.code == uferr.Code()
}

// 在调用处重新设置CallStack(skip == 3, deep = 20)
// 并在基础错误之后增加args...的错误信息
func (e *UFErrorStruct) New(args ...interface{}) UFError {
	extraMessage := fmt.Sprint(args...)
	ret := &UFErrorStruct{
		code:      e.code,
		message:   e.message,
		requestId: e.requestId,
		debugInfo: e.debugInfo,
	}
	ret.WithCallstack(3, 20)
	if len(extraMessage) > 0 {
		ret.message += ", " + extraMessage
	}
	return ret
}

func (e *UFErrorStruct) Newf(format string, args ...interface{}) UFError {
	extraMessage := fmt.Sprintf(format, args...)
	ret := &UFErrorStruct{
		code:      e.code,
		message:   e.message,
		requestId: e.requestId,
		debugInfo: e.debugInfo,
	}
	ret.WithCallstack(3, 20)
	if len(extraMessage) > 0 {
		ret.message += ", " + extraMessage
	}
	return ret
}

func NewFromError(err error) UFError {
	str := fmt.Sprintf("%v", err)
	kvs := strings.Split(str, ", ")
	ret := &UFErrorStruct{}

	ToInt := func(str string) int {
		var ret = 0
		i64, err := strconv.ParseInt(str, 10, 32)
		if nil == err {
			ret = int(i64)
		}
		return ret
	}
	for _, kvstr := range kvs {
		kv := strings.Split(kvstr, ": ")

		if len(kv) == 2 {
			switch kv[0] {
			case "code":
				ret.code = ToInt(kv[1])
			case "message":
				ret.message = kv[1]
			case "requestId":
				ret.requestId = kv[1]
			case "suggestHttpCode":
				ret.suggestHttpCode = ToInt(kv[1])
			case "debugInfo":
				ret.debugInfo = kv[1]
			case "callstack":
				json.Unmarshal([]byte(kv[1]), &ret.callstack)
			}
		}
	}
	return ret
}
