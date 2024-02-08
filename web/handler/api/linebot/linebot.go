package linebot

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linebot/internal"
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

	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}

	for i, lineBot := range lineBots {
		// 暗号化されたシークレットキーの復号化
		sercretKey, err := crypto.Decrypt(lineBot.LineBotSecret, keyBytes, lineBot.Iv)
		if err != nil {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
		inputSign := r.Header.Get("X-Line-Signature")
		// 受け取った署名をStringからByteへ変換
		inputSignByte, err := hex.DecodeString(inputSign)
		if err != nil {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}

		// macの生成
		mac := hmac.New(sha256.New, []byte(sercretKey))
		mac.Write(requestBodyByte)
		validSignByte := mac.Sum(nil)

		if hmac.Equal(inputSignByte, validSignByte) {
			break
		}
		// 最後の要素までループしても一致しなかった場合終了
		if i == len(lineBots)-1 {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
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
	// メッセージの種類によって処理を分岐
	if lineResponses.Events[0].Type == "text" {
	}
	// レスポンスの書き込み
	w.WriteHeader(http.StatusOK)
}
