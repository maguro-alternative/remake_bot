package permission

import (
	"context"
	"encoding/hex"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/line"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/permission/internal"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

func CheckLinePermission(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	indexService *service.IndexService,
	guildId string,
) (lineProfile line.LineProfile, lineLoginUser *model.LineOAuthSession, err error) {
	repo := internal.NewRepository(indexService.DB)

	// ログインユーザーの取得
	lineLoginUser, err = getoauth.GetLineOAuth(
		indexService.CookieStore,
		r,
		config.SessionSecret(),
	)
	if err != nil {
		slog.ErrorContext(ctx, "lineのユーザー取得に失敗しました", "エラー:", err.Error())
		return lineProfile, lineLoginUser, err
	}

	lineBotApi, err := repo.GetLineBot(ctx, guildId)
	if err != nil {
		slog.ErrorContext(ctx, "lineBotの取得に失敗しました", "エラー:", err.Error())
		return lineProfile, lineLoginUser, err
	}
	lineBotIv, err := repo.GetLineBotIv(ctx, guildId)
	if err != nil {
		slog.ErrorContext(ctx, "lineBotIvの取得に失敗しました", "エラー:", err.Error())
		return lineProfile, lineLoginUser, err
	}
	var lineBotDecrypt internal.LineBotDecrypt
	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(config.PrivateKey())
	if err != nil {
		slog.ErrorContext(ctx, "暗号化キーのバイトへの変換に失敗しました", "エラー:", err.Error())
		return lineProfile, lineLoginUser, err
	}

	lineNotifyTokenByte, err := crypto.Decrypt(lineBotApi.LineNotifyToken[0], keyBytes, lineBotIv.LineNotifyTokenIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "LineNotifyTokenの復号に失敗しました", "エラー:", err.Error())
		return lineProfile, lineLoginUser, err
	}
	lineBotTokenByte, err := crypto.Decrypt(lineBotApi.LineBotToken[0], keyBytes, lineBotIv.LineBotTokenIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "LineBotTokenの復号に失敗しました", "エラー:", err.Error())
		return lineProfile, lineLoginUser, err
	}
	lineGroupIDByte, err := crypto.Decrypt(lineBotApi.LineGroupID[0], keyBytes, lineBotIv.LineGroupIDIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "LineGuildIDの復号に失敗しました", "エラー:", err.Error())
		return lineProfile, lineLoginUser, err
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
		return lineProfile, lineLoginUser, err
	}
	return lineProfile, lineLoginUser, err

}

