package internal

type Config struct {
	PrivateKey      string `env:"PRIVATE_KEY" envDefault:"645E739A7F9F162725C1533DC2C5E827"`
	DiscordBotToken string `env:"DISCORD_BOT_TOKEN" envDefault:""`
	DBName          string `env:"DB_NAME" envDefault:"postgres"`
	DBUser			string `env:"DB_USER" envDefault:"postgres"`
	DBPassword		string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBHost			string `env:"DB_HOST" envDefault:"localhost"`
	DBPort			string `env:"DB_PORT" envDefault:"5432"`
	SessionSecret	string `env:"SESSION_SECRET" envDefault:""`
}
