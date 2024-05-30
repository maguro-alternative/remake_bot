package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type VcSignalJson struct {
	GuildID   string      `json:"guildId"`
	VcSignals []vcSignals `json:"vcSignals"`
}

type vcSignals struct {
	VcChannelID            string   `json:"vcChannelId"`
	SendSignal             bool     `json:"sendSignal"`
	SendChannelId          string   `json:"sendChannelId"`
	JoinBot                bool     `json:"joinBot"`
	EveryoneMention        bool     `json:"everyoneMention"`
	VcSignalNgUserIDs      []string `json:"vcSignalNgUserIds"`
	VcSignalNgRoleIDs      []string `json:"vcSignalNgRoleIds"`
	VcSignalMentionUserIDs []string `json:"vcSignalMentionUserIds"`
	VcSignalMentionRoleIDs []string `json:"vcSignalMentionRoleIds"`
}

func (g VcSignalJson) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.GuildID, validation.Required),
	)
}
