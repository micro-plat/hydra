package conf

import (
	"bytes"
	"encoding/json"

	"github.com/micro-plat/lib4go/security/md5"
	"github.com/micro-plat/lib4go/types"
)

//IConf 配置管理
type IConf interface {
	types.IXMap
	GetVersion() int32
	GetSignature() string
}

//EmptyRawConf 空的rawconf
var EmptyRawConf = &RawConf{
	version:   0,
	raw:       []byte("{}"),
	signature: md5.EncryptBytes([]byte("{}")),
	XMap:      types.NewXMap(),
}

//RawConf json配置文件
type RawConf struct {
	raw       []byte
	version   int32
	signature string
	types.XMap
}

//NewByMap 根据map初始化对象
func NewByMap(data map[string]interface{}, version int32) (c *RawConf, err error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	c = &RawConf{
		XMap:      types.NewXMapByMap(data),
		version:   version,
		raw:       buf,
		signature: md5.EncryptBytes(buf),
	}
	return c, nil
}

//NewByText 初始化RawConf
func NewByText(message []byte, version int32) (c *RawConf, err error) {
	c = &RawConf{
		raw:       message,
		signature: md5.EncryptBytes(message),
		version:   version,
	}
	switch {
	case bytes.HasPrefix(message, []byte("<?xml")):
		c.XMap, err = types.NewXMapByXML(string(message))
	case bytes.HasPrefix(message, []byte("{")) || bytes.HasPrefix(message, []byte("[")):
		c.XMap, err = types.NewXMapByJSON(string(message))
	}
	return c, err
}

//GetRaw 获取原串
func (j *RawConf) GetRaw() []byte {
	return j.raw
}

//GetVersion 获取当前配置的版本号
func (j *RawConf) GetVersion() int32 {
	return j.version
}

//GetSignature 获取当前对象的唯一标识
func (j *RawConf) GetSignature() string {
	return j.signature
}
