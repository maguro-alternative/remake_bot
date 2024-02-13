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

func PrivateKey() string {
	return cfg.PrivateKey
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
