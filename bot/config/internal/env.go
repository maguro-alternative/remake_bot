package internal

type Config struct {
	Token string `env:"TOKEN" envDefault:""`
}