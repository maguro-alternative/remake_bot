package getoauth

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

type DiscordOAuthSession struct {
	Token string      `json:"token"`
	User  DiscordUser `json:"user"`
}

type DiscordToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
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

func DiscordOAuthCheck(store *sessions.CookieStore, r *http.Request, sessionSecret string) (*DiscordUser, error) {
	session, err := store.Get(r, sessionSecret)
	if err != nil {
		return nil, err
	}
	// セッションに保存されているdiscorduserを取得
	discordUser, ok := session.Values["discord_user"].(*DiscordUser)
	if !ok {
		return nil, errors.New("session not found")
	}
	return discordUser, nil
}
