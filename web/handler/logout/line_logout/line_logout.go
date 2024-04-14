package linelogout

import (
	"context"
	"encoding/gob"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
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
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	gob.Register(&model.LineOAuthSession{})
	session, err := h.LineLogoutService.CookieStore.Get(r, config.SessionSecret())
	if err != nil {
		slog.ErrorContext(ctx, "sessionの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.Values["line_oauth_token"] = ""
	session.Values["guild_id"] = ""
	session.Values["line_user"] = model.LineOAuthSession{}
	err = session.Save(r, w)
	if err != nil {
		slog.ErrorContext(ctx, "セッションの初期化に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.LineLogoutService.CookieStore.Save(r, w, session)
	if err != nil {
		slog.ErrorContext(ctx, "セッションの初期化に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
