package linecallback

import (
	"context"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"strings"
	"log/slog"
	"net/http"
	"reflect"

	"golang.org/x/oauth2"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/handler/callback/line_callback/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

type LineCallbackHandler struct {
	svc *service.IndexService
}

func NewLineCallbackHandler(svc *service.IndexService) *LineCallbackHandler {
	return &LineCallbackHandler{
		svc: svc,
	}
}

func (h *LineCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var user model.LineUser
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	privateKey := config.PrivateKey()
	keyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		slog.ErrorContext(ctx, "暗号化キーのバイトへの変換に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	repo := internal.NewRepository(h.svc.DB)
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&model.LineUser{})
	session, err := h.svc.CookieStore.Get(r, config.SessionSecret())
	if err != nil {
		slog.InfoContext(ctx, "sessionの取得に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	state, ok := session.Values["line_state"].(string)
	if !ok {
		stateType := reflect.TypeOf(session.Values["line_state"]).String()
		slog.InfoContext(ctx, stateType)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	guildId, ok := session.Values["guild_id"].(string)
	if !ok {
		slog.InfoContext(ctx, "guild_idが取得できませんでした。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 2. 認可ページからリダイレクトされてきたときに送られてくるstateパラメータ
	if r.URL.Query().Get("line_state") != state {
		slog.InfoContext(ctx, "stateが一致しません。")
		session.Values["line_state"] = ""
		h.svc.CookieStore.Save(r, w, session)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.Values["line_state"] = ""
	// 1. 認可ページのURL
	code := r.URL.Query().Get("code")
	lineBot, err := repo.GetLineBot(ctx, guildId)
	if err != nil {
		slog.InfoContext(ctx, "line_botの取得に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineBotIv, err := repo.GetLineBotIv(ctx, lineBot.GuildID)
	if err != nil {
		slog.InfoContext(ctx, "line_bot_ivの取得に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineClientIDByte, err := crypto.Decrypt(lineBot.LineClientID[0], keyBytes, lineBotIv.LineClientIDIv[0])
	if err != nil {
		slog.InfoContext(ctx, "line_client_idの復号に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineClientSecretByte, err := crypto.Decrypt(lineBot.LineClientSecret[0], keyBytes, lineBotIv.LineClientSecretIv[0])
	if err != nil {
		slog.InfoContext(ctx, "line_client_secretの復号に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineClientID := string(lineClientIDByte)
	lineClientSecret := string(lineClientSecretByte)
	// 3. アクセストークンを取得するためのリクエスト
	conf := &oauth2.Config{
		ClientID:     lineClientID,
		ClientSecret: lineClientSecret,
		Scopes:       strings.Split("profile%%20openid%%20email", "%%20"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
			TokenURL: "https://api.line.me/oauth2/v2.1/token",
		},
		RedirectURL: config.ServerUrl() + "/callback/line-callback/",
	}
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		slog.InfoContext(ctx, "アクセストークンの取得に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 4. アクセストークンを使ってユーザー情報を取得する
	client := conf.Client(ctx, token)
	resp, err := client.Get("https://api.line.me/v2/profile")
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
	session.Values["line_user"] = user
	session.Values["line_oauth_token"] = token.AccessToken
	err = session.Save(r, w)
	if err != nil {
		slog.InfoContext(ctx, "sessionの保存に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, config.ServerUrl()+"/group/"+guildId, http.StatusFound)
}
