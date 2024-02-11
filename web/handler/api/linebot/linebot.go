package linebot

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/line"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linebot/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

type Repository interface {
	GetLineBots(ctx context.Context) ([]*internal.LineBot, error)
}

// A LineBotHandler handles requests for the line bot.
type LineBotHandler struct {
	IndexService *service.IndexService
}

// NewLineBotHandler returns new LineBotHandler.
func NewLineBotHandler(indexService *service.IndexService) *LineBotHandler {
	return &LineBotHandler{
		IndexService: indexService,
	}
}

// ServeHTTP handles HTTP requests.
func (h *LineBotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var lineResponses LineResponses
	// 暗号化キーの取得
	privateKey := config.PrivateKey()
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	repo := internal.NewRepository(h.IndexService.DB)
	lineBots, err := repo.GetLineBots(ctx)
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}

	// リクエストボディの読み込み
	requestBodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}

	// リクエストボディの検証
	lineBotDecrypt, err := internal.LineHmac(privateKey, requestBodyByte, lineBots, r.Header.Get("X-Line-Signature"))
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}

	// リクエストボディのバイトから構造体への変換
	err = json.Unmarshal(requestBodyByte, &lineResponses)
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}
	// バリデーションチェック
	if err := lineResponses.Validate(); err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}


	lineRequ := line.NewLineRequest(lineBotDecrypt.LineNotifyToken, lineBotDecrypt.LineBotToken, lineBotDecrypt.LineGroupID)
	// ユーザー情報の取得
	lineProfile, err := lineRequ.GetProfile(ctx, lineResponses.Events[0].Source.UserID)
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}

	// メッセージの種類によって処理を分岐
	if lineResponses.Events[0].Type == "text" {
		_, err = h.IndexService.DiscordSession.ChannelMessageSend(
			lineBotDecrypt.DefaultChannelID,
			lineProfile.DisplayName+"\n「 "+lineResponses.Events[0].Message.Text+" 」",
		)
		if err != nil {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
	}
	// レスポンスの書き込み
	w.WriteHeader(http.StatusOK)
}
