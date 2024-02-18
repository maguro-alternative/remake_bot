package linetoken

import (
	"net/http"
	"encoding/json"
	"encoding/hex"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linetoken/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

type LineTokenHandler struct {
	IndexService *service.IndexService
}

func NewLineTokenHandler(indexService *service.IndexService) *LineTokenHandler {
	return &LineTokenHandler{
		IndexService: indexService,
	}
}

func (h *LineTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var lineTokenJson internal.LineBotJson
	if err := json.NewDecoder(r.Body).Decode(&lineTokenJson); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := lineTokenJson.Validate(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	if ctx == nil {
		ctx = r.Context()
	}
	// 暗号化キーの取得
	privateKey := config.PrivateKey()
	lineBot, lineBotIv, err := lineBotJsonEncrypt(privateKey, &lineTokenJson)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 暗号化
	repo := internal.NewRepository(h.IndexService.DB)
	if err := repo.UpdateLineBot(ctx, lineBot); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := repo.UpdateLineBotIv(ctx, lineBotIv); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func lineBotJsonEncrypt(privateKey string, lineBotJson *internal.LineBotJson) (Bot *internal.LineBot, BotIv *internal.LineBotIv, err error) {
	var lineBot *internal.LineBot
	var lineBotIv *internal.LineBotIv
	key, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, nil, err
	}
	// 暗号化
	if len(lineBot.LineNotifyToken) > 0 {
		if lineBot.LineNotifyToken, lineBotIv.LineNotifyTokenIv, err = crypto.Encrypt([]byte(lineBotJson.LineNotifyToken), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBot.LineBotToken) > 0 {
		if lineBot.LineBotToken, lineBotIv.LineBotTokenIv, err = crypto.Encrypt([]byte(lineBotJson.LineBotToken), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBot.LineBotSecret) > 0 {
		if lineBot.LineBotSecret, lineBotIv.LineBotSecretIv, err = crypto.Encrypt([]byte(lineBotJson.LineBotSecret), key) ; err != nil {
			return nil, nil, err
		}
	}
	if len(lineBot.LineGroupID) > 0 {
		if lineBot.LineGroupID, lineBotIv.LineGroupIDIv, err = crypto.Encrypt([]byte(lineBotJson.LineGroupID), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBot.LineClientID) > 0 {
		if lineBot.LineClientID, lineBotIv.LineClientIDIv, err = crypto.Encrypt([]byte(lineBotJson.LineClientID), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBot.LineClientSecret) > 0 {
		if lineBot.LineClientSecret, lineBotIv.LineClientSecretIv, err = crypto.Encrypt([]byte(lineBotJson.LineClientSecret), key); err != nil {
			return nil, nil, err
		}
	}
	return lineBot, lineBotIv, nil
}
