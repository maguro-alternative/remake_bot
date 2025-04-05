package lineworkstoken

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/lib/pq"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/lineworks"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/lineworks_token/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

type LineWorksTokenHandler struct {
	indexService *service.IndexService
	repo         repository.RepositoryFunc
	aesCrypto    crypto.AESInterface
}

func NewLineWorksTokenHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
	aesCrypto crypto.AESInterface,
) *LineWorksTokenHandler {
	return &LineWorksTokenHandler{
		indexService: indexService,
		repo:         repo,
		aesCrypto:    aesCrypto,
	}
}

func (h *LineWorksTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "/api/lineworks-token Method Not Allowed")
		return
	}

	var lineWorksTokenJson internal.LineWorksTokenJson

	if err := json.NewDecoder(r.Body).Decode(&lineWorksTokenJson); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonの読み取りに失敗しました:"+err.Error())
		return
	}

	guildId := r.PathValue("guildId")
	if guildId == "" {
		guildId = lineWorksTokenJson.GuildID
	}
	if lineWorksTokenJson.GuildID == "" {
		lineWorksTokenJson.GuildID = guildId
	}

	worksBot, err := h.repo.GetLineWorksBotByGuildID(ctx, guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_works_botの取得に失敗しました:"+err.Error())
		return
	}
	worksBotIv, err := h.repo.GetLineWorksBotIVByGuildID(ctx, guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_works_bot_ivの取得に失敗しました:"+err.Error())
		return
	}
	worksBotInfo, err := h.repo.GetLineWorksBotInfoByGuildID(ctx, guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_works_bot_infoの取得に失敗しました:"+err.Error())
		return
	}
	worksBotInfoIv, err := h.repo.GetLineWorksBotInfoIVByGuildID(ctx, guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_works_bot_info_ivの取得に失敗しました:"+err.Error())
		return
	}

	lineWorksInfo := newLineWorksInfoGenerate(
		h.indexService.Client,
		h.aesCrypto,
		&lineWorksTokenJson,
		worksBot,
		worksBotIv,
		worksBotInfo,
		worksBotInfoIv,
	)

	lineWorksTokenResponse, err := lineWorksInfo.GetAccessToken(
		ctx,
		"bot.message bot.read user.profile.read",
	)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "トークンの取得に失敗しました:"+err.Error())
		return
	}

	expiresIn, err := strconv.Atoi(lineWorksTokenResponse.ExpiresIn)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "ExpiresInの変換に失敗しました:"+err.Error())
		return
	}

	lineWorksToken, lineWorksTokenIv, lineWorksTokenInfo, lineWorksTokenInfoIv, err := lineWorksTokenJsonEncrypt(
		h.aesCrypto,
		lineWorksTokenResponse.AccessToken,
		lineWorksTokenResponse.RefreshToken,
		time.Now().Add(time.Duration(expiresIn) * time.Second),
		&lineWorksTokenJson,
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "暗号化に失敗しました:"+err.Error())
		return
	}

	lineWorks := repository.NewLineWorksBot(
		lineWorksToken.GuildID,
		lineWorksToken.LineWorksBotToken,
		lineWorksToken.LineWorksRefreshToken,
		lineWorksToken.LineWorksGroupID,
		lineWorksToken.LineWorksBotID,
		lineWorksToken.LineWorksBotSecret,
		pq.NullTime{Time: lineWorksToken.RefreshTokenExpiresAt, Valid: true},
		lineWorksTokenJson.DefaultChannelID,
		lineWorksTokenJson.DebugMode,
	)

	if err := h.repo.UpdateLineWorksBot(ctx, lineWorks); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_works_botの更新に失敗しました:"+err.Error())
		return
	}

	lineWorksIv := repository.NewLineWorksBotIV(
		lineWorksToken.GuildID,
		lineWorksTokenIv.LineWorksBotTokenIv,
		lineWorksTokenIv.LineWorksRefreshTokenIv,
		lineWorksTokenIv.LineWorksGroupIDIv,
		lineWorksTokenIv.LineWorksBotIDIv,
		lineWorksTokenIv.LineWorksBotSecretIv,
	)
	if err := h.repo.UpdateLineWorksBotIV(ctx, lineWorksIv); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_works_bot_ivの更新に失敗しました:"+err.Error())
		return
	}

	lineWorksBotInfo := repository.NewLineWorksBotInfo(
		lineWorksTokenInfo.GuildID,
		lineWorksTokenInfo.LineWorksClientID,
		lineWorksTokenInfo.LineWorksClientSecret,
		lineWorksTokenInfo.LineWorksServiceAccount,
		lineWorksTokenInfo.LineWorksPrivateKey,
		lineWorksTokenInfo.LineWorksDomainID,
		lineWorksTokenInfo.LineWorksAdminID,
	)

	if err := h.repo.UpdateLineWorksBotInfo(ctx, lineWorksBotInfo); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_works_bot_infoの更新に失敗しました:"+err.Error())
		return
	}

	lineWorksBotInfoIv := repository.NewLineWorksBotInfoIV(
		lineWorksTokenInfo.GuildID,
		lineWorksTokenInfoIv.LineWorksClientIDIv,
		lineWorksTokenInfoIv.LineWorksClientSecretIv,
		lineWorksTokenInfoIv.LineWorksServiceAccountIv,
		lineWorksTokenInfoIv.LineWorksPrivateKeyIv,
		lineWorksTokenInfoIv.LineWorksDomainIDIv,
		lineWorksTokenInfoIv.LineWorksAdminIDIv,
	)
	if err := h.repo.UpdateLineWorksBotInfoIV(ctx, lineWorksBotInfoIv); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_works_bot_info_ivの更新に失敗しました:"+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func lineWorksTokenJsonEncrypt(
	aesCrypto crypto.AESInterface,
	lineWorksBotToken string,
	lineWorksRefreshToken string,
	refreshTokenExpiresAt time.Time,
	lineWorksTokenJson *internal.LineWorksTokenJson,
) (
	lineWorksToken *internal.LineWorksBot,
	lineWorksTokenIv *internal.LineWorksBotIv,
	lineWorksTokenInfo *internal.LineWorksBotInfo,
	lineWorksTokenInfoIv *internal.LineWorksBotInfoIv,
	err error,
) {
	lineWorksToken = &internal.LineWorksBot{
		GuildID:          lineWorksTokenJson.GuildID,
		LineWorksBotToken: make(pq.ByteaArray, 1),
		LineWorksRefreshToken: make(pq.ByteaArray, 1),
		LineWorksGroupID: make(pq.ByteaArray, 1),
		LineWorksBotID: make(pq.ByteaArray, 1),
		LineWorksBotSecret: make(pq.ByteaArray, 1),
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}
	lineWorksTokenIv = &internal.LineWorksBotIv{
		GuildID:          lineWorksTokenJson.GuildID,
		LineWorksBotTokenIv: make(pq.ByteaArray, 1),
		LineWorksRefreshTokenIv: make(pq.ByteaArray, 1),
		LineWorksGroupIDIv: make(pq.ByteaArray, 1),
		LineWorksBotIDIv: make(pq.ByteaArray, 1),
		LineWorksBotSecretIv: make(pq.ByteaArray, 1),
	}
	if len(lineWorksBotToken) > 0 {
		lineWorksTokenIv.LineWorksBotTokenIv[0], lineWorksToken.LineWorksBotToken[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksBotToken)}[0])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if len(lineWorksRefreshToken) > 0 {
		lineWorksTokenIv.LineWorksRefreshTokenIv[0], lineWorksToken.LineWorksRefreshToken[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksRefreshToken)}[0])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if len(lineWorksTokenJson.LineWorksGroupID) > 0 {
		lineWorksTokenIv.LineWorksGroupIDIv[0], lineWorksToken.LineWorksGroupID[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksTokenJson.LineWorksGroupID)}[0])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if len(lineWorksTokenJson.LineWorksBotID) > 0 {
		lineWorksTokenIv.LineWorksBotIDIv[0], lineWorksToken.LineWorksBotID[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksTokenJson.LineWorksBotID)}[0])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if len(lineWorksTokenJson.LineWorksBotSecret) > 0 {
		lineWorksTokenIv.LineWorksBotSecretIv[0], lineWorksToken.LineWorksBotSecret[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksTokenJson.LineWorksBotSecret)}[0])
		if err != nil {
			return nil, nil,nil, nil,  err
		}
	}
	lineWorksTokenIv.GuildID = lineWorksTokenJson.GuildID
	lineWorksToken.DefaultChannelID = lineWorksTokenJson.DefaultChannelID

	lineWorksTokenInfo = &internal.LineWorksBotInfo{
		GuildID:               lineWorksTokenJson.GuildID,
		LineWorksClientID:     make(pq.ByteaArray, 1),
		LineWorksClientSecret: make(pq.ByteaArray, 1),
		LineWorksServiceAccount: make(pq.ByteaArray, 1),
		LineWorksPrivateKey: make(pq.ByteaArray, 1),
		LineWorksDomainID: make(pq.ByteaArray, 1),
		LineWorksAdminID: make(pq.ByteaArray, 1),
	}
	lineWorksTokenInfoIv = &internal.LineWorksBotInfoIv{
		GuildID:               lineWorksTokenJson.GuildID,
		LineWorksClientIDIv:     make(pq.ByteaArray, 1),
		LineWorksClientSecretIv: make(pq.ByteaArray, 1),
		LineWorksServiceAccountIv: make(pq.ByteaArray, 1),
		LineWorksPrivateKeyIv: make(pq.ByteaArray, 1),
		LineWorksDomainIDIv: make(pq.ByteaArray, 1),
		LineWorksAdminIDIv: make(pq.ByteaArray, 1),
	}
	if len(lineWorksTokenJson.LineWorksClientID) > 0 {
		lineWorksTokenInfoIv.LineWorksClientIDIv[0], lineWorksTokenInfo.LineWorksClientID[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksTokenJson.LineWorksClientID)}[0])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if len(lineWorksTokenJson.LineWorksClientSecret) > 0 {
		lineWorksTokenInfoIv.LineWorksClientSecretIv[0], lineWorksTokenInfo.LineWorksClientSecret[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksTokenJson.LineWorksClientSecret)}[0])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if len(lineWorksTokenJson.LineWorksServiceAccount) > 0 {
		lineWorksTokenInfoIv.LineWorksServiceAccountIv[0], lineWorksTokenInfo.LineWorksServiceAccount[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksTokenJson.LineWorksServiceAccount)}[0])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if len(lineWorksTokenJson.LineWorksPrivateKey) > 0 {
		lineWorksTokenInfoIv.LineWorksPrivateKeyIv[0], lineWorksTokenInfo.LineWorksPrivateKey[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksTokenJson.LineWorksPrivateKey)}[0])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if len(lineWorksTokenJson.LineWorksDomainID) > 0 {
		lineWorksTokenInfoIv.LineWorksDomainIDIv[0], lineWorksTokenInfo.LineWorksDomainID[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksTokenJson.LineWorksDomainID)}[0])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if len(lineWorksTokenJson.LineWorksAdminID) > 0 {
		lineWorksTokenInfoIv.LineWorksAdminIDIv[0], lineWorksTokenInfo.LineWorksAdminID[0], err = aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineWorksTokenJson.LineWorksAdminID)}[0])
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	lineWorksTokenInfoIv.GuildID = lineWorksTokenJson.GuildID

	return lineWorksToken, lineWorksTokenIv, lineWorksTokenInfo, lineWorksTokenInfoIv, nil
}

