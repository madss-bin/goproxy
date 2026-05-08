package crypto

import (
	"animex/goproxy/config"
	"encoding/base64"
)

func xorTransform(data []byte) {
	key := config.XOR_KEY
	keyLen := len(key)
	if keyLen == 0 {
		return
	}
	for i := range data {
		data[i] ^= key[i%keyLen]
	}
}

func EncryptURL(url string) string {
	data := []byte(url)
	xorTransform(data)
	return base64.RawURLEncoding.EncodeToString(data)
}

func DecryptURL(encrypted string) (string, bool) {
	data, err := base64.RawURLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", false
	}
	xorTransform(data)
	return string(data), true
}
