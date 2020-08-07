package sha1

import (
	"crypto/sha1"
	"fmt"
)

// Encrypt SHA1加密
func Encrypt(content string) string {
	h := sha1.New()
	h.Write([]byte(content))
	bs := h.Sum(nil)
	h.Reset()
	r := fmt.Sprintf("%x", bs)
	return r
}