func newLineWorksInfoGenerate(
	client *http.Client,
	aesCrypto crypto.AESInterface,
	lineWorksTokenJson *internal.LineWorksTokenJson,
	worksBot *repository.LineWorksBot,
	worksBotIv *repository.LineWorksBotIV,
	worksBotInfo *repository.LineWorksBotInfo,
	worksBotInfoIv *repository.LineWorksBotInfoIV,
) *lineworks.LineWorksInfo {
	var clientIdByte, clientSecretByte, serviceAccountByte, privateKeyByte, domainIdByte, adminIdByte []byte
	var err error
	if len(worksBotInfo.LineWorksClientID) > 0 {
		clientIdByte, err = aesCrypto.Decrypt(worksBotInfo.LineWorksClientID[0], worksBotInfoIv.LineWorksClientIDIV[0])
		if err != nil {
			return nil
		}
	}
	if len(worksBotInfo.LineWorksClientSecret) > 0 {
		clientSecretByte, err = aesCrypto.Decrypt(worksBotInfo.LineWorksClientSecret[0], worksBotInfoIv.LineWorksClientSecretIV[0])
		if err != nil {
			return nil
		}
	}
	if len(worksBotInfo.LineWorksServiceAccount) > 0 {
		serviceAccountByte, err = aesCrypto.Decrypt(worksBotInfo.LineWorksServiceAccount[0], worksBotInfoIv.LineWorksServiceAccountIV[0])
		if err != nil {
			return nil
		}
	}
	if len(worksBotInfo.LineWorksPrivateKey) > 0 {
		privateKeyByte, err = aesCrypto.Decrypt(worksBotInfo.LineWorksPrivateKey[0], worksBotInfoIv.LineWorksPrivateKeyIV[0])
		if err != nil {
			return nil
		}
	}
	if len(worksBotInfo.LineWorksDomainID) > 0 {
		domainIdByte, err = aesCrypto.Decrypt(worksBotInfo.LineWorksDomainID[0], worksBotInfoIv.LineWorksDomainIDIV[0])
		if err != nil {
			return nil
		}
	}
	if len(worksBotInfo.LineWorksAdminID) > 0 {
		adminIdByte, err = aesCrypto.Decrypt(worksBotInfo.LineWorksAdminID[0], worksBotInfoIv.LineWorksAdminIDIV[0])
		if err != nil {
			return nil
		}
	}

	clientId := string(clientIdByte)
	clientSecret := string(clientSecretByte)
	serviceAccount := string(serviceAccountByte)
	privateKey := string(privateKeyByte)
	domainId := string(domainIdByte)
	adminId := string(adminIdByte)

	if lineWorksTokenJson.LineWorksClientID != "" {
		clientId = lineWorksTokenJson.LineWorksClientID
	}
	if lineWorksTokenJson.LineWorksClientSecret != "" {
		clientSecret = lineWorksTokenJson.LineWorksClientSecret
	}
	if lineWorksTokenJson.LineWorksServiceAccount != "" {
		serviceAccount = lineWorksTokenJson.LineWorksServiceAccount
	}
	if lineWorksTokenJson.LineWorksPrivateKey != "" {
		privateKey = lineWorksTokenJson.LineWorksPrivateKey
	}
	if lineWorksTokenJson.LineWorksDomainID != "" {
		domainId = lineWorksTokenJson.LineWorksDomainID
	}
	if lineWorksTokenJson.LineWorksAdminID != "" {
		adminId = lineWorksTokenJson.LineWorksAdminID
	}

	return lineworks.NewLineWorksInfo(
		*client,
		clientId,
		clientSecret,
		serviceAccount,
		privateKey,
		domainId,
		adminId,
	)
}
