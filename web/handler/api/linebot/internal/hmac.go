package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
)

func LineHmac(privateKey string, requestBodyByte []byte, lineBot *repository.LineBot, lineBotIv repository.LineBotIvNotClient, header string) (decrypt *LineBotDecrypt, err error) {
	lineBotDecrypt := &LineBotDecrypt{}
	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	lineBotSecretKey, err := crypto.Decrypt(lineBot.LineBotSecret[0], keyBytes, lineBotIv.LineBotSecretIv[0])
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
		return nil, nil
	}
	lineNotifyTokenByte, err := crypto.Decrypt(lineBot.LineNotifyToken[0], keyBytes, lineBotIv.LineNotifyTokenIv[0])
	if err != nil {
		return nil, err
	}
	lineBotTokenByte, err := crypto.Decrypt(lineBot.LineBotToken[0], keyBytes, lineBotIv.LineBotTokenIv[0])
	if err != nil {
		return nil, err
	}
	lineGroupByte, err := crypto.Decrypt(lineBot.LineGroupID[0], keyBytes, lineBotIv.LineGroupIDIv[0])
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

