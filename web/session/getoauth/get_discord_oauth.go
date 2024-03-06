package getoauth

import (
	"errors"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/gorilla/sessions"
)

type DiscordOAuthSession struct {
	Token oauth2.Token `json:"token"`
	User  DiscordUser  `json:"user"`
}

type DiscordUser struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	GlobalName       string `json:"global_name"`
	DisplayName      string `json:"display_name"`
	Avatar           string `json:"avatar"`
	AvatarDecoration string `json:"avatar_decoration"`
	Discriminator    string `json:"discriminator"`
	PublicFlags      int    `json:"public_flags"`
	Flags            int    `json:"flags"`
	Banner           string `json:"banner"`
	BannerColor      string `json:"banner_color"`
	AccentColor      int    `json:"accent_color"`
	Locale           string `json:"locale"`
	MfaEnabled       bool   `json:"mfa_enabled"`
	PremiumType      int    `json:"premium_type"`
	Email            string `json:"email"`
	Verified         bool   `json:"verified"`
	Bio              string `json:"bio"`
}

func GetDiscordOAuth(store *sessions.CookieStore, r *http.Request, sessionSecret string) (*DiscordOAuthSession, error) {
	session, err := store.Get(r, sessionSecret)
	if err != nil {
		return nil, err
	}
	// セッションに保存されているdiscorduserを取得
	discordUser, ok := session.Values["discord_user"].(*DiscordOAuthSession)
	if !ok {
		return nil, errors.New("session not found")
	}
	return discordUser, nil
}
