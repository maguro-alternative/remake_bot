package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineWorksValidateRequest(t *testing.T) {
	body := []byte(`{"type":"message"}`)
	botSecret := "secret"
	mac := hmac.New(sha256.New, []byte(botSecret))
	mac.Write(body)
	validSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	t.Run("署名が一致する場合trueを返すこと", func(t *testing.T) {
		assert.True(t, LineWorksValidateRequest(body, validSignature, botSecret))
	})

	t.Run("署名が一致しない場合falseを返すこと", func(t *testing.T) {
		assert.False(t, LineWorksValidateRequest(body, validSignature, "other"))
	})

	t.Run("署名がbase64でない場合falseを返すこと", func(t *testing.T) {
		assert.False(t, LineWorksValidateRequest(body, "not base64!!", botSecret))
	})

	t.Run("署名が空の場合falseを返すこと", func(t *testing.T) {
		assert.False(t, LineWorksValidateRequest(body, "", botSecret))
	})
}
