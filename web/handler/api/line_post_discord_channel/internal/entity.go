package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
	//"github.com/go-ozzo/ozzo-validation/is"
)

type LinePostDiscordChannelJson struct {
	GuildID  string `json:"guildId"`
	Channels []struct {
		ChannelID  string   `json:"channelId"`
		Ng         bool     `json:"ng"`
		BotMessage bool     `json:"botMessage"`
		NgTypes    []int    `json:"ngTypes"`
		NgUsers    []string `json:"ngUsers"`
		NgRoles    []string `json:"ngRoles"`
	} `json:"channels"`
}

func (g LinePostDiscordChannelJson) Validate() error {
	return validation.ValidateStruct(&g,
		//validation.Field(&g.GuildID, validation.Required),
		validation.Field(&g.Channels, validation.Required),
	)
}

type LinePostDiscordChannel struct {
	ChannelID  string `db:"channelId"`
	GuildID    string `db:"guildId"`
	Ng         bool   `db:"ng"`
	BotMessage bool   `db:"botMessage"`
}

type LineNgDiscordMessageType struct {
	ChannelID string `db:"channelId"`
	GuildID   string `db:"guildId"`
	Type      int    `db:"type"`
}
