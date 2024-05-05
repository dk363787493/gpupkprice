package common

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

func StackInfo(err error) string {
	if err == nil {
		return ""
	}

	_, file, line, _ := runtime.Caller(1)

	return fmt.Sprintf("%s:%d\n%s", file, line, string(debug.Stack()))
}
