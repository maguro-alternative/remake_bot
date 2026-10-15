package bot

import (
	"time"

	"github.com/maguro-alternative/remake_bot/bot/internal/health"

	"github.com/bwmarrin/discordgo"
)

// IsHealthy はBotのWebSocket接続が健全かどうかを返す。
func IsHealthy(s *discordgo.Session) bool {
	return health.IsHealthy(s)
}

// GetDisconnectInfo はヘルスチェック用の診断情報を返す。
func GetDisconnectInfo() (consecutiveDisconnects int64, lastDisconnect time.Time) {
	return health.GetDisconnectInfo()
}
