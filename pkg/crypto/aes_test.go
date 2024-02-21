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
		key, err := hex.DecodeString(keyString)
		assert.NoError(t, err)
		assert.Equal(t, []uint8([]byte{0x64, 0x5e, 0x73, 0x9a, 0x7f, 0x9f, 0x16, 0x27, 0x25, 0xc1, 0x53, 0x3d, 0xc2, 0xc5, 0xe8, 0x27}), key)

		iv, encrypted, err := Encrypt(plain, key)
		assert.NoError(t, err)

		fmt.Println("IV:", hex.EncodeToString(iv))
		fmt.Println("Encrypted:", base64.StdEncoding.EncodeToString(encrypted))
		decrypted, err := Decrypt(encrypted, key, iv)
		assert.NoError(t, err)
		fmt.Println(text, string(decrypted))
		assert.Equal(t, text, string(decrypted))
	})

	t.Run("空文字の暗号化と復号化が正常に行われること", func(t *testing.T) {
		key, err := hex.DecodeString(keyString)
		assert.NoError(t, err)
		assert.Equal(t, []uint8([]byte{0x64, 0x5e, 0x73, 0x9a, 0x7f, 0x9f, 0x16, 0x27, 0x25, 0xc1, 0x53, 0x3d, 0xc2, 0xc5, 0xe8, 0x27}), key)

		iv, encrypted, err := Encrypt(emptyPlain, key)
		assert.NoError(t, err)

		fmt.Println("IV:", hex.EncodeToString(iv))
		fmt.Println("Encrypted:", base64.StdEncoding.EncodeToString(encrypted))
		decrypted, err := Decrypt(encrypted, key, iv)
		assert.NoError(t, err)
		fmt.Println(emptyText, string(decrypted))
		assert.Equal(t, emptyText, string(decrypted))
	})

	t.Run("ivの暗号化と復号化が正常に行われること", func(t *testing.T) {
		key, err := hex.DecodeString(keyString)
		assert.NoError(t, err)
		assert.Equal(t, []uint8([]byte{0x64, 0x5e, 0x73, 0x9a, 0x7f, 0x9f, 0x16, 0x27, 0x25, 0xc1, 0x53, 0x3d, 0xc2, 0xc5, 0xe8, 0x27}), key)

		iv, encrypted, err := Encrypt(tokenPlain, key)
		assert.NoError(t, err)

		encodeIv := hex.EncodeToString(iv)
		decodeIv, err := hex.DecodeString(encodeIv)
		assert.NoError(t, err)

		fmt.Println("IV:", hex.EncodeToString(iv))
		fmt.Println("Encrypted:", base64.StdEncoding.EncodeToString(encrypted))
		decrypted, err := Decrypt(encrypted, key, decodeIv)
		assert.NoError(t, err)
		fmt.Println(tokenText, string(decrypted))
		assert.Equal(t, tokenText, string(decrypted))
	})

	t.Run("iv2の暗号化と復号化が正常に行われること", func(t *testing.T) {
		key, err := hex.DecodeString(keyString)
		assert.NoError(t, err)
		assert.Equal(t, []uint8([]byte{0x64, 0x5e, 0x73, 0x9a, 0x7f, 0x9f, 0x16, 0x27, 0x25, 0xc1, 0x53, 0x3d, 0xc2, 0xc5, 0xe8, 0x27}), key)

		decodeIv, err := hex.DecodeString("76a9cfafaaaf35c1d337ab5dc113d1ce")
		assert.NoError(t, err)

		decodeStr, err := base64.StdEncoding.DecodeString("Mwq69HWKIpjYluIwRFihWQ/y8/PYpbTFNR8dZ7GtdNg0N0lg0yZSt6m6iHjf8ZRcJBkCbMhrPM4cZ5spGtKhE0HUC76ud0NEHmEsLgu6LYRN8cuMgRRjbEy52+9BlDCXo12vqGpmL78GA3Yl/JfArbIicUjYgZxo1ofYXHUenb9xjMwyQYmd1sUg5g8ntEnAPEnVpxHxdRVfB0qtGFLn9RatBiAk/ZlrN0O7B9yPpAw=")
		assert.NoError(t, err)

		fmt.Println("IV:", decodeIv)
		fmt.Println("Encrypted:", decodeStr)
		decrypted, err := Decrypt(decodeStr, key, decodeIv)
		assert.NoError(t, err)
		fmt.Println(tokenText, string(decrypted))
		assert.Equal(t, tokenText, string(decrypted))
	})
}
