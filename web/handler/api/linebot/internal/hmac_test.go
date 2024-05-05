package internal

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestHmac_LineHmac(t *testing.T) {
	privateKey := "645E739A7F9F162725C1533DC2C5E827"
	decodeNotifyToken, err := hex.DecodeString("aa7c5fe80002633327f0fefe67a565de")
	assert.NoError(t, err)
	lineNotifyStr, err := base64.StdEncoding.DecodeString(string([]byte("X+P6kmO6DnEjM3TVqXkwNA==")))
	assert.NoError(t, err)

	decodeBotToken, err := hex.DecodeString("baeff317cb83ef55b193b6d3de194124")
	assert.NoError(t, err)
	lineBotStr, err := base64.StdEncoding.DecodeString(string([]byte("uy2qtvYTnSoB5qIntwUdVQ==")))
	assert.NoError(t, err)

	decodeBotSecret, err := hex.DecodeString("0ffa8ed72efcb5f1d834e4ce8463a62c")
	assert.NoError(t, err)
	lineBotSecretStr, err := base64.StdEncoding.DecodeString(string([]byte("i2uHQCyn58wRR/b03fRw6w==")))
	assert.NoError(t, err)

	decodeGroupID, err := hex.DecodeString("e14db710b23520766fd652c0f19d437a")
	assert.NoError(t, err)
	lineGroupStr, err := base64.StdEncoding.DecodeString(string([]byte("YgexFQQlLcaXmsw9mFN35Q==")))
	assert.NoError(t, err)

	lineBot := &repository.LineBot{
		LineNotifyToken: pq.ByteaArray{lineNotifyStr},
		LineBotToken:    pq.ByteaArray{lineBotStr},
		LineBotSecret:   pq.ByteaArray{lineBotSecretStr},
		LineGroupID:     pq.ByteaArray{lineGroupStr},
	}
	lineBotIv := repository.LineBotIvNotClient{
		LineNotifyTokenIv: pq.ByteaArray{decodeNotifyToken},
		LineBotTokenIv:    pq.ByteaArray{decodeBotToken},
		LineBotSecretIv:   pq.ByteaArray{decodeBotSecret},
		LineGroupIDIv:     pq.ByteaArray{decodeGroupID},
	}

	aesCrypto := &crypto.AESMock{
		EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
			return nil, nil, nil
		},
		DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
			return nil, nil
		},
	}
	t.Run("正常系", func(t *testing.T) {
		requestBodyByte := []byte(`{"destination":"U940f73798c755e4b1b86eb6a5adaba23","events":[]}`)
		header := "Q0v3M1tK3Ic7SsDvX6KAZJr2M9Jwm/2WCzoz2EwqjGc="

		decrypt, err := LineHmac(privateKey, requestBodyByte, aesCrypto, lineBot, lineBotIv, header)
		assert.NoError(t, err)
		assert.NotNil(t, decrypt)
	})
	t.Run("異常系", func(t *testing.T) {
		requestBodyByte := []byte(`{"events":[{"replyToken":"","type":"message","timestamp":0,"source":{"userId":"Udeadbw00dbaadbeefdeadbeefdeadbeef","type":"user"},"message":{"type":"text","id":"1234567890","text":"Hello, world"}}]}`)
		header := "12345678901234567890123456789012" //6eMInZT4CEsIf/P5Iv+9VmezoOPqXs1il6R4QjtUG4o=

		decrypt, err := LineHmac(privateKey, requestBodyByte, aesCrypto, lineBot, lineBotIv, header)
		assert.NoError(t, err)
		assert.Nil(t, decrypt)
	})
}
