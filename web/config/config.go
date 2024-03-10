package config

import (
	"fmt"
	"sync"

	"github.com/maguro-alternative/remake_bot/web/config/internal"

	"github.com/caarlos0/env/v7"
	"github.com/cockroachdb/errors"
	"github.com/joho/godotenv"
)

var (
	once sync.Once
	cfg  *internal.Config
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	once.Do(MustInit)
}

func MustInit() {
	cfg = &internal.Config{}
	if err := env.Parse(cfg); err != nil {
		xerr := errors.Wrap(err, "failed to env parse: ")
		fmt.Printf("panic: %+v", xerr)
		panic(xerr)
	}
}

func CsrfAuthKey() string {
	return cfg.CsrfAuthKey
}

func DatabaseName() string {
	return cfg.DBName
}

func DatabaseUser() string {
	return cfg.DBUser
}

func DatabasePassword() string {
	return cfg.DBPassword
}

func DatabaseHost() string {
	return cfg.DBHost
}

func DatabasePort() string {
	return cfg.DBPort
}

func DatabaseURL() string {
	return fmt.Sprintf("%s://%s:%s/%s?user=%s&password=%s&sslmode=disable", cfg.DBName, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword)
}


func DiscordClientID() string {
	return cfg.DiscordClientID
}

func DiscordClientSecret() string {
	return cfg.DiscordClientSecret
}

func DiscordCallbackUrl() string {
	return cfg.DiscordCallbackUrl
}

func DiscordScopes() string {
	return cfg.DiscordScopes
}

func PrivateKey() string {
	return cfg.PrivateKey
}

func Port() string {
	return cfg.Port
}

func ServerUrl() string {
	return cfg.ServerUrl
}

func SessionName() string {
	return cfg.SessionName
}

func SessionSecret() string {
	return cfg.SessionSecret
}

func YouTubeAPIKey() string {
	return cfg.YouTubeAPIKey
}

func YoutubeAccessToken() string {
	return cfg.YoutubeAccessToken
}

func YoutubeClientID() string {
	return cfg.YoutubeClientID
}

func YoutubeClientSecret() string {
	return cfg.YoutubeClientSecret
}

func YoutubeRefreshToken() string {
	return cfg.YoutubeRefreshToken
}

func YoutubeProjectID() string {
	return cfg.YoutubeProjectID
}

func YoutubeTokenExpiry() string {
	return cfg.YoutubeTokenExpiry
}
