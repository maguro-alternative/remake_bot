package api

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

type LineBot struct {
	LineNotifyToken   []byte `db:"line_notify_token"`
	LineBotToken      []byte `db:"line_bot_token"`
	LineBotSecret     []byte `db:"line_bot_secret"`
	LineGroupID       []byte `db:"line_group_id"`
	LineClientID      []byte `db:"line_client_id"`
	LineClientSercret []byte `db:"line_client_sercret"`
	Iv                []byte `db:"iv"`
	DefaultChannelID  string `db:"default_channel_id"`
	DebugMode         bool   `db:"debug_mode"`
}

type LineResponses struct {
	Events []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Source     struct {
			GroupID string `json:"groupId"`
			UserID  string `json:"userId"`
			Type    string `json:"type"`
		} `json:"source"`
		Timestamp int64 `json:"timestamp"`
		Message   struct {
			ID        string `json:"id"`
			Type      string `json:"type"`
			Text      string `json:"text"`
			ReplyToken string `json:"replyToken"`
		} `json:"message"`
	} `json:"events"`
}

// ServeHTTP handles HTTP requests.
func (h *LineBotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var lineBots []LineBot
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
	query := `
		SELECT
			line_notify_token,
			line_bot_token,
			line_bot_secret,
			line_group_id,
			line_client_id,
			line_client_sercret,
			iv,
			default_channel_id,
			debug_mode
		FROM
			line_bot
		WHERE
			line_notify_token IS NOT NULL
		AND
			line_bot_token IS NOT NULL
		AND
			line_bot_secret IS NOT NULL
		AND
			line_group_id IS NOT NULL
		AND
			line_client_id IS NOT NULL
		AND
			line_client_sercret IS NOT NULL
		AND
			iv IS NOT NULL
	`
	err = h.IndexService.DB.SelectContext(ctx, &lineBots, query)
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
	err = json.Unmarshal(requestBodyByte, &lineBots)
}
