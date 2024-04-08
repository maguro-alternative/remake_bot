package linecallback

import (
	"context"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

type LineCallbackHandler struct {
	svc  *service.IndexService
	repo repository.RepositoryFunc
}

func NewLineCallbackHandler(
	svc *service.IndexService,
	repo repository.RepositoryFunc,
) *LineCallbackHandler {
	return &LineCallbackHandler{
		svc:  svc,
		repo: repo,
	}
}

func (h *LineCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	privateKey := config.PrivateKey()
	keyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		slog.ErrorContext(ctx, "暗号化キーのバイトへの変換に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&model.LineIdTokenUser{})
	session, err := h.svc.CookieStore.Get(r, config.SessionSecret())
	if err != nil {
		slog.ErrorContext(ctx, "sessionの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	state, ok := session.Values["line_state"].(string)
	if !ok {
		stateType := reflect.TypeOf(session.Values["line_state"]).String()
		slog.ErrorContext(ctx, stateType+"型のstateが取得できませんでした。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	nonce, ok := session.Values["line_nonce"].(string)
	if !ok {
		slog.ErrorContext(ctx, "nonceが取得できませんでした。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	guildId, ok := session.Values["guild_id"].(string)
	if !ok {
		slog.ErrorContext(ctx, "guild_idが取得できませんでした。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 2. 認可ページからリダイレクトされてきたときに送られてくるstateパラメータ
	if r.URL.Query().Get("state") != state {
		slog.ErrorContext(ctx, "stateが一致しません。")
		session.Values["line_state"] = ""
		h.svc.CookieStore.Save(r, w, session)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.Values["line_state"] = ""
	session.Values["line_nonce"] = ""
	// 1. 認可ページのURL
	code := r.URL.Query().Get("code")
	lineBot, err := h.repo.GetAllColumnsLineBot(ctx, guildId)
	if err != nil {
		slog.ErrorContext(ctx, "line_botの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineBotIv, err := h.repo.GetAllColumnsLineBotIv(ctx, lineBot.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_bot_ivの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineClientIDByte, err := crypto.Decrypt(lineBot.LineClientID[0], keyBytes, lineBotIv.LineClientIDIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_client_idの復号に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineClientSecretByte, err := crypto.Decrypt(lineBot.LineClientSecret[0], keyBytes, lineBotIv.LineClientSecretIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_client_secretの復号に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineClientID := string(lineClientIDByte)
	lineClientSecret := string(lineClientSecretByte)
	// 3. アクセストークンを取得するためのリクエスト
	token, cleanupTokenBody, err := getIdToken(ctx, h.svc.Client, code, lineClientID, lineClientSecret)
	if err != nil {
		slog.ErrorContext(ctx, "ユーザー情報の取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cleanupTokenBody()
	user, cleanupIdtokenBody, err := verifyIdToken(ctx, h.svc.Client, token.IDToken, lineClientID, nonce)
	if err != nil {
		slog.ErrorContext(ctx, "id_tokenの検証に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cleanupIdtokenBody()
	session.Values["line_user"] = user
	session.Values["line_oauth_token"] = token.AccessToken
	err = session.Save(r, w)
	if err != nil {
		slog.ErrorContext(ctx, "sessionの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.svc.CookieStore.Save(r, w, session)
	if err != nil {
		slog.ErrorContext(ctx, "sessionの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, config.ServerUrl()+"/group/"+guildId, http.StatusFound)
}

func getIdToken(ctx context.Context, client *http.Client, code, clientID, clientSecret string) (*model.LineToken, func(), error) {
	u, err := url.ParseRequestURI("https://api.line.me/oauth2/v2.1/token")
	if err != nil {
		return nil, func() {}, err
	}
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", config.ServerUrl()+"/callback/line-callback/")
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
	var token model.LineToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, func() {}, err
	}
	return &token, func() { resp.Body.Close() }, nil
}

func verifyIdToken(ctx context.Context, nonceClient *http.Client, idToken, clientID, nonce string) (*model.LineIdTokenUser, func(), error) {
	verifyUrl := "https://api.line.me/oauth2/v2.1/verify"
	u, err := url.ParseRequestURI(verifyUrl)
	if err != nil {
		return nil, func() {}, err
	}
	form := url.Values{}
	form.Add("id_token", idToken)
	form.Add("client_id", clientID)
	form.Add("nonce", nonce)

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, func() {}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := nonceClient.Do(req)
	if err != nil {
		return nil, func() {}, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.InfoContext(ctx, resp.Status)
		slog.InfoContext(ctx, form.Encode())
		var e struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&e)
		return nil, func() {}, errors.New(e.ErrorDescription)
	}
	var user model.LineIdTokenUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, func() {}, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("ユーザー情報: %+v", user))
	return &user, func() { resp.Body.Close() }, nil
}
