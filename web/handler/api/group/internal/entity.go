package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type LineBotJson struct {
	GuildID          string `json:"guildId,omitempty" db:"guild_id"`
	DefaultChannelID string `json:"defaultChannelId,omitempty" db:"default_channel_id"`
	DebugMode        bool   `json:"debugMode,omitempty" db:"debug_mode"`
}

func (g LineBotJson) Validate() error {
	return validation.ValidateStruct(&g,
		//validation.Field(&g.GuildID, is.Digit),
		validation.Field(&g.DefaultChannelID, is.Digit),
		validation.Field(&g.DefaultChannelID, validation.Required),
	)
}
