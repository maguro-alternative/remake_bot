package internal

type DiscordChannel struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Position int    `db:"position"`
}

type DiscordChannelSet struct {
	ID         string
	Name       string
	Ng         bool
	BotMessage bool
	NgTypes    []int
	NgUsers    []string
	NgRoles    []string
}

type LineChannel struct {
	ChannelID  string `db:"channel_id"`
	GuildID    string `db:"guild_id"`
	Ng         bool   `db:"ng"`
	BotMessage bool   `db:"bot_message"`
}

type LineNgType struct {
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
