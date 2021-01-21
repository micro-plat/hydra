package conf

import (
	"encoding/hex"
	"fmt"

	"github.com/micro-plat/lib4go/security/des"
)

const confKey = "@-hydra*"
const confIV = "#*---iv*"
const hd = "encrypt"
const mode = "cbc/pkcs5"

//Encrypt 对加密内容进行des加密，并增加加密头
func Encrypt(input []byte) string {
	v, _ := des.EncryptBytes(input, confKey, []byte(confIV), mode)
	return fmt.Sprintf("%s:%s:%s", hd, mode, hex.EncodeToString(v))
}

//Decrypt 检查是否包含加密头，报含则根据加密头数据解密数据
func Decrypt(data []byte) ([]byte, error) {
	lheader := len(hd)
	lmode := len(mode)
	if len(data) <= lheader+lmode+2 {
		return data, nil
	}
	if string(data[0:lheader]) != hd {
		return data, nil
	}
	mode := string(data[lheader+1 : lheader+lmode+1])
	src := make([]byte, (len(data)-lheader-lmode-2)/2)
	if _, err := hex.Decode(src, data[(lheader+lmode+2):]); err != nil {
		return nil, err
	}
	return des.DecryptBytes(src, confKey, []byte(confIV), mode)
}
