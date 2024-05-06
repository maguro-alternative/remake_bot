package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/lib/pq"
)

type LineBotJson struct {
	GuildID          string `json:"guildId" db:"guild_id"`
	LineNotifyToken  string `json:"lineNotifyToken,omitempty" db:"line_notify_token"`
	LineBotToken     string `json:"lineBotToken,omitempty" db:"line_bot_token"`
	LineBotSecret    string `json:"lineBotSecret,omitempty" db:"line_bot_secret"`
	LineGroupID      string `json:"lineGroupId,omitempty" db:"line_group_id"`
	LineClientID     string `json:"lineClientId,omitempty" db:"line_client_id"`
	LineClientSecret string `json:"lineClientSecret,omitempty" db:"line_client_secret"`
	DefaultChannelID string `json:"defaultChannelId,omitempty" db:"default_channel_id"`
	DebugMode        bool   `json:"debugMode,omitempty" db:"debug_mode"`
}

func (g LineBotJson) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.LineNotifyToken, is.Alphanumeric),
		validation.Field(&g.LineBotSecret, is.Alphanumeric),
		validation.Field(&g.LineGroupID, is.Alphanumeric),
		validation.Field(&g.LineClientID, is.Alphanumeric),
		validation.Field(&g.LineClientSecret, is.Alphanumeric),
		validation.Field(&g.DefaultChannelID, is.Digit),
	)
}

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
	GuildID            string        `db:"guild_id"`
	LineNotifyTokenIv  pq.ByteaArray `db:"line_notify_token_iv"`
	LineBotTokenIv     pq.ByteaArray `db:"line_bot_token_iv"`
	LineBotSecretIv    pq.ByteaArray `db:"line_bot_secret_iv"`
	LineGroupIDIv      pq.ByteaArray `db:"line_group_id_iv"`
	LineClientIDIv     pq.ByteaArray `db:"line_client_id_iv"`
	LineClientSecretIv pq.ByteaArray `db:"line_client_secret_iv"`
}
