package discordlogout

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
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

// Discordの認証情報を削除し、ログアウトする
func (h *DiscordOAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	sessionStore, err := session.NewSessionStore(r, h.DiscordOAuth2Service.CookieStore, config.SessionSecret())
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sessionStore.CleanupDiscordUser()
	sessionStore.CleanupDiscordOAuthToken()
	err = sessionStore.SessionSave(r, w)
	if err != nil {
		slog.ErrorContext(ctx, "セッションの初期化に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.DiscordOAuth2Service.CookieStore.Save(r, w, sessionStore.GetSession())
	if err != nil {
		slog.ErrorContext(ctx, "セッションの初期化に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
