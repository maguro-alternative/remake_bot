package discordcallback

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"golang.org/x/oauth2"

	"github.com/maguro-alternative/remake_bot/web/handler/callback/discord_callback/internal"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
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
	var user internal.DiscordUser
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&internal.DiscordUser{})
	gob.Register(&oauth2.Token{})
	session, err := h.svc.CookieStore.Get(r, config.SessionSecret())
	if err != nil {
		slog.InfoContext(ctx, "sessionの取得に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	state, ok := session.Values["state"].(string)
	if !ok {
		stateType := reflect.TypeOf(session.Values["state"]).String()
		slog.InfoContext(ctx, stateType)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 2. 認可ページからリダイレクトされてきたときに送られてくるstateパラメータ
	if r.URL.Query().Get("state") != state {
		slog.InfoContext(ctx, "stateが一致しません。")
		session.Values["state"] = ""
		h.svc.CookieStore.Save(r, w, session)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.Values["state"] = ""
	// 1. 認可ページのURL
	code := r.URL.Query().Get("code")
	conf := h.svc.OAuth2Conf
	// 2. アクセストークンの取得
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		slog.InfoContext(ctx, "アクセストークンの取得に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.Values["discord_oauth_token"] = &token
	// 3. ユーザー情報の取得
	client := conf.Client(ctx, token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		slog.InfoContext(ctx, "ユーザー情報の取得に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		slog.InfoContext(ctx, "ユーザー情報のデコードに失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// セッションに保存
	session.Values["discord_user"] = user
	err = session.Save(r, w)
	if err != nil {
		slog.InfoContext(ctx, "セッションの保存に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.svc.CookieStore.Save(r, w, session)
	if err != nil {
		slog.InfoContext(ctx, "セッションの保存に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.InfoContext(ctx, fmt.Sprintf("ユーザー情報: %+v", user))
	// 4. ログイン後のページに遷移
	http.Redirect(w, r, "/", http.StatusFound)
}