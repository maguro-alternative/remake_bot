package internal

type Config struct {
	DiscordBotToken string `env:"DISCORD_BOT_TOKEN" envDefault:""`
}