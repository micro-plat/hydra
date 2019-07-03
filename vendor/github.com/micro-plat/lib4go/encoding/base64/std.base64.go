package base64

import "encoding/base64"

// EncodeBytes 把一个[]byte通过base64编码成string
func EncodeBytes(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// DecodeBytes 把一个string通过base64解码成[]byte
func DecodeBytes(src string) (s []byte, err error) {
	s, err = base64.StdEncoding.DecodeString(src)
	return
}

// Encode 把一个string通过base64编码
func Encode(src string) string {
	return EncodeBytes([]byte(src))
}

// Decode 把一个string通过base64解码
func Decode(src string) (s string, err error) {
	buf, err := DecodeBytes(src)
	if err != nil {
		return
	}
	s = string(buf)
	return
}
