package des

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/micro-plat/lib4go/security/padding"
	"github.com/micro-plat/lib4go/types"
)

const (
	DesECB = "ECB"
	DesCBC = "CBC"
)

// Encrypt DES加密
// mode 加密类型/填充模式,不传默认为:CFB/ZERO
// input 要加密的字符串	skey 加密使用的秘钥[字符串长度必须是8]
func Encrypt(input string, skey string, mode ...string) (r string, err error) {
	origData := []byte(input)
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	crypted, err := EncryptBytes(origData, skey, iv, mode...)
	if err != nil {
		return
	}
	r = strings.ToUpper(hex.EncodeToString(crypted))
	return
}

// Decrypt DES解密
// mode 加密类型/填充模式,不传默认为:CFB/ZERO
// input 要解密的字符串	skey 加密使用的秘钥[字符串长度必须是8]
func Decrypt(input string, skey string, mode ...string) (r string, err error) {
	crypted, err := hex.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("hex DecodeString err:%v", err)
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	origData, err := DecryptBytes(crypted, skey, iv, mode...)
	if err != nil {
		return
	}
	r = string(origData)
	return
}

// EncryptBytes DES加密
// mode 加密类型/填充模式,不传默认为:CFB/ZERO
// input 要加密的字符串	skey 加密使用的秘钥[字符串长度必须是8]
func EncryptBytes(origData []byte, skey string, iv []byte, mode ...string) (crypted []byte, err error) {
	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		err = fmt.Errorf("des NewCipher err:%v", err)
		return
	}
	cmode := types.GetStringByIndex(mode, 0, fmt.Sprintf("%s/%s", DesECB, padding.PaddingZero))

	m, p, err := padding.GetModePadding(cmode)
	if err != nil {
		return nil, err
	}

	var blockMode cipher.BlockMode
	switch m {
	case DesECB:
		blockMode = NewECBEncrypter(block)
	case DesCBC:
		blockMode = cipher.NewCBCEncrypter(block, iv)
	default:
		err = fmt.Errorf("加密模式不支持:%s", m)
		return
	}
	switch p {
	case padding.PaddingPkcs5:
		origData = padding.PKCS5Padding(origData, block.BlockSize())
	case padding.PaddingPkcs7:
		origData = padding.PKCS7Padding(origData)
	case padding.PaddingZero:
		origData = padding.ZeroPadding(origData, block.BlockSize())
	default:
		err = fmt.Errorf("填充模式不支持:%s", p)
		return
	}
	crypted = make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return
}

// DecryptBytes DES解密
// mode 加密类型/填充模式,不传默认为:CFB/ZERO
// input 要解密的字符串	skey 加密使用的秘钥[字符串长度必须是8]
func DecryptBytes(crypted []byte, skey string, iv []byte, mode ...string) (r []byte, err error) {

	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		err = fmt.Errorf("des NewCipher err:%v", err)
		return
	}
	cmode := types.GetStringByIndex(mode, 0, fmt.Sprintf("%s/%s", DesECB, padding.PaddingZero))

	m, p, err := padding.GetModePadding(cmode)
	if err != nil {
		return nil, err
	}

	var blockMode cipher.BlockMode
	switch m {
	case DesCBC:
		blockMode = cipher.NewCBCDecrypter(block, iv)
	case DesECB:
		blockMode = NewECBDecrypter(block)
	default:
		err = fmt.Errorf("加密模式不支持:%s", m)
		return
	}
	r = make([]byte, len(crypted))
	blockMode.CryptBlocks(r, crypted)
	switch p {
	case padding.PaddingPkcs5:
		r = padding.PKCS5UnPadding(r)
	case padding.PaddingPkcs7:
		r = padding.PKCS7UnPadding(r)
	case padding.PaddingZero:
		r = padding.ZeroUnPadding(r)
	default:
		err = fmt.Errorf("填充模式不支持:%s", p)
		return
	}

	return
}
