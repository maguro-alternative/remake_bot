package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type LineBotJson struct {
	GuildID          string `json:"guild_id" db:"guild_id"`
	LineNotifyToken  string `json:"line_notify_token,omitempty" db:"line_notify_token"`
	LineBotToken     string `json:"line_bot_token,omitempty" db:"line_bot_token"`
	LineBotSecret    string `json:"line_bot_secret,omitempty" db:"line_bot_secret"`
	LineGroupID      string `json:"line_group_id,omitempty" db:"line_group_id"`
	LineClientID     string `json:"line_client_id,omitempty" db:"line_client_id"`
	LineClientSecret string `json:"line_client_secret,omitempty" db:"line_client_secret"`
	DefaultChannelID string `json:"default_channel_id,omitempty" db:"default_channel_id"`
	DebugMode        bool   `json:"debug_mode,omitempty" db:"debug_mode"`
}

func (g LineBotJson) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.GuildID, validation.Required),
		validation.Field(&g.LineNotifyToken, is.Alphanumeric),
		validation.Field(&g.LineBotToken, is.Alphanumeric),
		validation.Field(&g.LineBotSecret, is.Alphanumeric),
		validation.Field(&g.LineGroupID, is.Alphanumeric),
		validation.Field(&g.LineClientID, is.Alphanumeric),
		validation.Field(&g.LineClientSecret, is.Alphanumeric),
		validation.Field(&g.DefaultChannelID, is.Digit),
	)
}

type LineBot struct {
	GuildID          string `db:"guild_id"`
	LineNotifyToken  []byte `db:"line_notify_token"`
	LineBotToken     []byte `db:"line_bot_token"`
	LineBotSecret    []byte `db:"line_bot_secret"`
	LineGroupID      []byte `db:"line_group_id"`
	LineClientID     []byte `db:"line_client_id"`
	LineClientSecret []byte `db:"line_client_secret"`
	DefaultChannelID string `db:"default_channel_id"`
	DebugMode        bool   `db:"debug_mode"`
}

type LineBotIv struct {
	GuildID            string `db:"guild_id"`
	LineNotifyTokenIv  []byte `db:"line_notify_token_iv"`
	LineBotTokenIv     []byte `db:"line_bot_token_iv"`
	LineBotSecretIv    []byte `db:"line_bot_secret_iv"`
	LineGroupIDIv      []byte `db:"line_group_id_iv"`
	LineClientIDIv     []byte `db:"line_client_id_iv"`
	LineClientSecretIv []byte `db:"line_client_secret_iv"`
}
