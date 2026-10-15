package watchdog

import (
	"context"
	"log/slog"
	"time"
)

// Run はcheckIntervalごとにisHealthyを確認し、
// unhealthyな状態がunhealthyLimit以上続いた場合にonUnhealthyを呼ぶ。
func Run(
	ctx context.Context,
	checkInterval time.Duration,
	unhealthyLimit time.Duration,
	isHealthy func() bool,
	onUnhealthy func(),
) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	var unhealthySince time.Time
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		if isHealthy() {
			if !unhealthySince.IsZero() {
				slog.InfoContext(ctx, "healthyに復旧しました",
					"unhealthy継続時間", time.Since(unhealthySince).String(),
				)
				unhealthySince = time.Time{}
			}
			continue
		}

		if unhealthySince.IsZero() {
			unhealthySince = time.Now()
			slog.WarnContext(ctx, "unhealthyになりました")
			continue
		}

		if elapsed := time.Since(unhealthySince); elapsed >= unhealthyLimit {
			slog.ErrorContext(ctx, "unhealthyのまま復旧しないため、プロセスを終了します",
				"unhealthy継続時間", elapsed.String(),
			)
			onUnhealthy()
			return
		}
	}
}
