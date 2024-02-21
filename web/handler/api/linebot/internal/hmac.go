package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
)

func LineHmac(privateKey string, requestBodyByte []byte, lineBot LineBot, lineBotIv LineBotIv, header string) (decrypt *LineBotDecrypt, err error) {
	var lineBotDecrypt *LineBotDecrypt
	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	decodeBotSecret, err := hex.DecodeString(string(lineBotIv.LineBotSecretIv[0]))
	if err != nil {
		return nil, err
	}
	decodeNotifyToken, err := hex.DecodeString(string(lineBotIv.LineNotifyTokenIv[0]))
	if err != nil {
		return
	}
	decodeBotToken, err := hex.DecodeString(string(lineBotIv.LineBotTokenIv[0]))
	if err != nil {
		return
	}
	decodeGroupID, err := hex.DecodeString(string(lineBotIv.LineGroupIDIv[0]))
	if err != nil {
		return
	}
	lineBotSecretStr, err := base64.StdEncoding.DecodeString(string(lineBot.LineBotSecret[0]))
	if err != nil {
		return nil, err
	}
	lineNotifyStr, err := base64.StdEncoding.DecodeString(string(lineBot.LineNotifyToken[0]))
	if err != nil {
		return
	}
	lineBotTokenStr, err := base64.StdEncoding.DecodeString(string(lineBot.LineBotToken[0]))
	if err != nil {
		return
	}
	lineGroupStr, err := base64.StdEncoding.DecodeString(string(lineBot.LineGroupID[0]))
	if err != nil {
		return
	}

	// 暗号化されたシークレットキーの復号化
	sercretKey, err := crypto.Decrypt(lineBotSecretStr, keyBytes, decodeBotSecret)
	if err != nil {
		return nil, err
	}

	// macの生成
	mac := hmac.New(sha256.New, []byte(sercretKey))
	mac.Write(requestBodyByte)
	validSignByte := mac.Sum(nil)

	signature := base64.StdEncoding.EncodeToString(validSignByte)

	// 署名が一致しない場合は両方nilを返す
	if header != signature {
		return nil, nil
	}
	lineNotifyTokenByte, err := crypto.Decrypt(lineNotifyStr, keyBytes, decodeNotifyToken)
	if err != nil {
		return nil, err
	}
	lineBotTokenByte, err := crypto.Decrypt(lineBotTokenStr, keyBytes, decodeBotToken)
	if err != nil {
		return nil, err
	}
	lineGroupByte, err := crypto.Decrypt(lineGroupStr, keyBytes, decodeGroupID)
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

