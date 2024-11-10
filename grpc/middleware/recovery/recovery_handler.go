package recovery

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HandlerContext recovery处理器
func HandlerContext(ctx context.Context, panicInfo any) (err error) {
	stackTrace := takeStacktrace(4)
	panicStr := fmt.Sprintf("%+v", panicInfo)
	err = status.Error(codes.Internal, panicStr+"\n"+stackTrace)
	return
}

// takeStacktrace 获取栈信息
func takeStacktrace(skip int) string {
	pcs := make([]uintptr, 64)
	var numFrames int
	for {
		// Skip the call to runtime.Callers and takeStacktrace so that the
		// program counters start at the caller of takeStacktrace.
		numFrames = runtime.Callers(skip+2, pcs)
		if numFrames < len(pcs) {
			break
		}
		// Don't put the too-short counter slice back into the pool; this lets
		// the pool adjust if we consistently take deep stacktraces.
		pcs = make([]uintptr, len(pcs)*2)
	}

	i := 0
	frames := runtime.CallersFrames(pcs[:numFrames])
	buffer := bytes.NewBuffer([]byte{})
	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if i != 0 {
			buffer.WriteByte('\n')
		}
		i++
		buffer.WriteString(frame.Function)
		buffer.WriteByte('\n')
		buffer.WriteByte('\t')
		buffer.WriteString(frame.File)
		buffer.WriteByte(':')
		buffer.WriteString(strconv.Itoa(frame.Line))
	}

	return buffer.String()
}
