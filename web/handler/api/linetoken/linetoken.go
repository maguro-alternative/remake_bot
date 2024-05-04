package linetoken

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/lib/pq"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linetoken/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

type LineTokenHandler struct {
	indexService *service.IndexService
	repo         repository.RepositoryFunc
}

func NewLineTokenHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *LineTokenHandler {
	return &LineTokenHandler{
		indexService: indexService,
		repo:         repo,
	}
}

func (h *LineTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "/api/linetoken Method Not Allowed")
		return
	}
	var lineTokenJson internal.LineBotJson

	if err := json.NewDecoder(r.Body).Decode(&lineTokenJson); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonの読み取りに失敗しました:"+err.Error())
		return
	}
	if err := lineTokenJson.Validate(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonのバリデーションに失敗しました:"+err.Error())
		return
	}

	// 暗号化キーの取得
	privateKey := config.PrivateKey()
	lineBot, lineBotIv, err := lineBotJsonEncrypt(privateKey, &lineTokenJson)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "暗号化に失敗しました:"+err.Error())
		return
	}

	if err := h.repo.UpdateLineBot(ctx, lineBot); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_botの更新に失敗しました:"+err.Error())
		return
	}
	if err := h.repo.UpdateLineBotIv(ctx, lineBotIv); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_bot_ivの更新に失敗しました:"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}

func lineBotJsonEncrypt(privateKey string, lineBotJson *internal.LineBotJson) (Bot *repository.LineBot, BotIv *repository.LineBotIv, err error) {
	lineBot := repository.LineBot{
		LineNotifyToken:  make(pq.ByteaArray, 1),
		LineBotToken:     make(pq.ByteaArray, 1),
		LineBotSecret:    make(pq.ByteaArray, 1),
		LineGroupID:      make(pq.ByteaArray, 1),
		LineClientID:     make(pq.ByteaArray, 1),
		LineClientSecret: make(pq.ByteaArray, 1),
	}
	lineBotIv := repository.LineBotIv{
		LineNotifyTokenIv:  make(pq.ByteaArray, 1),
		LineBotTokenIv:     make(pq.ByteaArray, 1),
		LineBotSecretIv:    make(pq.ByteaArray, 1),
		LineGroupIDIv:      make(pq.ByteaArray, 1),
		LineClientIDIv:     make(pq.ByteaArray, 1),
		LineClientSecretIv: make(pq.ByteaArray, 1),
	}
	key, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, nil, err
	}
	// 暗号化
	if len(lineBotJson.LineNotifyToken) > 0 {
		if lineBotIv.LineNotifyTokenIv[0], lineBot.LineNotifyToken[0], err = crypto.Encrypt([]byte(lineBotJson.LineNotifyToken), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineBotToken) > 0 {
		if lineBotIv.LineBotTokenIv[0], lineBot.LineBotToken[0], err = crypto.Encrypt([]byte(lineBotJson.LineBotToken), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineBotSecret) > 0 {
		if lineBotIv.LineBotSecretIv[0], lineBot.LineBotSecret[0], err = crypto.Encrypt([]byte(lineBotJson.LineBotSecret), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineGroupID) > 0 {
		if lineBotIv.LineGroupIDIv[0], lineBot.LineGroupID[0], err = crypto.Encrypt([]byte(lineBotJson.LineGroupID), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineClientID) > 0 {
		if lineBotIv.LineClientIDIv[0], lineBot.LineClientID[0], err = crypto.Encrypt([]byte(lineBotJson.LineClientID), key); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineClientSecret) > 0 {
		if lineBotIv.LineClientSecretIv[0], lineBot.LineClientSecret[0], err = crypto.Encrypt([]byte(lineBotJson.LineClientSecret), key); err != nil {
			return nil, nil, err
		}
	}
	lineBotIv.GuildID = lineBotJson.GuildID
	lineBot.GuildID = lineBotJson.GuildID
	lineBot.DefaultChannelID = lineBotJson.DefaultChannelID
	lineBot.DebugMode = lineBotJson.DebugMode
	return &lineBot, &lineBotIv, nil
}
