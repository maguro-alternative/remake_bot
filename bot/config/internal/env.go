package internal

type Config struct {
	PrivateKey      string `env:"PRIVATE_KEY" envDefault:""`
	DiscordBotToken string `env:"DISCORD_BOT_TOKEN" envDefault:""`
	DBName          string `env:"DB_NAME" envDefault:"postgres"`
	DBUser			string `env:"DB_USER" envDefault:"postgres"`
	DBPassword		string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBHost			string `env:"DB_HOST" envDefault:"localhost"`
	DBPort			string `env:"DB_PORT" envDefault:"5432"`
	SessionSecret	string `env:"SESSION_SECRET" envDefault:""`
}
