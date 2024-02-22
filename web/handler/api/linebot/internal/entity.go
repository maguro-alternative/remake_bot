package internal

import "github.com/lib/pq"

type LineBot struct {
	GuildID          string        `db:"guild_id"`
	LineNotifyToken  pq.ByteaArray `db:"line_notify_token"`
	LineBotToken     pq.ByteaArray `db:"line_bot_token"`
	LineBotSecret    pq.ByteaArray `db:"line_bot_secret"`
	LineGroupID      pq.ByteaArray `db:"line_group_id"`
	LineClientID     pq.ByteaArray `db:"line_client_id"`
	LineClientSecret pq.ByteaArray `db:"line_client_secret"`
	DefaultChannelID string        `db:"default_channel_id"`
	DebugMode        bool          `db:"debug_mode"`
}

type LineBotIv struct {
	LineNotifyTokenIv   pq.ByteaArray `db:"line_notify_token_iv"`
	LineBotTokenIv      pq.ByteaArray `db:"line_bot_token_iv"`
	LineBotSecretIv     pq.ByteaArray `db:"line_bot_secret_iv"`
	LineGroupIDIv       pq.ByteaArray `db:"line_group_id_iv"`
	LineClientIDIv      pq.ByteaArray `db:"line_client_id_iv"`
	LineClientSercretIv pq.ByteaArray `db:"line_client_sercret_iv"`
}

type LineBotDecrypt struct {
	LineNotifyToken   string
	LineBotToken      string
	LineGroupID       string
	LineClientID      string
	LineClientSercret string
	DefaultChannelID  string
	DebugMode         bool
}
