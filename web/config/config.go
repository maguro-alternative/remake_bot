package config

import (
	"fmt"
	"strings"
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

func DatabaseURLWithSslmode() string {
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

func LineCallBackUrl() string {
	return cfg.LineCallBackUrl
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

// 開発用の既定値。本番でこの値のまま起動するとセッション偽造やトークン復号が可能になる
const (
	insecureDefaultPrivateKey    = "645E739A7F9F162725C1533DC2C5E827"
	insecureDefaultSessionSecret = "test"
	minSessionSecretLength       = 32
)

// ValidateSecrets は秘密鍵が既定値や弱い値のままでないかを検証します。
// 起動時に呼び出し、エラーの場合は起動を中止してください。
func ValidateSecrets() error {
	if cfg.PrivateKey == "" || cfg.PrivateKey == insecureDefaultPrivateKey {
		return errors.New("PRIVATE_KEY が未設定または既定値です。ランダムな値を設定してください")
	}
	if cfg.SessionSecret == "" || cfg.SessionSecret == insecureDefaultSessionSecret {
		return errors.New("SESSION_SECRET が未設定または既定値です。ランダムな値を設定してください")
	}
	if len(cfg.SessionSecret) < minSessionSecretLength {
		return errors.Newf("SESSION_SECRET は %d 文字以上にしてください", minSessionSecretLength)
	}
	return nil
}

// IsSecureServer はサーバーURLがHTTPSかどうかを返します。
func IsSecureServer() bool {
	return strings.HasPrefix(cfg.ServerUrl, "https://")
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
