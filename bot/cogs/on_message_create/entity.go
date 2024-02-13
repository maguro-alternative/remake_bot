package on_message_create

type LineChannel struct {
	Ng         bool   `db:"ng"`
	BotMessage bool   `db:"bot_message"`
}

type LineBot struct {
	LineNotifyToken   []byte `db:"line_notify_token"`
	LineBotToken      []byte `db:"line_bot_token"`
	LineBotSecret     []byte `db:"line_bot_secret"`
	LineGroupID       []byte `db:"line_group_id"`
	Iv                []byte `db:"iv"`
	DefaultChannelID  string `db:"default_channel_id"`
	DebugMode         bool   `db:"debug_mode"`
}

type LineBotDecrypt struct {
	LineNotifyToken   string
	LineBotToken      string
	LineGroupID       string
	DefaultChannelID  string
	DebugMode         bool
}
