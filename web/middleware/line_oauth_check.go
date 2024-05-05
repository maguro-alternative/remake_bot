package middleware

import (
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/line"
	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"
)

type LineBotDecrypt struct {
	LineNotifyToken  string
	LineBotToken     string
	LineGroupID      string
	LineClientID     string
	LineClientSecret string
	DefaultChannelID string
	DebugMode        bool
}

func LineOAuthCheckMiddleware(
	indexService service.IndexService,
	repo Repository,
	aesCrypto crypto.AESInterface,
	loginRequiredFlag bool,
) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var lineProfile line.LineProfile
			var lineLoginUser *model.LineOAuthSession
			ctx := r.Context()

			sessionStore, err := session.NewSessionStore(r, indexService.CookieStore, config.SessionSecret())
			if err != nil {
				slog.ErrorContext(r.Context(), "sessionの取得に失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// ログインユーザーの取得
			lineUser, err := sessionStore.GetLineUser()
			if err != nil && loginRequiredFlag {
				slog.ErrorContext(ctx, "lineのユーザー取得に失敗しました", "エラー:", err.Error())
				http.Redirect(w, r, "/login/line", http.StatusFound)
				return
			}
			sessionGuildId, err := sessionStore.GetGuildID()
			if err != nil && loginRequiredFlag {
				slog.ErrorContext(ctx, "guildIdの取得に失敗しました", "エラー:", err.Error())
				http.Redirect(w, r, "/login/line", http.StatusFound)
				return
			}
			lineToken, err := sessionStore.GetLineOAuthToken()
			if err != nil && loginRequiredFlag {
				slog.ErrorContext(ctx, "lineのトークン取得に失敗しました", "エラー:", err.Error())
				http.Redirect(w, r, "/login/line", http.StatusFound)
				return
			}

			// 認証情報が必要ないかつ、認証情報がない場合
			if lineUser == nil {
				lineLoginUser = &model.LineOAuthSession{
					User: model.LineIdTokenUser{},
				}
				ctx = ctxvalue.ContextWithLineProfile(ctx, &lineProfile)
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			lineLoginUser = &model.LineOAuthSession{
				User:           *lineUser,
				Token:          lineToken,
				DiscordGuildID: sessionGuildId,
			}
			ctx = ctxvalue.ContextWithLineUser(ctx, lineLoginUser)

			lineBotApi, err := repo.GetLineBotNotClient(ctx, sessionGuildId)
			if err != nil {
				slog.ErrorContext(ctx, "lineBotの取得に失敗しました", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			lineBotIv, err := repo.GetLineBotIvNotClient(ctx, sessionGuildId)
			if err != nil {
				slog.ErrorContext(ctx, "lineBotIvの取得に失敗しました", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			var lineBotDecrypt LineBotDecrypt

			lineNotifyTokenByte, err := aesCrypto.Decrypt(lineBotApi.LineNotifyToken[0], lineBotIv.LineNotifyTokenIv[0])
			if err != nil {
				slog.ErrorContext(ctx, "LineNotifyTokenの復号に失敗しました", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			lineBotTokenByte, err := aesCrypto.Decrypt(lineBotApi.LineBotToken[0], lineBotIv.LineBotTokenIv[0])
			if err != nil {
				slog.ErrorContext(ctx, "LineBotTokenの復号に失敗しました", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			lineGroupIDByte, err := aesCrypto.Decrypt(lineBotApi.LineGroupID[0], lineBotIv.LineGroupIDIv[0])
			if err != nil {
				slog.ErrorContext(ctx, "LineGuildIDの復号に失敗しました", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			lineBotDecrypt.LineNotifyToken = string(lineNotifyTokenByte)
			lineBotDecrypt.LineBotToken = string(lineBotTokenByte)
			lineBotDecrypt.LineGroupID = string(lineGroupIDByte)
			lineBotDecrypt.DefaultChannelID = lineBotApi.DefaultChannelID
			lineBotDecrypt.DebugMode = lineBotApi.DebugMode

			lineRequ := line.NewLineRequest(
				*indexService.Client,
				lineBotDecrypt.LineNotifyToken,
				lineBotDecrypt.LineBotToken,
				lineBotDecrypt.LineGroupID,
			)
			lineProfile, err = lineRequ.GetProfileInGroup(ctx, lineLoginUser.User.Sub)
			if err != nil && loginRequiredFlag {
				slog.ErrorContext(ctx, "LineProfileの取得に失敗しました", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			ctx = ctxvalue.ContextWithLineProfile(ctx, &lineProfile)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
