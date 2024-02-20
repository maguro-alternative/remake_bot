package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
	//"github.com/go-ozzo/ozzo-validation/is"
)

type LineChannelJson struct {
	GuildID  string `json:"guild_id"`
	Channels []struct {
		ChannelID  string   `json:"channel_id"`
		NG         bool     `json:"ng"`
		BotMessage bool     `json:"bot_message"`
		NGTypes    []string `json:"ng_types"`
		NGUsers    []string `json:"ng_users"`
		NGRoles    []string `json:"ng_roles"`
	} `json:"channels"`
}

func (g LineChannelJson) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.GuildID, validation.Required),
		validation.Field(&g.Channels, validation.Required),
	)
}

type LineChannel struct {
	ChannelID  string `db:"channel_id"`
	GuildID    string `db:"guild_id"`
	NG         bool   `db:"ng"`
	BotMessage bool   `db:"bot_message"`
}

type LineNgType struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	Type      string `db:"type"`
}

type LineNgID struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	ID        string `db:"id"`
	IDType    string `db:"id_type"`
}
