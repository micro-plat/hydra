package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// Encrypt MD5加密
func Encrypt(s string) string {
	return EncryptBytes([]byte(s))
}
func EncryptBytes(buffer []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(buffer)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
