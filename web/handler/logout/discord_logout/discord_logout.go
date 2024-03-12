package discordlogout

import (
	"context"
	"encoding/gob"
	"net/http"
	"log/slog"

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

// Discordの認証情報を削除し、ログアウトする
func (h *DiscordOAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	gob.Register(&model.DiscordUser{})
	session, err := h.DiscordOAuth2Service.CookieStore.Get(r, config.SessionSecret())
	if err != nil {
		slog.InfoContext(r.Context(), "sessionの取得に失敗しました。"+err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.Values["discord_oauth_token"] = ""
	session.Values["discord_user"] = model.DiscordUser{}
	err = session.Save(r, w)
	if err != nil {
		slog.InfoContext(ctx, "セッションの初期化に失敗しました。"+err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.DiscordOAuth2Service.CookieStore.Save(r, w, session)
	if err != nil {
		slog.InfoContext(ctx, "セッションの初期化に失敗しました。"+err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
