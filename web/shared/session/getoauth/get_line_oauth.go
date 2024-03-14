package getoauth

import (
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/shared/session/model"

	"github.com/gorilla/sessions"
)

func init() {
	// 本番では削除すること
	gob.Register(&model.LineIdTokenUser{})
}

func GetLineOAuth(store *sessions.CookieStore, r *http.Request, sessionSecret string) (*model.LineOAuthSession, error) {
	session, err := store.Get(r, sessionSecret)
	if err != nil {
		return nil, err
	}
	// セッションに保存されているlineuserを取得
	lineUser, ok := session.Values["line_user"].(*model.LineIdTokenUser)
	if !ok {
		return nil, errors.New("session not found")
	}
	lineToken, ok := session.Values["line_oauth_token"].(*model.LineToken)
	if !ok {
		return nil, errors.New("session not found")
	}
	guildId, ok := session.Values["guild_id"].(string)
	if !ok {
		return nil, errors.New("session not found")
	}
	lineSession := model.LineOAuthSession{
		Token:          lineToken.AccessToken,
		DiscordGuildID: guildId,
		User:           *lineUser,
	}
	return &lineSession, nil
}
