package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log/slog"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
)

func LineHmac(
	privateKey string,
	requestBodyByte []byte,
	aesCrypto crypto.AESInterface,
	lineBot *repository.LineBot,
	lineBotIv repository.LineBotIvNotClient,
	header string,
) (decrypt *LineBotDecrypt, err error) {
	lineBotDecrypt := &LineBotDecrypt{}

	lineBotSecretKey, err := aesCrypto.Decrypt(lineBot.LineBotSecret[0], lineBotIv.LineBotSecretIv[0])
	if err != nil {
		return nil, err
	}

	// macの生成
	mac := hmac.New(sha256.New, []byte(lineBotSecretKey))
	mac.Write(requestBodyByte)
	validSignByte := mac.Sum(nil)

	signature := base64.StdEncoding.EncodeToString(validSignByte)

	// 署名が一致しない場合は両方nilを返す
	if header != signature {
		slog.Error("signature is not match", "header", header, "signature", signature)
		return nil, nil
	}
	lineNotifyTokenByte, err := aesCrypto.Decrypt(lineBot.LineNotifyToken[0], lineBotIv.LineNotifyTokenIv[0])
	if err != nil {
		return nil, err
	}
	lineBotTokenByte, err := aesCrypto.Decrypt(lineBot.LineBotToken[0], lineBotIv.LineBotTokenIv[0])
	if err != nil {
		return nil, err
	}
	lineGroupByte, err := aesCrypto.Decrypt(lineBot.LineGroupID[0], lineBotIv.LineGroupIDIv[0])
	if err != nil {
		return nil, err
	}
	lineBotDecrypt.LineNotifyToken = string(lineNotifyTokenByte)
	lineBotDecrypt.LineBotToken = string(lineBotTokenByte)
	lineBotDecrypt.LineGroupID = string(lineGroupByte)
	lineBotDecrypt.DefaultChannelID = lineBot.DefaultChannelID
	lineBotDecrypt.DebugMode = lineBot.DebugMode
	return lineBotDecrypt, nil
}
