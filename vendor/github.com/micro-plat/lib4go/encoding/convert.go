package encoding

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//Encode 将UTF8字符串编码为gbk或gb2312格式
func Encode(content string, e string) (result []byte, err error) {
	return EncodeBytes([]byte(content), e)
}

//EncodeBytes 将UTF8字符串编码为gbk或gb2312格式
func EncodeBytes(buff []byte, e string) (result []byte, err error) {
	reader := GetEncodeReader(buff, e)
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("编码转换失败:content:%s, err:%+v", string(buff), err)
		return
	}
	return d, nil
}

//Decode 根据编码进行解码操作
func Decode(content string, e string) (result []byte, err error) {
	return DecodeBytes([]byte(content), e)
}

//DecodeBytes 根据编码进行解码操作
func DecodeBytes(buff []byte, e string) (result []byte, err error) {
	reader := GetDecodeReader(buff, e)
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("编码转换失败:content:%s, err:%+v", string(buff), err)
		return
	}
	return d, nil
}

// GetDecodeReader 获取
// charset不区分大小写
func GetDecodeReader(buff []byte, charset string) io.Reader {
	charset = strings.ToLower(charset)
	if strings.EqualFold(charset, "gbk") || strings.EqualFold(charset, "gb2312") {
		return transform.NewReader(bytes.NewReader(buff), simplifiedchinese.GBK.NewDecoder())
	}
	return strings.NewReader(string(buff))
}

// GetEncodeReader 获取
func GetEncodeReader(buff []byte, charset string) io.Reader {
	charset = strings.ToLower(charset)
	if strings.EqualFold(charset, "gbk") || strings.EqualFold(charset, "gb2312") {
		return transform.NewReader(bytes.NewReader(buff), simplifiedchinese.GBK.NewEncoder())
	}
	return strings.NewReader(string(buff))
}
