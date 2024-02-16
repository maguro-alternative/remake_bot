package internal

type LineBot struct {
	GuildID           string `db:"guild_id"`
	LineNotifyToken   []byte `db:"line_notify_token"`
	LineBotToken      []byte `db:"line_bot_token"`
	LineBotSecret     []byte `db:"line_bot_secret"`
	LineGroupID       []byte `db:"line_group_id"`
	LineClientID      []byte `db:"line_client_id"`
	LineClientSercret []byte `db:"line_client_sercret"`
	DefaultChannelID  string `db:"default_channel_id"`
	DebugMode         bool   `db:"debug_mode"`
}

type LineBotIv struct {
	LineNotifyTokenIv   []byte `db:"line_notify_token_iv"`
	LineBotTokenIv      []byte `db:"line_bot_token_iv"`
	LineBotSecretIv     []byte `db:"line_bot_secret_iv"`
	LineGroupIDIv       []byte `db:"line_group_id_iv"`
	LineClientIDIv      []byte `db:"line_client_id_iv"`
	LineClientSercretIv []byte `db:"line_client_sercret_iv"`
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
