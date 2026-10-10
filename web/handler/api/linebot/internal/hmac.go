package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/samber/mo"

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
) (decrypt mo.Option[*LineBotDecrypt], err error) {
	lineBotDecrypt := &LineBotDecrypt{}

	lineBotSecretKey, err := aesCrypto.Decrypt(lineBot.LineBotSecret[0], lineBotIv.LineBotSecretIv[0])
	if err != nil {
		return mo.None[*LineBotDecrypt](), err
	}

	// macの生成
	mac := hmac.New(sha256.New, []byte(lineBotSecretKey))
	mac.Write(requestBodyByte)
	validSignByte := mac.Sum(nil)

	// 署名が一致しない場合は None を返す
	// 正しい署名をログに出すと偽造に使われるため出力しない。比較はタイミング攻撃を防ぐため定数時間で行う
	headerSignByte, err := base64.StdEncoding.DecodeString(header)
	if err != nil || !hmac.Equal(headerSignByte, validSignByte) {
		return mo.None[*LineBotDecrypt](), nil
	}
	lineNotifyTokenByte, err := aesCrypto.Decrypt(lineBot.LineNotifyToken[0], lineBotIv.LineNotifyTokenIv[0])
	if err != nil {
		return mo.None[*LineBotDecrypt](), err
	}
	lineBotTokenByte, err := aesCrypto.Decrypt(lineBot.LineBotToken[0], lineBotIv.LineBotTokenIv[0])
	if err != nil {
		return mo.None[*LineBotDecrypt](), err
	}
	lineGroupByte, err := aesCrypto.Decrypt(lineBot.LineGroupID[0], lineBotIv.LineGroupIDIv[0])
	if err != nil {
		return mo.None[*LineBotDecrypt](), err
	}
	lineBotDecrypt.LineNotifyToken = string(lineNotifyTokenByte)
	lineBotDecrypt.LineBotToken = string(lineBotTokenByte)
	lineBotDecrypt.LineGroupID = string(lineGroupByte)
	lineBotDecrypt.DefaultChannelID = lineBot.DefaultChannelID
	lineBotDecrypt.DebugMode = lineBot.DebugMode
	return mo.Some(lineBotDecrypt), nil
}
