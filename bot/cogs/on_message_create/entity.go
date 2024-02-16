package on_message_create

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

type LineBot struct {
	LineNotifyToken  []byte `db:"line_notify_token"`
	LineBotToken     []byte `db:"line_bot_token"`
	LineBotSecret    []byte `db:"line_bot_secret"`
	LineGroupID      []byte `db:"line_group_id"`
	DefaultChannelID string `db:"default_channel_id"`
	DebugMode        bool   `db:"debug_mode"`
}

type LineBotIv struct {
	LineNotifyTokenIv  []byte `db:"line_notify_token_iv"`
	LineBotTokenIv     []byte `db:"line_bot_token_iv"`
	LineBotSecretIv    []byte `db:"line_bot_secret_iv"`
	LineGroupIDIv      []byte `db:"line_group_id_iv"`
}

type LineBotDecrypt struct {
	LineNotifyToken  string
	LineBotToken     string
	LineGroupID      string
	DefaultChannelID string
	DebugMode        bool
}
