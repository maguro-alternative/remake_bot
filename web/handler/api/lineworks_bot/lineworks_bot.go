package lineworksbot

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/lineworks"

	"github.com/maguro-alternative/remake_bot/web/handler/api/lineworks_bot/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

// A LineBotHandler handles requests for the line bot.
type LineWorksHandler struct {
	indexService *service.IndexService
	repo         repository.RepositoryFunc
	aesCrypto crypto.AESInterface
}

// NewLineBotHandler returns new LineBotHandler.
func NewLineWorksHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
	aesCrypto crypto.AESInterface,
) *LineWorksHandler {
	return &LineWorksHandler{
		indexService: indexService,
		repo:         repo,
		aesCrypto:    aesCrypto,
	}
}

func (h *LineWorksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var lineWorksResponses internal.LineWorksResponses
	var lineWorksBotDecrypt *internal.LineWorksBotDecrypt
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	lineWorksBots, err := h.repo.GetAllLineWorksBots(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "line_works_botの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if len(lineWorksBots) == 0 {
		slog.ErrorContext(ctx, "line_works_botが存在しません。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// リクエストボディの読み込み
	requestBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "リクエストボディの読み込みに失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// リクエストボディの復号化
	for i, lineWorksBot := range lineWorksBots {
		lineWorksBotIv, err := h.repo.GetLineWorksBotIVByGuildID(ctx, lineWorksBot.GuildID)
		if err != nil {
			slog.ErrorContext(ctx, "line_works_bot_ivの取得に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		lineWorksBotSecretKey, err := h.aesCrypto.Decrypt(lineWorksBot.LineWorksBotSecret[0], lineWorksBotIv.LineWorksBotSecretIV[0])
		if err != nil {
			slog.ErrorContext(ctx, "リクエストボディの復号化に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if internal.LineWorksValidateRequest(requestBodyBytes, r.Header.Get("X-WORKS-Signature"), string(lineWorksBotSecretKey)) {
			lineWorksBotToken, err := h.aesCrypto.Decrypt(lineWorksBot.LineWorksBotToken[0], lineWorksBotIv.LineWorksBotTokenIV[0])
			if err != nil {
				slog.ErrorContext(ctx, "リクエストボディの復号化に失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			lineWorksBotRefreshToken, err := h.aesCrypto.Decrypt(lineWorksBot.LineWorksRefreshToken[0], lineWorksBotIv.LineWorksRefreshTokenIV[0])
			if err != nil {
				slog.ErrorContext(ctx, "リクエストボディの復号化に失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			lineWorksBotGroupID, err := h.aesCrypto.Decrypt(lineWorksBot.LineWorksGroupID[0], lineWorksBotIv.LineWorksGroupIDIV[0])
			if err != nil {
				slog.ErrorContext(ctx, "リクエストボディの復号化に失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			lineWorksBotID, err := h.aesCrypto.Decrypt(lineWorksBot.LineWorksBotID[0], lineWorksBotIv.LineWorksBotIDIV[0])
			if err != nil {
				slog.ErrorContext(ctx, "リクエストボディの復号化に失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			lineWorksBotDecrypt = &internal.LineWorksBotDecrypt{
				LineWorksBotToken:     string(lineWorksBotToken),
				LineWorksRefreshToken: string(lineWorksBotRefreshToken),
				LineWorksGroupID:      string(lineWorksBotGroupID),
				LineWorksBotID:        string(lineWorksBotID),
				RefreshTokenExpiresAt: lineWorksBot.RefreshTokenExpiresAt.Time,
				DefaultChannelID:      lineWorksBot.DefaultChannelID,
			}
			break
		}
		if i == len(lineWorksBots)-1 {
			slog.ErrorContext(ctx, "line_works_botの情報が見つかりませんでした。")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	// リクエストボディのバイトから構造体への変換
	err = json.Unmarshal(requestBodyBytes, &lineWorksResponses)
	if err != nil {
		slog.ErrorContext(ctx, "jsonの読み込みに失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// バリデーションチェック
	//err = lineWorksResponses.Validate()
	//if err != nil {
	//	slog.ErrorContext(ctx, "バリデーションチェックに失敗しました。", "エラー:", err.Error())
	//	http.Error(w, "Bad Request", http.StatusBadRequest)
	//	return
	//}

	lineWorks := lineworks.NewLineWorks(
		*h.indexService.Client,
		lineWorksBotDecrypt.LineWorksBotToken,
		lineWorksBotDecrypt.LineWorksRefreshToken,
		lineWorksBotDecrypt.LineWorksGroupID,
		lineWorksBotDecrypt.LineWorksBotID,
	)

	// ユーザー情報の取得
	lineWorksProfile, err := lineWorks.GetUserProfile(ctx, lineWorksResponses.Source.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "LINEユーザー情報の取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.InfoContext(ctx, "%v", lineWorksProfile)

	switch lineWorksResponses.Content.Type {
	case "text":
		_, err = h.indexService.DiscordSession.ChannelMessageSend(
			lineWorksBotDecrypt.DefaultChannelID,
			lineWorksProfile.UserName.FirstName+lineWorksProfile.UserName.LastName+": "+lineWorksResponses.Content.Text,
		)
		if err != nil {
			slog.ErrorContext(ctx, "Discordメッセージの送信に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case "image":
	}
	// レスポンスの書き込み
	w.WriteHeader(http.StatusOK)
}
