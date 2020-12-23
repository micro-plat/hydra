package des

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
	"strings"
)

// Encrypt DES加密
// input 要加密的字符串	skey 加密使用的秘钥[字符串长度必须是8的倍数]
func Encrypt(input string, skey string, mode string) (r string, err error) {
	origData := []byte(input)
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	crypted, err := EncryptBytes(origData, skey, iv, mode)
	if err != nil {
		return
	}
	r = strings.ToUpper(hex.EncodeToString(crypted))
	return
}

// Decrypt DES解密
// input 要解密的字符串	skey 加密使用的秘钥[字符串长度必须是8的倍数]
func Decrypt(input string, skey string, mode string) (r string, err error) {
	crypted, err := hex.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("hex DecodeString err:%v", err)
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	origData, err := DecryptBytes(crypted, skey, iv, mode)
	if err != nil {
		return
	}
	r = string(origData)
	return
}

// EncryptBytes DES加密
// input 要加密的字符串	skey 加密使用的秘钥[字符串长度必须是8的倍数]
func EncryptBytes(origData []byte, skey string, iv []byte, mode string) (crypted []byte, err error) {
	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		err = fmt.Errorf("des NewCipher err:%v", err)
		return
	}
	m, p, err := getModePadding(mode)
	if err != nil {
		return
	}
	var blockMode cipher.BlockMode
	switch m {
	case "ECB":
		blockMode = NewECBEncrypter(block)
	case "CBC":
		blockMode = cipher.NewCBCEncrypter(block, iv)
	default:
		err = fmt.Errorf("加密模式不支持:%s", m)
		return
	}
	switch p {
	case "PKCS5":
		origData = PKCS5Padding(origData, block.BlockSize())
	case "PKCS7":
		origData = PKCS7Padding(origData)
	case "ZERO":
		origData = ZeroPadding(origData, block.BlockSize())
	default:
		err = fmt.Errorf("填充模式不支持:%s", p)
		return
	}
	crypted = make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return
}

// DecryptBytes DES解密
// input 要解密的字符串	skey 加密使用的秘钥[字符串长度必须是8的倍数]
func DecryptBytes(crypted []byte, skey string, iv []byte, mode string) (r []byte, err error) {

	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		err = fmt.Errorf("des NewCipher err:%v", err)
		return
	}
	m, p, err := getModePadding(mode)
	if err != nil {
		return
	}
	var blockMode cipher.BlockMode
	switch m {
	case "CBC":
		blockMode = cipher.NewCBCDecrypter(block, iv)
	case "ECB":
		blockMode = NewECBDecrypter(block)
	default:
		err = fmt.Errorf("加密模式不支持:%s", m)
		return
	}
	r = make([]byte, len(crypted))
	blockMode.CryptBlocks(r, crypted)
	switch p {
	case "PKCS5":
		r = PKCS5UnPadding(r)
	case "PKCS7":
		r = PKCS7UnPadding(r)
	case "ZERO":
		r = ZeroUnPadding(r)
	default:
		err = fmt.Errorf("填充模式不支持:%s", p)
		return
	}

	return
}

func getModePadding(name string) (mode, padding string, err error) {
	names := strings.Split(name, "/")
	if len(names) != 2 {
		err = fmt.Errorf("输入模式不正确:%s", name)
		return
	}
	mode = strings.ToUpper(names[0])
	padding = strings.ToUpper(names[1])
	if mode != "CBC" && mode != "ECB" {
		err = fmt.Errorf("加密模式不支持:%s", mode)
		return
	}
	if padding != "PKCS5" && padding != "PKCS7" && padding != "ZERO" {
		err = fmt.Errorf("填充模式不支持:%s", padding)
		return
	}
	return
}
