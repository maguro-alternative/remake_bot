package internal

type Config struct {
	PrivateKey    string `env:"PRIVATE_KEY" envDefault:""`
}
