package internal

type Config struct {
	PrivateKey      string `env:"PRIVATE_KEY" envDefault:""`
	DiscordBotToken string `env:"DISCORD_BOT_TOKEN" envDefault:""`
}
