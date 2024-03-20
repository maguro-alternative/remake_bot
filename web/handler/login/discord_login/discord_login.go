package discordlogin

import (
	"encoding/gob"
	"net/http"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

type DiscordOAuth2Handler struct {
	DiscordOAuth2Service *service.DiscordOAuth2Service
}

func NewDiscordOAuth2Handler(discordOAuth2Service *service.DiscordOAuth2Service) *DiscordOAuth2Handler {
	return &DiscordOAuth2Handler{
		DiscordOAuth2Service: discordOAuth2Service,
	}
}

// stateを生成し、Discordの認可ページのURLにリダイレクトする
func (h *DiscordOAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&model.DiscordUser{})
	uuid := uuid.New().String()
	session, err := h.DiscordOAuth2Service.CookieStore.Get(r, config.SessionSecret())
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.Values["discord_state"] = uuid
	// セッションに保存
	session.Save(r, w)
	h.DiscordOAuth2Service.CookieStore.Save(r, w, session)
	conf := h.DiscordOAuth2Service.OAuth2Conf
	// 1. 認可ページのURL
	url := conf.AuthCodeURL(uuid, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusSeeOther)
}