package des

import (
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// EncryptECB DES加密
// input 要加密的字符串	skey 加密使用的秘钥[字符串长度必须是8的倍数]
func EncryptECB(input string, skey string) (r string, err error) {
	origData := []byte(input)
	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("des NewCipher err:%v", err)
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := NewECBEncrypter(block)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	r = strings.ToUpper(hex.EncodeToString(crypted))
	return
}

// DecryptECB DES解密
// input 要解密的字符串	skey 加密使用的秘钥[字符串长度必须是8的倍数]
func DecryptECB(input string, skey string) (r string, err error) {
	/*add by champly 2016年11月16日17:35:03*/
	if len(input) < 1 {
		return "", errors.New("解密的对象长度必须大于0")
	}
	/*end*/

	crypted, err := hex.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("hex DecodeString err:%v", err)
	}
	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("des NewCipher err:%v", err)
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	r = string(origData)
	return
}
