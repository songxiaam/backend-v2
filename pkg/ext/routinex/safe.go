package routinex

import (
	"context"
	"fmt"
	"runtime/debug"
)

func SafeBlockGo(ctx context.Context, f func() error) error {
	var err error
	done := make(chan struct{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic recovered: %v\n%s", r, debug.Stack())
			}
			close(done)
		}()
		err = f()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return err
	}
}

// SafeGo 非阻塞的安全协程执行方法
func SafeGo(ctx context.Context, f func() error) chan error {
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("panic recovered: %v\n%s", r, debug.Stack())
			}
			close(errChan)
		}()

		if err := f(); err != nil {
			errChan <- err
		}
	}()

	return errChan
}

// SafeFunc 封装一个适配err group的安全执行方法
func SafeFunc(f func() error) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic recovered: %v\n%s", r, debug.Stack())
			}
		}()
		return f()
	}
}
