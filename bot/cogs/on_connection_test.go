package cogs

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

// resetState はテスト間でグローバル状態をリセットする。
func resetState() {
	disconnectCount.Store(0)
	lastDisconnectTime.Store(0)
}

func TestIsHealthy(t *testing.T) {
	t.Run("DataReady=trueかつ切断なしならhealthy", func(t *testing.T) {
		resetState()
		s := &discordgo.Session{DataReady: true}

		assert.True(t, IsHealthy(s))
	})

	t.Run("DataReady=falseならunhealthy", func(t *testing.T) {
		resetState()
		s := &discordgo.Session{DataReady: false}

		assert.False(t, IsHealthy(s))
	})

	t.Run("連続切断がMaxConsecutiveDisconnectsに達したらunhealthy", func(t *testing.T) {
		resetState()
		disconnectCount.Store(MaxConsecutiveDisconnects)
		s := &discordgo.Session{DataReady: true}

		assert.False(t, IsHealthy(s))
	})

	t.Run("連続切断がMaxConsecutiveDisconnects未満ならhealthy", func(t *testing.T) {
		resetState()
		disconnectCount.Store(MaxConsecutiveDisconnects - 1)
		s := &discordgo.Session{DataReady: true}

		assert.True(t, IsHealthy(s))
	})

	t.Run("DataReady=falseかつ連続切断超過の場合もunhealthy", func(t *testing.T) {
		resetState()
		disconnectCount.Store(MaxConsecutiveDisconnects + 1)
		s := &discordgo.Session{DataReady: false}

		assert.False(t, IsHealthy(s))
	})
}

func TestGetDisconnectInfo(t *testing.T) {
	t.Run("初期状態ではカウント0かつ時刻はゼロ値", func(t *testing.T) {
		resetState()

		count, lastDisconnect := GetDisconnectInfo()

		assert.Equal(t, int64(0), count)
		assert.True(t, lastDisconnect.IsZero())
	})

	t.Run("切断情報が設定されていれば正しく返す", func(t *testing.T) {
		resetState()
		now := time.Now()
		disconnectCount.Store(3)
		lastDisconnectTime.Store(now.Unix())

		count, lastDisconnect := GetDisconnectInfo()

		assert.Equal(t, int64(3), count)
		assert.Equal(t, now.Unix(), lastDisconnect.Unix())
	})

	t.Run("Resume後もlastDisconnectTimeは保持される", func(t *testing.T) {
		resetState()

		onDisconnect(true)
		onDisconnect(true)
		_, lastDisconnect := GetDisconnectInfo()
		assert.False(t, lastDisconnect.IsZero())

		onResumed()
		count, lastDisconnectAfterResume := GetDisconnectInfo()
		assert.Equal(t, int64(0), count)
		assert.False(t, lastDisconnectAfterResume.IsZero())
	})
}

func TestOnDisconnect(t *testing.T) {
	t.Run("呼び出すとカウントが1増える", func(t *testing.T) {
		resetState()

		onDisconnect(true)

		assert.Equal(t, int64(1), disconnectCount.Load())
		assert.NotEqual(t, int64(0), lastDisconnectTime.Load())
	})

	t.Run("複数回呼び出すとカウントが累積する", func(t *testing.T) {
		resetState()

		onDisconnect(true)
		onDisconnect(true)
		onDisconnect(true)

		assert.Equal(t, int64(3), disconnectCount.Load())
	})

	t.Run("lastDisconnectTimeに現在時刻が記録される", func(t *testing.T) {
		resetState()
		before := time.Now().Unix()

		onDisconnect(true)

		after := time.Now().Unix()
		ts := lastDisconnectTime.Load()
		assert.GreaterOrEqual(t, ts, before)
		assert.LessOrEqual(t, ts, after)
	})
}

func TestOnResumed(t *testing.T) {
	t.Run("カウントが0にリセットされる", func(t *testing.T) {
		resetState()
		disconnectCount.Store(3)

		onResumed()

		assert.Equal(t, int64(0), disconnectCount.Load())
	})
}

func TestOnReady(t *testing.T) {
	t.Run("カウントが0にリセットされる", func(t *testing.T) {
		resetState()
		disconnectCount.Store(2)

		onReady()

		assert.Equal(t, int64(0), disconnectCount.Load())
	})

	t.Run("カウントが0の場合もそのまま0", func(t *testing.T) {
		resetState()

		onReady()

		assert.Equal(t, int64(0), disconnectCount.Load())
	})
}

func TestDisconnectからの復旧(t *testing.T) {
	t.Run("切断→unhealthy→Resume→healthyに復旧", func(t *testing.T) {
		resetState()
		s := &discordgo.Session{DataReady: true}

		for i := 0; i < int(MaxConsecutiveDisconnects); i++ {
			onDisconnect(true)
		}
		assert.False(t, IsHealthy(s))

		onResumed()
		assert.True(t, IsHealthy(s))
	})

	t.Run("切断→unhealthy→Ready→healthyに復旧", func(t *testing.T) {
		resetState()
		s := &discordgo.Session{DataReady: true}

		for i := 0; i < int(MaxConsecutiveDisconnects); i++ {
			onDisconnect(true)
		}
		assert.False(t, IsHealthy(s))

		onReady()
		assert.True(t, IsHealthy(s))
	})
}
