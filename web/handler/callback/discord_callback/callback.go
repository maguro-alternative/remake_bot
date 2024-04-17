package discordcallback

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"
)

type DiscordCallbackHandler struct {
	svc *service.DiscordOAuth2Service
}

func NewDiscordCallbackHandler(svc *service.DiscordOAuth2Service) *DiscordCallbackHandler {
	return &DiscordCallbackHandler{
		svc: svc,
	}
}

func (h *DiscordCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var user model.DiscordUser
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	sessionStore, err := session.NewSessionStore(r, h.svc.CookieStore, config.SessionSecret())
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	state, err := sessionStore.GetDiscordState()
	if err != nil {
		slog.ErrorContext(ctx, "stateの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 2. 認可ページからリダイレクトされてきたときに送られてくるstateパラメータ
	if r.URL.Query().Get("state") != state {
		slog.ErrorContext(ctx, "stateが一致しません。", "state:", state, "r.URL.Query()", r.URL.Query().Get("state"))
		sessionStore.CleanupDiscordState()
		err = sessionStore.StoreSave(r, w, h.svc.CookieStore)
		if err != nil {
			slog.ErrorContext(ctx, "セッションの保存に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sessionStore.CleanupDiscordState()
	// 1. 認可ページのURL
	code := r.URL.Query().Get("code")
	conf := h.svc.OAuth2Conf
	// 2. アクセストークンの取得
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		slog.ErrorContext(ctx, "アクセストークンの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sessionStore.SetDiscordOAuthToken(token.AccessToken)
	// 3. ユーザー情報の取得
	client := conf.Client(ctx, token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		slog.ErrorContext(ctx, "ユーザー情報の取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		slog.ErrorContext(ctx, "ユーザー情報のデコードに失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// セッションに保存
	sessionStore.SetDiscordUser(&user)
	err = sessionStore.SessionSave(r, w)
	if err != nil {
		slog.ErrorContext(ctx, "セッションの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = sessionStore.StoreSave(r, w, h.svc.CookieStore)
	if err != nil {
		slog.ErrorContext(ctx, "セッションの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.InfoContext(ctx, fmt.Sprintf("ユーザー情報: %+v", user))
	// 4. ログイン後のページに遷移
	http.Redirect(w, r, "/guilds", http.StatusFound)
}
