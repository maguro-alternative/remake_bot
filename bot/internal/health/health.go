package health

import (
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
)

// disconnectCount はWebSocket切断の連続回数を追跡する。
// Resumed/Readyで0にリセットされる。
var disconnectCount atomic.Int64

// lastDisconnectTime は最後にDisconnectイベントを受信した時刻。
var lastDisconnectTime atomic.Int64

// MaxConsecutiveDisconnects を超えて連続切断が発生した場合、
// ヘルスチェックが失敗を返すようになる。
const MaxConsecutiveDisconnects = 5

// OnDisconnect はWebSocket切断時に呼ばれる。
func OnDisconnect(shouldReconnect bool) {
	count := disconnectCount.Add(1)
	lastDisconnectTime.Store(time.Now().Unix())
	slog.Warn("Discord WebSocket切断を検知しました",
		"連続切断回数", count,
		"自動再接続待ち", shouldReconnect,
	)
}

// OnResumed はWebSocket再接続(Resume)成功時に呼ばれる。
func OnResumed() {
	prev := disconnectCount.Swap(0)
	slog.Info("Discord WebSocket再接続(Resume)に成功しました",
		"復旧までの切断回数", prev,
	)
}

// OnReady はWebSocket再接続(Ready)成功時に呼ばれる。
func OnReady() {
	prev := disconnectCount.Swap(0)
	if prev > 0 {
		slog.Info("Discord WebSocket再接続(Ready)に成功しました",
			"復旧までの切断回数", prev,
		)
	}
}

// IsHealthy はBotのWebSocket接続が健全かどうかを返す。
// - DataReadyがfalseの場合はfalse
// - 連続切断回数がMaxConsecutiveDisconnectsを超えている場合はfalse
func IsHealthy(s *discordgo.Session) bool {
	if !s.DataReady {
		return false
	}
	if disconnectCount.Load() >= MaxConsecutiveDisconnects {
		return false
	}
	return true
}

// GetDisconnectInfo はヘルスチェック用の診断情報を返す。
func GetDisconnectInfo() (consecutiveDisconnects int64, lastDisconnect time.Time) {
	count := disconnectCount.Load()
	ts := lastDisconnectTime.Load()
	var t time.Time
	if ts > 0 {
		t = time.Unix(ts, 0)
	}
	return count, t
}
