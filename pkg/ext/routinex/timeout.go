package routinex

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"runtime/debug"
	"time"
)

func GoWithTimeout(fc func(), name string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorf("go func %s panic: %v\n%s", name, r, debug.Stack())
			}
			close(done)
		}()

		fc()
	}()

	select {
	case <-ctx.Done():
		logx.Errorf("go func %s timeout", name)
	case <-done:
		logx.Debugf("go func %s done", name)
	}
}

func GoWithTimeoutNonBlocking(fc func(), name string, timeout time.Duration) {
	go GoWithTimeout(fc, name, timeout)
}

func GoWithRetryAndTimeout(fc func(), name string, timeout time.Duration, retries int) {
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		done := make(chan error, 1) // 用error通道区分正常执行或panic

		go func() {
			defer func() {
				if r := recover(); r != nil {
					// 将panic信息发送到通道
					done <- fmt.Errorf("panic: %v\n%s", r, debug.Stack())
				}
			}()
			// 正常执行完成后发送nil，表示成功
			fc()
			done <- nil
		}()

		select {
		case <-ctx.Done():
			// 当前调用超时
			logx.Errorf("go func %s timeout", name)
		case err := <-done:
			// 根据err判断是正常完成还是panic
			if err == nil {
				// 正常完成，不需要重试，直接返回
				logx.Debugf("go func %s done", name)
				cancel()
				return
			} else {
				// 出现panic，需要进行重试
				logx.Errorf("go func %s encounter panic: %v", name, err)
			}
		}

		// 释放资源
		cancel()

		// 输出重试日志
		logx.Infof("retrying %s (%d/%d)", name, i+1, retries)

		// 指数退避机制：避免过于频繁的重试
		time.Sleep(time.Second << uint(i))
	}

	// 所有重试失败后，记录日志
	logx.Errorf("go func %s failed after %d retries", name, retries)
}

func GoWithRetryAndTimeoutNonBlocking(fc func(), name string, timeout time.Duration, retries int) {
	go GoWithRetryAndTimeout(fc, name, timeout, retries)
}
