package discordlogin

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session"
)

type DiscordOAuth2Handler struct {
	indexService *service.IndexService
}

func NewDiscordOAuth2Handler(indexService *service.IndexService) *DiscordOAuth2Handler {
	return &DiscordOAuth2Handler{
		indexService: indexService,
	}
}

// stateを生成し、Discordの認可ページのURLにリダイレクトする
func (h *DiscordOAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New().String()
	sessionStore, err := session.NewSessionStore(r, h.indexService.CookieStore, config.SessionSecret())
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sessionStore.SetDiscordState(uuid)
	// セッションに保存
	err = sessionStore.SessionSave(r, w)
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = sessionStore.StoreSave(r, w, h.indexService.CookieStore)
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	redirectUri := config.ServerUrl() + "/callback/discord-callback/"
	// 1. 認可ページのURL
	url := fmt.Sprintf("https://discord.com/api/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=%s", config.DiscordClientID(), redirectUri, uuid, config.DiscordScopes())
	http.Redirect(w, r, url, http.StatusSeeOther)
}
