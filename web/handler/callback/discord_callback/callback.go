package discordcallback

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
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
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&model.DiscordUser{})
	sessionsSession, err := h.svc.CookieStore.Get(r, config.SessionSecret())
	if err != nil {
		slog.ErrorContext(ctx, "sessionの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	state, ok := sessionsSession.Values["discord_state"].(string)
	if !ok {
		stateType := reflect.TypeOf(sessionsSession.Values["discord_state"]).String()
		slog.ErrorContext(ctx, stateType+"型のstateが取得できませんでした。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 2. 認可ページからリダイレクトされてきたときに送られてくるstateパラメータ
	if r.URL.Query().Get("state") != state {
		slog.ErrorContext(ctx, "stateが一致しません。")
		sessionsSession.Values["discord_state"] = ""
		h.svc.CookieStore.Save(r, w, sessionsSession)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sessionsSession.Values["discord_state"] = ""
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
	sessionsSession.Values["discord_oauth_token"] = token.AccessToken
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
	sessionsSession.Values["discord_user"] = user
	err = sessionsSession.Save(r, w)
	if err != nil {
		slog.ErrorContext(ctx, "セッションの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.svc.CookieStore.Save(r, w, sessionsSession)
	if err != nil {
		slog.ErrorContext(ctx, "セッションの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.InfoContext(ctx, fmt.Sprintf("ユーザー情報: %+v", user))
	// 4. ログイン後のページに遷移
	http.Redirect(w, r, "/guilds", http.StatusFound)
}
