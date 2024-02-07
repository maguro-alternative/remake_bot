package api

import (
	"net/http"
	"context"
	"io"
	"log"
	"crypto/hmac"
    "crypto/sha256"
	"encoding/hex"

	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/pkg/crypto"
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
}

// ServeHTTP handles HTTP requests.
func (h *LineBotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var iv []byte
	privateKey := config.PrivateKey()
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	requestBodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to Load Request")
		//response.BadRequest(writer, "Failed to Load Request")
		return
	}

	keyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		log.Println("Failed to Load Request")
		return
	}

	var lineClientSecrets []LineSecret
	err = h.IndexService.DB.SelectContext(ctx, &lineClientSecrets, "TestSercretKey")
	if err != nil {
		log.Println("Failed to Load Request")
		return
	}

	for _, lineClientSecret := range lineClientSecrets {
		err = h.IndexService.DB.GetContext(ctx, &iv, "TestSercret")
		if err != nil {
			log.Println("Failed to Load Request")
			return
		}

		sercretKey, err := crypto.Decrypt(lineClientSecret.ClientSercret, keyBytes, iv)
		if err != nil {
			log.Println("Failed to Load Request")
			return
		}
		inputSign := r.Header.Get("X-Line-Signature")
		// 受け取った署名をStringからByteへ変換
		inputSignByte, err := hex.DecodeString(inputSign)
		if err != nil {
			return
		}

		// macの生成
		mac := hmac.New(sha256.New, []byte(sercretKey))
		mac.Write(requestBodyByte)
		validSignByte := mac.Sum(nil)

		if hmac.Equal(inputSignByte, validSignByte) {
			break
		}
	}
}