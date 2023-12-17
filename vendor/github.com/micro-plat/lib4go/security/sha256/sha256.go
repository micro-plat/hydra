package sha256

import (
	"crypto/sha256"
	"fmt"
)

// Encrypt sha256加密
func Encrypt(content string) string {
	h := sha256.New()
	h.Write([]byte(content))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

//EncryptBytes bytes加密
func EncryptBytes(content []byte) []byte {
	h := sha256.New()
	h.Write(content)
	return h.Sum(nil)
}
