package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func LineWorksValidateRequest(body []byte, signature string, botSecret string) bool {
	// Convert botSecret to byte array
	secretKey := []byte(botSecret)

	// Create HMAC-SHA256 hash
	h := hmac.New(sha256.New, secretKey)
	h.Write(body)
	encodedBody := h.Sum(nil)

	// Compare signatures in constant time to prevent timing attacks
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	return hmac.Equal(encodedBody, decodedSignature)
}
