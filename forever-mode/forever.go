package forever_mode

import (
	"context"
	"time"
)

type runFn func() error

// RunForever 永远执行 fn
func RunForever(interval time.Duration, fn runFn) {

	timer := time.NewTimer(interval)
	defer func() {
		timer.Stop()
	}()

	// 立即调用
	if err := fn(); err != nil {
		panic(err)
	}

	for {
		select {
		case <-timer.C:
			if err := fn(); err != nil {
				panic(err)
			}
			timer.Reset(interval)
		}
	}
}

// RunWithChannel 执行 fn，当外部关闭 stopC chan 会退出
func RunWithChannel(interval time.Duration, fn runFn, stopC chan struct{}) {

	timer := time.NewTimer(interval)
	defer func() {
		timer.Stop()
	}()

	if err := fn(); err != nil {
		panic(err)
	}

	for {
		select {
		case <-timer.C:
			if err := fn(); err != nil {
				panic(err)
			}
			timer.Reset(interval)
		case <-stopC:
			return
		}
	}
}

// RunWithContext 调用 fn，外部传入的 ctx 能够自由关闭
func RunWithContext(interval time.Duration, fn runFn, ctx context.Context) {

	timer := time.NewTimer(interval)
	defer func() {
		timer.Stop()
	}()

	if err := fn(); err != nil {
		panic(err)
	}

	for {
		select {
		case <-timer.C:
			if err := fn(); err != nil {
				panic(err)
			}
			timer.Reset(interval)
		case <-ctx.Done():
			return
		}
	}
}
