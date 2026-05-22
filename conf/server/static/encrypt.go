package static

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"net/http"
)

// computeKeyFingerprint 生成密钥指纹（SHA256 前8字节 hex）
func computeKeyFingerprint(key string) string {
	h := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", h[:8])
}

// padKey 填充或截断 key 到 32 字节
func padKey(key string) []byte {
	keyBytes := []byte(key)
	if len(keyBytes) >= 32 {
		return keyBytes[:32]
	}
	padded := make([]byte, 32)
	copy(padded, keyBytes)
	return padded
}

// padIV 填充或截断 IV 到 12 字节
func padIV(iv string) []byte {
	ivBytes := []byte(iv)
	if len(ivBytes) >= 12 {
		return ivBytes[:12]
	}
	padded := make([]byte, 12)
	copy(padded, ivBytes)
	return padded
}

// EncryptAndCompress 先 Gzip 压缩再 AES-256-GCM 加密
// 使用固定 IV，输出仅 [密文+Tag]（IV 不在响应体中）
func EncryptAndCompress(data []byte, key string, iv string) ([]byte, error) {
	// 1. Gzip 压缩
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if _, err := gw.Write(data); err != nil {
		return nil, fmt.Errorf("gzip压缩失败: %w", err)
	}
	if err := gw.Close(); err != nil {
		return nil, fmt.Errorf("gzip关闭失败: %w", err)
	}
	compressed := buf.Bytes()

	// 2. AES-256-GCM 加密
	block, err := aes.NewCipher(padKey(key))
	if err != nil {
		return nil, fmt.Errorf("创建AES失败: %w", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM失败: %w", err)
	}

	nonce := padIV(iv)
	sealed := aesGCM.Seal(nil, nonce, compressed, nil)
	return sealed, nil
}

// DoEncrypt 加密并压缩数据（使用配置的密钥和IV）
func (s *Static) DoEncrypt(data []byte) ([]byte, error) {
	return EncryptAndCompress(data, s.EncryptKey, s.EncryptIV)
}

// IsRemoteEncrypted 检查远程响应是否已加密且密钥匹配
func IsRemoteEncrypted(remoteHeaders http.Header, localKeyFingerprint string) bool {
	if remoteHeaders.Get("X-Content-Crypto") != "aes-256-gcm" {
		return false
	}
	remoteFingerprint := remoteHeaders.Get("X-Content-Crypto-Key")
	return remoteFingerprint != "" && remoteFingerprint == localKeyFingerprint
}
