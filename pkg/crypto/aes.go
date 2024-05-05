package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

func generateIV() ([]byte, error) {
	/*IV (128 bit) の生成*/
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil
}

func pkcs7Pad(data []byte) []byte {
	/*平文 []byte の長さが16の倍数ではない場合、16の倍数にするためにパディングする*/
	length := aes.BlockSize - (len(data) % aes.BlockSize)
	trailing := bytes.Repeat([]byte{byte(length)}, length)
	return append(data, trailing...)
}

func pkcs7Unpad(data []byte) []byte {
	/*パディングを削除する*/
	dataLength := len(data)
	padLength := int(data[dataLength-1])
	return data[:dataLength-padLength]
}

type AES struct {
	key []byte
}

func NewAESCrypto(key string) (*AES, error) {
	k, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return &AES{key: k}, nil
}

// Encrypt は、dataをAES暗号化します。
func (a AES) Encrypt(data []byte) (iv []byte, encrypted []byte, err error) {
	/*AES暗号化*/
	// IV (Initialization Vector) の生成
	iv, err = generateIV()
	if err != nil {
		return nil, nil, err
	}
	// 暗号化
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, nil, err
	}
	// パディング
	padded := pkcs7Pad(data)
	encrypted = make([]byte, len(padded))
	// 暗号化
	cbcEncrypter := cipher.NewCBCEncrypter(block, iv)
	cbcEncrypter.CryptBlocks(encrypted, padded)
	return iv, encrypted, nil
}

// Decrypt は、dataをAES復号化します。
func (a AES) Decrypt(data []byte, iv []byte) ([]byte, error) {
	/*AES復号化*/
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	// 復号化
	cbcDecrypter := cipher.NewCBCDecrypter(block, iv)
	cbcDecrypter.CryptBlocks(decrypted, data)
	return pkcs7Unpad(decrypted), nil
}

type AESMock struct {
	EncryptFunc func(data []byte) (iv []byte, encrypted []byte, err error)
	DecryptFunc func(data []byte, iv []byte) (decrypted []byte, err error)
}

func (a AESMock) Encrypt(data []byte) (iv []byte, encrypted []byte, err error) {
	return a.EncryptFunc(data)
}

func (a AESMock) Decrypt(data []byte, iv []byte) (decrypted []byte, err error) {
	return a.DecryptFunc(data, iv)
}

type AESInterface interface {
	Encrypt(data []byte) (iv []byte, encrypted []byte, err error)
	Decrypt(data []byte, iv []byte) (decrypted []byte, err error)
}

var (
	_ AESInterface = (*AES)(nil)
	_ AESInterface = (*AESMock)(nil)
)
