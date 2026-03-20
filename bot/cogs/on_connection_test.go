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

func TestIsHealthy_DataReadyTrue_NoDisconnects(t *testing.T) {
	resetState()
	s := &discordgo.Session{DataReady: true}

	assert.True(t, IsHealthy(s))
}

func TestIsHealthy_DataReadyFalse(t *testing.T) {
	resetState()
	s := &discordgo.Session{DataReady: false}

	assert.False(t, IsHealthy(s))
}

func TestIsHealthy_ExceedsMaxConsecutiveDisconnects(t *testing.T) {
	resetState()
	disconnectCount.Store(MaxConsecutiveDisconnects)
	s := &discordgo.Session{DataReady: true}

	assert.False(t, IsHealthy(s))
}

func TestIsHealthy_BelowMaxConsecutiveDisconnects(t *testing.T) {
	resetState()
	disconnectCount.Store(MaxConsecutiveDisconnects - 1)
	s := &discordgo.Session{DataReady: true}

	assert.True(t, IsHealthy(s))
}

func TestIsHealthy_DataReadyFalseAndExceedsMax(t *testing.T) {
	resetState()
	disconnectCount.Store(MaxConsecutiveDisconnects + 1)
	s := &discordgo.Session{DataReady: false}

	assert.False(t, IsHealthy(s))
}

func TestGetDisconnectInfo_NoDisconnects(t *testing.T) {
	resetState()

	count, lastDisconnect := GetDisconnectInfo()

	assert.Equal(t, int64(0), count)
	assert.True(t, lastDisconnect.IsZero())
}

func TestGetDisconnectInfo_WithDisconnects(t *testing.T) {
	resetState()
	now := time.Now()
	disconnectCount.Store(3)
	lastDisconnectTime.Store(now.Unix())

	count, lastDisconnect := GetDisconnectInfo()

	assert.Equal(t, int64(3), count)
	assert.Equal(t, now.Unix(), lastDisconnect.Unix())
}

func TestOnDisconnect_IncrementsCount(t *testing.T) {
	resetState()

	onDisconnect(true)

	assert.Equal(t, int64(1), disconnectCount.Load())
	assert.NotEqual(t, int64(0), lastDisconnectTime.Load())
}

func TestOnDisconnect_MultipleCalls(t *testing.T) {
	resetState()

	onDisconnect(true)
	onDisconnect(true)
	onDisconnect(true)

	assert.Equal(t, int64(3), disconnectCount.Load())
}

func TestOnDisconnect_SetsLastDisconnectTime(t *testing.T) {
	resetState()
	before := time.Now().Unix()

	onDisconnect(true)

	after := time.Now().Unix()
	ts := lastDisconnectTime.Load()
	assert.GreaterOrEqual(t, ts, before)
	assert.LessOrEqual(t, ts, after)
}

func TestOnResumed_ResetsCount(t *testing.T) {
	resetState()
	disconnectCount.Store(3)

	onResumed()

	assert.Equal(t, int64(0), disconnectCount.Load())
}

func TestOnReady_ResetsCount(t *testing.T) {
	resetState()
	disconnectCount.Store(2)

	onReady()

	assert.Equal(t, int64(0), disconnectCount.Load())
}

func TestOnReady_NoOpWhenCountIsZero(t *testing.T) {
	resetState()

	onReady()

	assert.Equal(t, int64(0), disconnectCount.Load())
}

func TestDisconnectThenResumed_ResetsHealthy(t *testing.T) {
	resetState()
	s := &discordgo.Session{DataReady: true}

	// 5回切断 → unhealthy
	for i := 0; i < int(MaxConsecutiveDisconnects); i++ {
		onDisconnect(true)
	}
	assert.False(t, IsHealthy(s))

	// Resume成功 → healthy
	onResumed()
	assert.True(t, IsHealthy(s))
}

func TestDisconnectThenReady_ResetsHealthy(t *testing.T) {
	resetState()
	s := &discordgo.Session{DataReady: true}

	// 5回切断 → unhealthy
	for i := 0; i < int(MaxConsecutiveDisconnects); i++ {
		onDisconnect(true)
	}
	assert.False(t, IsHealthy(s))

	// Ready成功 → healthy
	onReady()
	assert.True(t, IsHealthy(s))
}

func TestGetDisconnectInfo_AfterDisconnectAndResume(t *testing.T) {
	resetState()

	onDisconnect(true)
	onDisconnect(true)
	_, lastDisconnect := GetDisconnectInfo()
	assert.False(t, lastDisconnect.IsZero())

	onResumed()
	count, lastDisconnectAfterResume := GetDisconnectInfo()
	assert.Equal(t, int64(0), count)
	// lastDisconnectTimeはResumeでリセットされない（最後の切断時刻として保持）
	assert.False(t, lastDisconnectAfterResume.IsZero())
}
