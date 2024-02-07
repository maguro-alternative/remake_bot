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
	keyString := "645E739A7F9F162725C1533DC2C5E827"
	plain := []byte(text)
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
		assert.Equal(t, text, string(decrypted))
	})
}
