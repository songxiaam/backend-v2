package routinex

import "github.com/zeromicro/go-zero/core/logx"

// Recovery 捕获panic
func Recovery() {
	if err := recover(); err != nil {
		logx.Errorf("panic: %v", err)
	}
}
