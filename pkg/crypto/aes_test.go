package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateIV(t *testing.T) {
	text := "this is test message"
	emptyText := ""
	tokenText := "JXYMWQDKCVFLBZTEINAHORSGUPXJWODLKCVMFRGBYAHITZSNQUOJIPXVBLRWCDZGAFYXJTKUHJXYMWQDKCVFLBZTEINAHORSGUPXJWODLKCVMFRGBYAHITZSNQUOJIPXVBLRWCDZGAFYXJTKUHaskwsrhuivoanufinaoenfuyvow"
	keyString := "645E739A7F9F162725C1533DC2C5E827"
	plain := []byte(text)
	emptyPlain := []byte(emptyText)
	tokenPlain := []byte(tokenText)
	t.Run("暗号化と復号化が正常に行われること", func(t *testing.T) {
		aesCrypto, err := NewAESCrypto(keyString)
		assert.NoError(t, err)

		iv, encrypted, err := aesCrypto.Encrypt(plain)
		assert.NoError(t, err)

		fmt.Println("IV:", hex.EncodeToString(iv))
		fmt.Println("Encrypted:", base64.StdEncoding.EncodeToString(encrypted))
		decrypted, err := aesCrypto.Decrypt(encrypted, iv)
		assert.NoError(t, err)
		fmt.Println(text, string(decrypted))
		assert.Equal(t, text, string(decrypted))
	})

	t.Run("空文字の暗号化と復号化が正常に行われること", func(t *testing.T) {
		aesCrypto, err := NewAESCrypto(keyString)
		assert.NoError(t, err)

		iv, encrypted, err := aesCrypto.Encrypt(emptyPlain)
		assert.NoError(t, err)

		fmt.Println("IV:", hex.EncodeToString(iv))
		fmt.Println("Encrypted:", base64.StdEncoding.EncodeToString(encrypted))
		decrypted, err := aesCrypto.Decrypt(encrypted, iv)
		assert.NoError(t, err)
		fmt.Println(emptyText, string(decrypted))
		assert.Equal(t, emptyText, string(decrypted))
	})

	t.Run("ivの暗号化と復号化が正常に行われること", func(t *testing.T) {
		aesCrypto, err := NewAESCrypto(keyString)
		assert.NoError(t, err)

		iv, encrypted, err := aesCrypto.Encrypt(tokenPlain)
		assert.NoError(t, err)

		encodeIv := hex.EncodeToString(iv)
		decodeIv, err := hex.DecodeString(encodeIv)
		assert.NoError(t, err)

		fmt.Println("IV:", hex.EncodeToString(iv))
		fmt.Println("Encrypted:", base64.StdEncoding.EncodeToString(encrypted))
		decrypted, err := aesCrypto.Decrypt(encrypted, decodeIv)
		assert.NoError(t, err)
		fmt.Println(tokenText, string(decrypted))
		assert.Equal(t, tokenText, string(decrypted))
	})

	t.Run("iv2の暗号化と復号化が正常に行われること", func(t *testing.T) {
		aesCrypto, err := NewAESCrypto(keyString)
		assert.NoError(t, err)

		decodeIv, err := hex.DecodeString("76a9cfafaaaf35c1d337ab5dc113d1ce")
		assert.NoError(t, err)

		decodeStr, err := base64.StdEncoding.DecodeString("Mwq69HWKIpjYluIwRFihWQ/y8/PYpbTFNR8dZ7GtdNg0N0lg0yZSt6m6iHjf8ZRcJBkCbMhrPM4cZ5spGtKhE0HUC76ud0NEHmEsLgu6LYRN8cuMgRRjbEy52+9BlDCXo12vqGpmL78GA3Yl/JfArbIicUjYgZxo1ofYXHUenb9xjMwyQYmd1sUg5g8ntEnAPEnVpxHxdRVfB0qtGFLn9RatBiAk/ZlrN0O7B9yPpAw=")
		assert.NoError(t, err)

		fmt.Println("IV:", decodeIv)
		fmt.Println("Encrypted:", decodeStr)
		decrypted, err := aesCrypto.Decrypt(decodeStr, decodeIv)
		assert.NoError(t, err)
		fmt.Println(tokenText, string(decrypted))
		assert.Equal(t, tokenText, string(decrypted))
	})
}
