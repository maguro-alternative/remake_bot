package discordlogin

import (
	"encoding/gob"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"
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
	sessionsSession, err := h.DiscordOAuth2Service.CookieStore.Get(r, config.SessionSecret())
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//sessionsSession.Values["discord_state"] = uuid
	session.SetDiscordState(sessionsSession, uuid)
	// セッションに保存
	err = sessionsSession.Save(r, w)
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.DiscordOAuth2Service.CookieStore.Save(r, w, sessionsSession)
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	conf := h.DiscordOAuth2Service.OAuth2Conf
	// 1. 認可ページのURL
	url := conf.AuthCodeURL(uuid, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
