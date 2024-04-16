package linelogout

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session"
)

type LineLogoutHandler struct {
	LineLogoutService *service.IndexService
}

func NewLineLogoutHandler(lineLogoutService *service.IndexService) *LineLogoutHandler {
	return &LineLogoutHandler{
		LineLogoutService: lineLogoutService,
	}
}

// Lineの認証情報を削除し、ログアウトする
func (h *LineLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	sessionStore, err := session.NewSessionStore(r, h.LineLogoutService.CookieStore, config.SessionSecret())
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sessionStore.CleanupLineOAuthToken()
	sessionStore.CleanupGuildID()
	sessionStore.CleanupLineUser()
	err = sessionStore.SessionSave(r, w)
	if err != nil {
		slog.ErrorContext(ctx, "セッションの初期化に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = sessionStore.StoreSave(r, w, h.LineLogoutService.CookieStore)
	if err != nil {
		slog.ErrorContext(ctx, "セッションの初期化に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
