package watchdog

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("unhealthyがunhealthyLimit以上続いたらonUnhealthyが呼ばれる", func(t *testing.T) {
		called := make(chan struct{})

		go Run(context.Background(), 5*time.Millisecond, 20*time.Millisecond,
			func() bool { return false },
			func() { close(called) },
		)

		select {
		case <-called:
		case <-time.After(time.Second):
			t.Fatal("onUnhealthyが呼ばれなかった")
		}
	})

	t.Run("healthyのままならonUnhealthyは呼ばれない", func(t *testing.T) {
		var called atomic.Bool
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		Run(ctx, 5*time.Millisecond, 20*time.Millisecond,
			func() bool { return true },
			func() { called.Store(true) },
		)

		assert.False(t, called.Load())
	})

	t.Run("unhealthyLimit未満で復旧すればonUnhealthyは呼ばれない", func(t *testing.T) {
		var called atomic.Bool
		var healthy atomic.Bool
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		// 30msだけunhealthyにしてから復旧させる(unhealthyLimitは100ms)
		go func() {
			time.Sleep(30 * time.Millisecond)
			healthy.Store(true)
		}()

		Run(ctx, 5*time.Millisecond, 100*time.Millisecond,
			healthy.Load,
			func() { called.Store(true) },
		)

		assert.False(t, called.Load())
	})

	t.Run("contextがキャンセルされたら終了する", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})

		go func() {
			Run(ctx, 5*time.Millisecond, time.Hour,
				func() bool { return true },
				func() {},
			)
			close(done)
		}()
		cancel()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("Runが終了しなかった")
		}
	})
}
