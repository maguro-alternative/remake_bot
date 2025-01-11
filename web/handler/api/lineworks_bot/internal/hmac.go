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

	// Encode to BASE64
	encodedB64Body := base64.StdEncoding.EncodeToString(encodedBody)

	// Compare signatures
	return encodedB64Body == signature
}
