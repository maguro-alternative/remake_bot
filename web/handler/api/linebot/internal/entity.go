package internal

type LineBot struct {
	LineNotifyToken   []byte `db:"line_notify_token"`
	LineBotToken      []byte `db:"line_bot_token"`
	LineBotSecret     []byte `db:"line_bot_secret"`
	LineGroupID       []byte `db:"line_group_id"`
	LineClientID      []byte `db:"line_client_id"`
	LineClientSercret []byte `db:"line_client_sercret"`
	Iv                []byte `db:"iv"`
	DefaultChannelID  string `db:"default_channel_id"`
	DebugMode         bool   `db:"debug_mode"`
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
