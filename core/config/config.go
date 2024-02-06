package config

import (
	"fmt"
	"sync"

	"github.com/maguro-alternative/remake_bot/core/config/internal"

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

func DiscordBotToken() string {
	return cfg.DiscordBotToken
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

func DatabaseURLWithUserAndPassword() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", DatabaseUser(), DatabasePassword(), DatabaseName(), DatabaseHost(), DatabasePort())
}

func DatabaseURLWithoutUserAndPassword() string {
	return fmt.Sprintf("dbname=%s host=%s port=%s sslmode=disable", DatabaseName(), DatabaseHost(), DatabasePort())
}

func DatabaseURLWithoutUserAndPasswordForMigration() string {
	return fmt.Sprintf("dbname=%s host=%s port=%s sslmode=disable", "postgres", "localhost", "5432")
}

func DatabaseURLWithUserAndPasswordForMigration() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", "postgres", "", "postgres", "localhost", "5432")
}

func DatabaseURLForMigration() string {
	if DatabaseUser() == "" && DatabasePassword() == "" {
		return DatabaseURLWithoutUserAndPasswordForMigration()
	}
	return DatabaseURLWithUserAndPasswordForMigration()
}

func SessionSecret() string {
	return cfg.SessionSecret
}
