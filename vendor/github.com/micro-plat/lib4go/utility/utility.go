package utility

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/micro-plat/lib4go/security/md5"
)

// GetGUID 生成Guid字串
func GetGUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return md5.Encrypt(base64.URLEncoding.EncodeToString(b))
}
