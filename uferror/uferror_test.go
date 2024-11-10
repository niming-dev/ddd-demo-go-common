package uferror

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/go-playground/assert/v2"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func generalCheck(t *testing.T, ue UFError, code int, message, requestId string, httpCode int, debugInfo string) {
	assert.Equal(t, ue.Code(), code)
	assert.Equal(t, ue.Message(), message)
	assert.Equal(t, ue.RequestId(), requestId)
	assert.Equal(t, ue.SuggestHttpCode(), httpCode)
	assert.Equal(t, ue.DebugInfo(), debugInfo)
}

func TestSimple(t *testing.T) {
	ue1 := New(10, "test")
	generalCheck(t, ue1, 10, "test", "", 0, "")

	e1 := ue1.Error()
	stre1 := fmt.Sprintf("%v", e1)
	assert.Equal(t, stre1, "code: 10, message: test")

	e2 := ue1.ToDebugError()
	stre2 := fmt.Sprintf("%v", e2)
	assert.Equal(t, stre2, "code: 10, message: test")
}

func TestSimpleWithRequestId(t *testing.T) {
	ue1 := New(10, "test").WithRequestId("request1234")
	generalCheck(t, ue1, 10, "test", "request1234", 0, "")

	e1 := ue1.Error()
	stre1 := fmt.Sprintf("%v", e1)
	assert.Equal(t, stre1, "code: 10, message: test, requestId: request1234")
}

func TestWithHttpCode(t *testing.T) {
	ue1 := New(10, "test").WithSuggestHttpCode(404).WithRequestId("request123")
	generalCheck(t, ue1, 10, "test", "request123", 404, "")

	e1 := ue1.Error()
	stre1 := fmt.Sprintf("%v", e1)
	assert.Equal(t, stre1, "code: 10, message: test, requestId: request123")

	e2 := ue1.ToDebugError()
	stre2 := fmt.Sprintf("%v", e2)
	assert.Equal(t, stre2, "code: 10, message: test, requestId: request123, suggestHttpCode: 404")
}

func TestDebugInfo(t *testing.T) {
	ue1 := New(10, "test").WithDebugInfof("mysql execute %s error", "select").WithRequestId("request123")
	generalCheck(t, ue1, 10, "test", "request123", 0, "mysql execute select error")

	e1 := ue1.Error()
	stre1 := fmt.Sprintf("%v", e1)
	assert.Equal(t, stre1, "code: 10, message: test, requestId: request123")

	e2 := ue1.ToDebugError()
	stre2 := fmt.Sprintf("%v", e2)
	expectHeadStr := "code: 10, message: test, requestId: request123, debugInfo: mysql execute select error"
	assert.MatchRegex(t, stre2, regexp.MustCompile(expectHeadStr))
	callstackString := "["
	for _, frame := range ue1.Callstack() {
		callstackString += `"` + frame + `" `
	}
	callstackString = callstackString[:len(callstackString)-1] + "]"
	if len(callstackString) == 0 {
		expectHeadStr += "callstack: " + callstackString
		assert.Equal(t, stre2, expectHeadStr)
	}
}

func TestCallstack(t *testing.T) {
	ue1 := New(10, "test").WithDebugInfof("mysql execute %s error", "select").WithRequestId("request123").
		WithCallstack(1, 3)
	generalCheck(t, ue1, 10, "test", "request123", 0, "mysql execute select error")

	e1 := ue1.Error()
	stre1 := fmt.Sprintf("%v", e1)
	assert.Equal(t, stre1, "code: 10, message: test, requestId: request123")

	e2 := ue1.ToDebugError()
	stre2 := fmt.Sprintf("%v", e2)
	expectHeadStr := "code: 10, message: test, requestId: request123, debugInfo: mysql execute select error"
	assert.MatchRegex(t, stre2, regexp.MustCompile(expectHeadStr))
	callstackString := "["
	for _, frame := range ue1.Callstack() {
		callstackString += `"` + frame + `" `
	}
	callstackString = callstackString[:len(callstackString)-1] + "]"
	if len(callstackString) == 0 {
		expectHeadStr += "callstack: " + callstackString
		assert.Equal(t, stre2, expectHeadStr)
	}
}

func Test_Convert(t *testing.T) {
	A := func() UFError {
		return New(-10, "abccd").WithDebugInfo("mmmm")
	}

	a := A()
	log.Printf("%#v", a)

	b := a.New("hello")
	log.Printf("%#v", b)

	assert.Equal(t, a.IsKindOf(b), true)
}
