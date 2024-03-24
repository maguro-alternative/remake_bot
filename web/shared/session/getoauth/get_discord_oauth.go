package getoauth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/shared/session/model"

	"github.com/gorilla/sessions"
)

type OAuthStore struct {
	Store *sessions.CookieStore
	Secret string
}

func NewOAuthStore(store *sessions.CookieStore, secret string) *OAuthStore {
	return &OAuthStore{
		Store: store,
		Secret: secret,
	}
}

func (o *OAuthStore) GetDiscordOAuth(ctx context.Context, r *http.Request) (*model.DiscordOAuthSession, error) {
	session, err := o.Store.Get(r, o.Secret)
	if err != nil {
		slog.ErrorContext(ctx, "sessionの取得に失敗しました。", "エラー:", err.Error())
		return nil, err
	}
	// セッションに保存されているdiscorduserを取得
	discordUser, ok := session.Values["discord_user"].(*model.DiscordUser)
	if !ok {
		slog.ErrorContext(ctx, "user session not found")
		return nil, errors.New("session not found")
	}
	// セッションに保存されているtokenを取得
	token, ok := session.Values["discord_oauth_token"].(string)
	if !ok {
		slog.ErrorContext(ctx, "token session not found")
		return nil, errors.New("session not found")
	}
	discordOAuth := &model.DiscordOAuthSession{
		Token: token,
		User:  *discordUser,
	}
	return discordOAuth, nil
}
