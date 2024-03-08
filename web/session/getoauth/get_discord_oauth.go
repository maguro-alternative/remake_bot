package getoauth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/session/model"

	"github.com/gorilla/sessions"
)

func GetDiscordOAuth(ctx context.Context, store *sessions.CookieStore, r *http.Request, sessionSecret string) (*model.DiscordOAuthSession, error) {
	session, err := store.Get(r, sessionSecret)
	if err != nil {
		slog.InfoContext(ctx, "sessionの取得に失敗しました。"+err.Error())
		return nil, err
	}
	// セッションに保存されているdiscorduserを取得
	discordUser, ok := session.Values["discord_user"].(*model.DiscordUser)
	if !ok {
		slog.InfoContext(ctx, "user session not found")
		return nil, errors.New("session not found")
	}
	// セッションに保存されているtokenを取得
	token, ok := session.Values["discord_oauth_token"].(string)
	if !ok {
		slog.InfoContext(ctx, "token session not found")
		return nil, errors.New("session not found")
	}
	discordOAuth := &model.DiscordOAuthSession{
		Token: token,
		User:  *discordUser,
	}
	return discordOAuth, nil
}
