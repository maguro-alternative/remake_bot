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

type LineBotByteEntered struct {
	LineNotifyToken  [][]byte
	LineBotToken     [][]byte
	LineBotSecret    [][]byte
	LineGroupID      [][]byte
	LineClientID     [][]byte
	LineClientSecret [][]byte
	LineDebugMode    bool
}

type LineEntered struct {
	LineNotifyToken  string
	LineBotToken     string
	LineBotSecret    string
	LineGroupID      string
	LineClientID     string
	LineClientSecret string
	LineDebugMode    bool
}

type PermissionCode struct {
	GuildID string
	Type    string
	Code    int64
}

type PermissionID struct {
	GuildID    string
	Type       string
	TargetType string
	TargetID   string
	Permission string
}
