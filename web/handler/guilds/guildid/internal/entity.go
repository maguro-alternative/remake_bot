package internal

import (
	"github.com/lib/pq"
)

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
