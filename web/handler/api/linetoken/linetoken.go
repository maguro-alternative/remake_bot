package linetoken

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/web/config"
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
		json.NewEncoder(w).Encode("Method Not Allowed")
		return
	}
	var lineTokenJson internal.LineBotJson
	if err := json.NewDecoder(r.Body).Decode(&lineTokenJson); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	if err := lineTokenJson.Validate(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	ctx := r.Context()
	if ctx == nil {
		ctx = r.Context()
	}
	lineTokenJson.GuildID = r.PathValue("guildId")
	// 暗号化キーの取得
	privateKey := config.PrivateKey()
	fmt.Println(&lineTokenJson)
	lineBot, lineBotIv, err := lineBotJsonEncrypt(privateKey, &lineTokenJson)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	// 暗号化
	repo := internal.NewRepository(h.IndexService.DB)
	if err := repo.UpdateLineBot(ctx, lineBot); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	if err := repo.UpdateLineBotIv(ctx, lineBotIv); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}

func lineBotJsonEncrypt(privateKey string, lineBotJson *internal.LineBotJson) (Bot *internal.LineBot, BotIv *internal.LineBotIv, err error) {
	var lineBot internal.LineBot
	var lineBotIv internal.LineBotIv
	fmt.Println(lineBotJson.GuildID)
	key, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, nil, err
	}
	// 暗号化
	if len(lineBotJson.LineNotifyToken) > 0 {
		if lineBot.LineNotifyToken[0], lineBotIv.LineNotifyTokenIv[0], err = crypto.Encrypt([]byte(lineBotJson.LineNotifyToken), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineBotToken) > 0 {
		if lineBot.LineBotToken[0], lineBotIv.LineBotTokenIv[0], err = crypto.Encrypt([]byte(lineBotJson.LineBotToken), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineBotSecret) > 0 {
		if lineBot.LineBotSecret[0], lineBotIv.LineBotSecretIv[0], err = crypto.Encrypt([]byte(lineBotJson.LineBotSecret), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineGroupID) > 0 {
		if lineBot.LineGroupID[0], lineBotIv.LineGroupIDIv[0], err = crypto.Encrypt([]byte(lineBotJson.LineGroupID), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineClientID) > 0 {
		if lineBot.LineClientID[0], lineBotIv.LineClientIDIv[0], err = crypto.Encrypt([]byte(lineBotJson.LineClientID), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineClientSecret) > 0 {
		if lineBot.LineClientSecret[0], lineBotIv.LineClientSecretIv[0], err = crypto.Encrypt([]byte(lineBotJson.LineClientSecret), key); err != nil {
			return nil, nil, err
		}
	}
	lineBot.GuildID = lineBotJson.GuildID
	lineBot.DefaultChannelID = lineBotJson.DefaultChannelID
	lineBot.DebugMode = lineBotJson.DebugMode
	return &lineBot, &lineBotIv, nil
}
