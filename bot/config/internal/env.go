package internal

type Config struct {
	Port string `env:"PORT" envDefault:"8080"`
}