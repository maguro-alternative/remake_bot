package linetoken

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/lib/pq"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/line"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/linetoken/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

type LineTokenHandler struct {
	indexService *service.IndexService
	repo         repository.RepositoryFunc
	aesCrypto    crypto.AESInterface
}

func NewLineTokenHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
	aesCrypto crypto.AESInterface,
) *LineTokenHandler {
	return &LineTokenHandler{
		indexService: indexService,
		repo:         repo,
		aesCrypto:    aesCrypto,
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
		http.Error(w, "Bad Request", http.StatusUnprocessableEntity)
		slog.ErrorContext(ctx, "jsonのバリデーションに失敗しました:"+err.Error())
		return
	}

	guildId := r.PathValue("guildId")
	if lineTokenJson.GuildID == "" {
		lineTokenJson.GuildID = guildId
	}

	if err := verifyLineToken(ctx, h.repo, h.aesCrypto, h.indexService.Client, &lineTokenJson); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "トークンの検証に失敗しました:"+err.Error())
		return
	}

	// 暗号化キーの取得
	lineBot, lineBotIv, err := lineBotJsonEncrypt(h.aesCrypto, &lineTokenJson)
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

func lineBotJsonEncrypt(aesCrypto crypto.AESInterface, lineBotJson *internal.LineBotJson) (Bot *repository.LineBot, BotIv *repository.LineBotIv, err error) {
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
	// 暗号化
	if len(lineBotJson.LineNotifyToken) > 0 {
		if lineBotIv.LineNotifyTokenIv[0], lineBot.LineNotifyToken[0], err = aesCrypto.Encrypt([]byte(lineBotJson.LineNotifyToken)); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineBotToken) > 0 {
		if lineBotIv.LineBotTokenIv[0], lineBot.LineBotToken[0], err = aesCrypto.Encrypt([]byte(lineBotJson.LineBotToken)); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineBotSecret) > 0 {
		if lineBotIv.LineBotSecretIv[0], lineBot.LineBotSecret[0], err = aesCrypto.Encrypt([]byte(lineBotJson.LineBotSecret)); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineGroupID) > 0 {
		if lineBotIv.LineGroupIDIv[0], lineBot.LineGroupID[0], err = aesCrypto.Encrypt([]byte(lineBotJson.LineGroupID)); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineClientID) > 0 {
		if lineBotIv.LineClientIDIv[0], lineBot.LineClientID[0], err = aesCrypto.Encrypt([]byte(lineBotJson.LineClientID)); err != nil {
			return nil, nil, err
		}
	}
	if len(lineBotJson.LineClientSecret) > 0 {
		if lineBotIv.LineClientSecretIv[0], lineBot.LineClientSecret[0], err = aesCrypto.Encrypt([]byte(lineBotJson.LineClientSecret)); err != nil {
			return nil, nil, err
		}
	}
	lineBotIv.GuildID = lineBotJson.GuildID
	lineBot.GuildID = lineBotJson.GuildID
	lineBot.DefaultChannelID = lineBotJson.DefaultChannelID
	lineBot.DebugMode = lineBotJson.DebugMode
	return &lineBot, &lineBotIv, nil
}

func verifyLineToken(
	ctx context.Context,
	repo repository.RepositoryFunc,
	aesCrypto crypto.AESInterface,
	client *http.Client,
	lineTokenJson *internal.LineBotJson,
) error {
	lineNotifyToken := lineTokenJson.LineNotifyToken
	lineBotToken := lineTokenJson.LineBotToken
	lineGroupID := lineTokenJson.LineGroupID
	linebot, err := repo.GetAllColumnsLineBot(ctx, lineTokenJson.GuildID)
	if err != nil {
		return err
	}
	linebotIv, err := repo.GetAllColumnsLineBotIv(ctx, lineTokenJson.GuildID)
	if err != nil {
		return err
	}
	if lineNotifyToken == "" && len(linebot.LineNotifyToken) > 0 && len(linebotIv.LineNotifyTokenIv) > 0 {
		notifyByte, err := aesCrypto.Decrypt(linebot.LineNotifyToken[0], linebotIv.LineNotifyTokenIv[0])
		if err != nil {
			return err
		}
		lineTokenJson.LineNotifyToken = string(notifyByte)
	}
	if lineBotToken == "" && len(linebot.LineBotToken) > 0 && len(linebotIv.LineBotTokenIv) > 0 {
		botByte, err := aesCrypto.Decrypt(linebot.LineBotToken[0], linebotIv.LineBotTokenIv[0])
		if err != nil {
			return err
		}
		lineTokenJson.LineBotToken = string(botByte)
	}
	if lineGroupID == "" && len(linebot.LineGroupID) > 0 && len(linebotIv.LineGroupIDIv) > 0 {
		groupByte, err := aesCrypto.Decrypt(linebot.LineGroupID[0], linebotIv.LineGroupIDIv[0])
		if err != nil {
			return err
		}
		lineTokenJson.LineGroupID = string(groupByte)
	}
	lineRequ := line.NewLineRequest(
		*client,
		lineTokenJson.LineNotifyToken,
		lineTokenJson.LineBotToken,
		lineTokenJson.LineGroupID,
	)
	// Line Notify Tokenの検証
	if lineTokenJson.LineNotifyToken != "" {
		err := lineRequ.PushMessageNotify(ctx, "トークン検証のテストメッセージです。")
		if err != nil {
			return err
		}
	}
	// Line Bot Tokenの検証
	if lineTokenJson.LineBotToken != "" {
		botInfo, err := lineRequ.GetBotInfo(ctx)
		if err != nil {
			return err
		}
		if botInfo.Message != "" {
			return errors.New(botInfo.Message)
		}
	}
	// Line Group IDの検証
	if lineTokenJson.LineGroupID == "" {
		return nil
	}
	_, err = lineRequ.GetGroupUserCount(ctx)
	return err
}
