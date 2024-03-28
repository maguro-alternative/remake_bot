package middleware

import (
	"encoding/hex"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/ctxvalue"
	"github.com/maguro-alternative/remake_bot/pkg/line"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
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
) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var lineProfile line.LineProfile
			var lineLoginUser *model.LineOAuthSession
			ctx := r.Context()
			guildId := r.PathValue("guildId")
			oauthStore := getoauth.NewOAuthStore(indexService.CookieStore, config.SessionSecret())
			// ログインユーザーの取得
			lineLoginUser, err := oauthStore.GetLineOAuth(r)
			if err != nil {
				slog.ErrorContext(ctx, "lineのユーザー取得に失敗しました", "エラー:", err.Error())
				return
			}
			ctx = ctxvalue.ContextWithLineUser(ctx, lineLoginUser)

			lineBotApi, err := repo.GetLineBotNotClient(ctx, guildId)
			if err != nil {
				slog.ErrorContext(ctx, "lineBotの取得に失敗しました", "エラー:", err.Error())
				return
			}
			lineBotIv, err := repo.GetLineBotIvNotClient(ctx, guildId)
			if err != nil {
				slog.ErrorContext(ctx, "lineBotIvの取得に失敗しました", "エラー:", err.Error())
				return
			}
			var lineBotDecrypt LineBotDecrypt
			// 暗号化キーのバイトへの変換
			keyBytes, err := hex.DecodeString(config.PrivateKey())
			if err != nil {
				slog.ErrorContext(ctx, "暗号化キーのバイトへの変換に失敗しました", "エラー:", err.Error())
				return
			}

			lineNotifyTokenByte, err := crypto.Decrypt(lineBotApi.LineNotifyToken[0], keyBytes, lineBotIv.LineNotifyTokenIv[0])
			if err != nil {
				slog.ErrorContext(ctx, "LineNotifyTokenの復号に失敗しました", "エラー:", err.Error())
				return
			}
			lineBotTokenByte, err := crypto.Decrypt(lineBotApi.LineBotToken[0], keyBytes, lineBotIv.LineBotTokenIv[0])
			if err != nil {
				slog.ErrorContext(ctx, "LineBotTokenの復号に失敗しました", "エラー:", err.Error())
				return
			}
			lineGroupIDByte, err := crypto.Decrypt(lineBotApi.LineGroupID[0], keyBytes, lineBotIv.LineGroupIDIv[0])
			if err != nil {
				slog.ErrorContext(ctx, "LineGuildIDの復号に失敗しました", "エラー:", err.Error())
				return
			}
			lineBotDecrypt.LineNotifyToken = string(lineNotifyTokenByte)
			lineBotDecrypt.LineBotToken = string(lineBotTokenByte)
			lineBotDecrypt.LineGroupID = string(lineGroupIDByte)
			lineBotDecrypt.DefaultChannelID = lineBotApi.DefaultChannelID
			lineBotDecrypt.DebugMode = lineBotApi.DebugMode

			lineRequ := line.NewLineRequest(
				lineBotDecrypt.LineNotifyToken,
				lineBotDecrypt.LineBotToken,
				lineBotDecrypt.LineGroupID,
			)
			lineProfile, err = lineRequ.GetProfileInGroup(ctx, lineLoginUser.User.Sub)
			if err != nil {
				slog.ErrorContext(ctx, "LineProfileの取得に失敗しました", "エラー:", err.Error())
				return
			}
			ctx = ctxvalue.ContextWithLineProfile(ctx, &lineProfile)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
