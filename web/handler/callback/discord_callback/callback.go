package discordcallback

import (
	"context"
	"encoding/json"
	"encoding/gob"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"
)

func init(){
	gob.Register(&model.DiscordUser{})
}

type DiscordCallbackHandler struct {
	svc *service.IndexService
}

func NewDiscordCallbackHandler(svc *service.IndexService) *DiscordCallbackHandler {
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
	//conf := h.svc.OAuth2Conf
	// 2. アクセストークンの取得
	token, cleanupTokenBody, err := getToken(ctx, h.svc.Client, code, config.DiscordClientID(), config.DiscordClientSecret())
	if err != nil {
		slog.ErrorContext(ctx, "アクセストークンの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cleanupTokenBody()
	sessionStore.SetDiscordOAuthToken(token.AccessToken)
	// 3. ユーザー情報の取得
	resp, err := h.svc.Client.Get("https://discord.com/api/users/@me")
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

func getToken(ctx context.Context, client *http.Client, code, clientID, clientSecret string) (*model.DiscordToken, func(), error) {
	u, err := url.ParseRequestURI("https://discord.com/api/oauth2/token")
	if err != nil {
		return nil, func() {}, err
	}
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", config.ServerUrl()+"/callback/discord-callback/")
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)
	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, func() {}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, func() {}, err
	}
	var token model.DiscordToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, func() {}, err
	}
	return &token, func() { resp.Body.Close() }, nil
}

