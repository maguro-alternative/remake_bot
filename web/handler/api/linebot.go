package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
)

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

type LineSecret struct {
	ClientSercret []byte `db:"line_client_sercret"`
	Iv            []byte `db:"iv"`
}

// ServeHTTP handles HTTP requests.
func (h *LineBotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// 暗号化キーの取得
	privateKey := config.PrivateKey()
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
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

	var lineClientSecrets []LineSecret
	query := `
		SELECT
			line_client_sercret,
			iv
		FROM
			line_bot
		WHERE
			line_client_sercret IS NOT NULL
		AND
			iv IS NOT NULL
	`
	err = h.IndexService.DB.SelectContext(ctx, &lineClientSecrets, query)
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}

	for i, lineClientSecret := range lineClientSecrets {
		// 暗号化されたシークレットキーの復号化
		sercretKey, err := crypto.Decrypt(lineClientSecret.ClientSercret, keyBytes, lineClientSecret.Iv)
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
		if i == len(lineClientSecrets)-1 {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
	}
}
