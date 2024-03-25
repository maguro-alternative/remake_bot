package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type LineBotJson struct {
	DefaultChannelID string `json:"default_channel_id,omitempty" db:"default_channel_id"`
	DebugMode        bool   `json:"debug_mode,omitempty" db:"debug_mode"`
}

func (g LineBotJson) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.DefaultChannelID, is.Digit),
		validation.Field(&g.DefaultChannelID, validation.Required),
	)
}

type LineBot struct {
	GuildID          string        `db:"guild_id"`
	DefaultChannelID string        `db:"default_channel_id"`
}