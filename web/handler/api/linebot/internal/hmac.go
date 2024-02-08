package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
)

func LineHmac(privateKey string, requestBodyByte []byte, lineBots []*LineBot, header string) (LineBotDecrypt, error) {
	var lineBotDecrypt LineBotDecrypt
	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return lineBotDecrypt, err
	}

	for _, lineBot := range lineBots {
		// 暗号化されたシークレットキーの復号化
		sercretKey, err := crypto.Decrypt(lineBot.LineBotSecret, keyBytes, lineBot.Iv)
		if err != nil {
			return lineBotDecrypt, err
		}

		// macの生成
		mac := hmac.New(sha256.New, []byte(sercretKey))
		mac.Write(requestBodyByte)
		validSignByte := mac.Sum(nil)

		signature := base64.StdEncoding.EncodeToString(validSignByte)

		if header != signature {
			continue
		}
		lineNotifyTokenByte, err := crypto.Decrypt(lineBot.LineNotifyToken, keyBytes, lineBot.Iv)
		if err != nil {
			return lineBotDecrypt, err
		}
		lineBotTokenByte, err := crypto.Decrypt(lineBot.LineBotToken, keyBytes, lineBot.Iv)
		if err != nil {
			return lineBotDecrypt, err
		}
		lineGroupByte, err := crypto.Decrypt(lineBot.LineGroupID, keyBytes, lineBot.Iv)
		if err != nil {
			return lineBotDecrypt, err
		}
		lineClientIDByte, err := crypto.Decrypt(lineBot.LineClientID, keyBytes, lineBot.Iv)
		if err != nil {
			return lineBotDecrypt, err
		}
		lineClientSercretByte, err := crypto.Decrypt(lineBot.LineClientSercret, keyBytes, lineBot.Iv)
		if err != nil {
			return lineBotDecrypt, err
		}
		lineBotDecrypt.LineNotifyToken = string(lineNotifyTokenByte)
		lineBotDecrypt.LineBotToken = string(lineBotTokenByte)
		lineBotDecrypt.LineGroupID = string(lineGroupByte)
		lineBotDecrypt.LineClientID = string(lineClientIDByte)
		lineBotDecrypt.LineClientSercret = string(lineClientSercretByte)
		lineBotDecrypt.DefaultChannelID = lineBot.DefaultChannelID
		lineBotDecrypt.DebugMode = lineBot.DebugMode
		return lineBotDecrypt, nil
	}
	return lineBotDecrypt, nil
}

