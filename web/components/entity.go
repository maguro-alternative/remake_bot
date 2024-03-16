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
