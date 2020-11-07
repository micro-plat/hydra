package conf

import (
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

//EmptyJSONConf 空的jsonconf
var EmptyJSONConf = &RawConf{
	raw:       []byte("{}"),
	signature: md5.EncryptBytes([]byte("{}")),
	XMap:      types.NewXMap(),
}

//RawConf json配置文件
type RawConf struct {
	raw       json.RawMessage
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

//NewByJSON 初始化JsonConf
func NewByJSON(message []byte, version int32) (c *RawConf, err error) {
	c = &RawConf{
		raw:       json.RawMessage(message),
		signature: md5.EncryptBytes(message),
		version:   version,
	}
	c.XMap, err = types.NewXMapByJSON(string(message))
	return c, err
}

//GetVersion 获取当前配置的版本号
func (j *RawConf) GetVersion() int32 {
	return j.version
}

//GetSignature 获取当前对象的唯一标识
func (j *RawConf) GetSignature() string {
	return j.signature
}
