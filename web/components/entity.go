package components

type DiscordChannelSet struct {
	ID         string
	Name       string
	Ng         bool
	BotMessage bool
	NgTypes    []int
	NgUsers    []string
	NgRoles    []string
}

type DiscordChannel struct {
	ID       string
	Name     string
	Position int
}

type DiscordChannelSelect struct {
	ID   string
	Name string
}
