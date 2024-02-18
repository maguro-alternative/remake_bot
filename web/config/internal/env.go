package internal

type Config struct {
	DBName              string `env:"DB_NAME" envDefault:"postgres"`
	DBUser              string `env:"DB_USER" envDefault:"postgres"`
	DBPassword          string `env:"DB_PASSWORD" envDefault:""`
	DBHost              string `env:"DB_HOST" envDefault:"localhost"`
	DBPort              string `env:"DB_PORT" envDefault:"5432"`
	DiscordClientID     string `env:"DISCORD_CLIENT_ID" envDefault:""`
	DiscordClientSecret string `env:"DISCORD_CLIENT_SECRET" envDefault:""`
	DiscordCallbackUrl  string `env:"DISCORD_CALLBACK_URL" envDefault:""`
	DiscordScopes       string `env:"DISCORD_SCOPES" envDefault:""`
	PrivateKey          string `env:"PRIVATE_KEY" envDefault:""`
	SessionName         string `env:"SESSION_NAME" envDefault:""`
	SessionSecret       string `env:"SESSION_SECRET" envDefault:""`
	YouTubeAPIKey       string `env:"YOUTUBE_API_KEY" envDefault:""`
	YoutubeAccessToken  string `env:"YOUTUBE_ACCESS_TOKEN" envDefault:""`
	YoutubeClientID     string `env:"YOUTUBE_CLIENT_ID" envDefault:""`
	YoutubeClientSecret string `env:"YOUTUBE_CLIENT_SECRET" envDefault:""`
	YoutubeRefreshToken string `env:"YOUTUBE_REFRESH_TOKEN" envDefault:""`
	YoutubeProjectID    string `env:"YOUTUBE_PROJECT_ID" envDefault:""`
	YoutubeTokenExpiry  string `env:"YOUTUBE_TOKEN_EXPIRY" envDefault:""`
}
