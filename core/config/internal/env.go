package internal

type Config struct {
	DiscordBotToken   string `env:"DISCORD_BOT_TOKEN" envDefault:""`
	DBName            string `env:"DB_NAME" envDefault:"postgres"`
	DBUser            string `env:"DB_USER" envDefault:"postgres"`
	DBPassword        string `env:"DB_PASSWORD" envDefault:""`
	DBHost            string `env:"DB_HOST" envDefault:"localhost"`
	DBPort            string `env:"DB_PORT" envDefault:"5432"`
	DBURL             string `env:"DATABASE_URL" envDefault:""`
	SessionSecret     string `env:"SESSION_SECRET" envDefault:""`
	LineWorksID       string `env:"LINE_WORKS_ID" envDefault:""`
	LineWorksPassword string `env:"LINE_WORKS_PASSWORD" envDefault:""`
}
