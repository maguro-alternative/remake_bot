package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
	//"github.com/go-ozzo/ozzo-validation/is"
)

type LinePostDiscordChannelJson struct {
	GuildID  string `json:"guild_id"`
	Channels []struct {
		ChannelID  string   `json:"channel_id"`
		Ng         bool     `json:"ng"`
		BotMessage bool     `json:"bot_message"`
		NgTypes    []int    `json:"ng_types"`
		NgUsers    []string `json:"ng_users"`
		NgRoles    []string `json:"ng_roles"`
	} `json:"channels"`
}

func (g LinePostDiscordChannelJson) Validate() error {
	return validation.ValidateStruct(&g,
		//validation.Field(&g.GuildID, validation.Required),
		validation.Field(&g.Channels, validation.Required),
	)
}

type LinePostDiscordChannel struct {
	ChannelID  string `db:"channel_id"`
	GuildID    string `db:"guild_id"`
	Ng         bool   `db:"ng"`
	BotMessage bool   `db:"bot_message"`
}

type LineNgDiscordMessageType struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	Type      int    `db:"type"`
}

type LineNgID struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	ID        string `db:"id"`
	IDType    string `db:"id_type"`
}
