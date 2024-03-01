package on_message_create

import (
	"github.com/lib/pq"
)

type TestLineChannel struct {
	ChannelID  string `db:"channel_id"`
	GuildID    string `db:"guild_id"`
	Ng         bool   `db:"ng"`
	BotMessage bool   `db:"bot_message"`
}

type LineChannel struct {
	Ng         bool `db:"ng"`
	BotMessage bool `db:"bot_message"`
}

type LineNgID struct {
	ID    string `db:"id"`
	IDType string `db:"id_type"`
}

type LineBot struct {
	LineNotifyToken  pq.ByteaArray `db:"line_notify_token"`
	LineBotToken     pq.ByteaArray `db:"line_bot_token"`
	LineBotSecret    pq.ByteaArray `db:"line_bot_secret"`
	LineGroupID      pq.ByteaArray `db:"line_group_id"`
	DefaultChannelID string        `db:"default_channel_id"`
	DebugMode        bool          `db:"debug_mode"`
}

type LineBotIv struct {
	LineNotifyTokenIv pq.ByteaArray `db:"line_notify_token_iv"`
	LineBotTokenIv    pq.ByteaArray `db:"line_bot_token_iv"`
	LineBotSecretIv   pq.ByteaArray `db:"line_bot_secret_iv"`
	LineGroupIDIv     pq.ByteaArray `db:"line_group_id_iv"`
}

type LineBotDecrypt struct {
	LineNotifyToken  string
	LineBotToken     string
	LineGroupID      string
	DefaultChannelID string
	DebugMode        bool
}
