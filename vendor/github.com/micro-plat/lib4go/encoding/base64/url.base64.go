package base64

import "encoding/base64"

//URLEncodeBytes 把一个[]byte通过base64编码成string
func URLEncodeBytes(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

//URLDecodeBytes 把一个string通过base64解码成[]byte
func URLDecodeBytes(src string) (s []byte, err error) {
	s, err = base64.URLEncoding.DecodeString(src)
	return
}

//URLEncode 把一个string通过base64编码
func URLEncode(src string) string {
	return URLEncodeBytes([]byte(src))
}

//URLDecode 把一个string通过base64解码
func URLDecode(src string) (s string, err error) {
	buf, err := URLDecodeBytes(src)
	if err != nil {
		return
	}
	s = string(buf)
	return
}
