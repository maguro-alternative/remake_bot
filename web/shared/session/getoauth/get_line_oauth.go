package getoauth

import (
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

func init() {
	// 本番では削除すること
	gob.Register(&model.LineIdTokenUser{})
}

func (o *OAuthStore) GetLineOAuth(r *http.Request) (*model.LineOAuthSession, error) {
	session, err := o.Store.Get(r, o.Secret)
	if err != nil {
		return nil, err
	}
	// セッションに保存されているlineuserを取得
	lineUser, ok := session.Values["line_user"].(*model.LineIdTokenUser)
	if !ok {
		return nil, errors.New("session not found")
	}
	lineToken, ok := session.Values["line_oauth_token"].(string)
	if !ok {
		return nil, errors.New("session not found")
	}
	guildId, ok := session.Values["guild_id"].(string)
	if !ok {
		return nil, errors.New("session not found")
	}
	lineSession := model.LineOAuthSession{
		Token:          lineToken,
		DiscordGuildID: guildId,
		User:           *lineUser,
	}
	return &lineSession, nil
}
