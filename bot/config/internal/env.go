package internal

type Config struct {
	PrivateKey               string `env:"PRIVATE_KEY" envDefault:"645E739A7F9F162725C1533DC2C5E827"`
	DiscordBotToken          string `env:"DISCORD_BOT_TOKEN" envDefault:""`
	DBName                   string `env:"DB_NAME" envDefault:"postgres"`
	DBUser                   string `env:"DB_USER" envDefault:"postgres"`
	DBPassword               string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBHost                   string `env:"DB_HOST" envDefault:"localhost"`
	DBPort                   string `env:"DB_PORT" envDefault:"5432"`
	SessionSecret            string `env:"SESSION_SECRET" envDefault:""`
	VoiceVoxKey              string `env:"VOICEVOX_KEY" envDefault:""`
	InternalURL              string `env:"INTERNAL_URL" envDefault:"http://localhost:8080"`
	ChannelNo                string `env:"CHANNEL_NO" envDefault:""`
	SlachCommandDebugGuildID string `env:"SLACK_COMMAND_DEBUG_GUILD_ID" envDefault:""`
	LineWorksID       string `env:"LINE_WORKS_ID" envDefault:""`
	LineWorksPassword string `env:"LINE_WORKS_PASSWORD" envDefault:""`
}
