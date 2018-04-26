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

// GetReader 获取
// charset不区分大小写
func GetReader(content string, charset string) io.Reader {
	charset = strings.ToLower(charset)
	if strings.EqualFold(charset, "gbk") || strings.EqualFold(charset, "gb2312") {
		return transform.NewReader(bytes.NewReader([]byte(content)), simplifiedchinese.GBK.NewDecoder())
	}
	return strings.NewReader(content)

}
func ConvertBytes(data []byte, encoding string) (buffer []byte, err error) {
	encoding = strings.ToLower(encoding)
	if !strings.EqualFold(encoding, "gbk") && !strings.EqualFold(encoding, "gb2312") &&
		!strings.EqualFold(encoding, "utf-8") {
		err = fmt.Errorf("不支持的编码方式：%s", encoding)
		return
	}
	//转换utf-8格式
	if strings.EqualFold(encoding, "utf-8") {
		buffer = data
		return
	}

	//转换gbk gb2312格式
	buffer, err = ioutil.ReadAll(transform.NewReader(bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder()))
	return
}

// Convert []byte转换为字符串
// encoding 将utf-8格式数据，转换为其它格式 支持gbk，gb2312，utf-8	不区分大小写
func Convert(data []byte, encoding string) (content string, err error) {
	buffer, err := ConvertBytes(data, encoding)
	if err != nil {
		return
	}
	content = string(buffer)
	return
}
