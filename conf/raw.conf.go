package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/micro-plat/lib4go/security/md5"
	"github.com/micro-plat/lib4go/types"
)

//IConf 配置管理
type IConf interface {
	GetVersion() int32
	GetString(key string, def ...string) (r string)
	GetStrings(key string, def ...string) (r []string)
	GetInt(key string, def ...int) int
	GetArray(key string, def ...interface{}) (r []interface{})
	GetBool(key string, def ...bool) (r bool)
	GetJSON(section string) (r []byte, version int32, err error)
	GetSection(section string) (c *RawConf, err error)
	HasSection(section string) bool
	GetRaw() []byte
	Unmarshal(v interface{}) error
	GetSignature() string
}

//EmptyJSONConf 空的jsonconf
var EmptyJSONConf = &RawConf{
	raw:       []byte("{}"),
	signature: md5.EncryptBytes([]byte("{}")),
	data:      map[string]interface{}{},
}

//RawConf json配置文件
type RawConf struct {
	raw       json.RawMessage
	version   int32
	signature string
	data      map[string]interface{}
}

//NewRawConfByMap 根据map初始化对象
func NewRawConfByMap(data map[string]interface{}, version int32) (c *RawConf, err error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	c = &RawConf{
		data:      data,
		version:   version,
		raw:       buf,
		signature: md5.EncryptBytes(buf),
	}
	return c, nil
}

//NewRawConfByJson 初始化JsonConf
func NewRawConfByJson(message []byte, version int32) (c *RawConf, err error) {
	c = &RawConf{
		raw:       json.RawMessage(message),
		signature: md5.EncryptBytes(message),
		version:   version,
	}

	if bytes.HasPrefix(message, []byte("{")) && bytes.HasSuffix(message, []byte("}")) {
		if err = json.Unmarshal(message, &c.data); err != nil {
			return nil, err
		}
	}

	return c, nil
}

//Unmarshal 将当前[]byte反序列化为对象
func (j *RawConf) Unmarshal(v interface{}) error {
	return json.Unmarshal(j.raw, v)
}

//GetVersion 获取当前配置的版本号
func (j *RawConf) GetVersion() int32 {
	return j.version
}

//GetRaw 获取json数据
func (j *RawConf) GetRaw() []byte {
	if len(j.raw) > 0 {
		return j.raw
	}
	j.raw, _ = json.Marshal(j.data)
	return j.raw
}

//GetString 获取字符串
func (j *RawConf) GetString(key string, def ...string) (r string) {
	if val, ok := j.data[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		case map[string]interface{}:
			buffer, _ := json.Marshal(val)
			return string(buffer)
		default:
			return fmt.Sprint(val)
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

//GetInt 获取数字
func (j *RawConf) GetInt(key string, def ...int) int {
	if v, err := strconv.Atoi(j.GetString(key)); err == nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//GetStrings 获取字符串数组
func (j *RawConf) GetStrings(key string, def ...string) (r []string) {
	if v := j.GetString(key); v != "" {
		if r = strings.Split(v, ";"); len(r) > 0 {
			return r
		}
	}
	if len(def) > 0 {
		return def
	}
	return nil
}

//GetArray 获取数组对象
func (j *RawConf) GetArray(key string, def ...interface{}) (r []interface{}) {
	if _, ok := j.data[key]; !ok {
		if len(def) > 0 {
			return def
		}
		return nil
	}
	if r, ok := j.data[key].([]interface{}); ok {
		return r
	}
	return nil
}

//GetBool 获取bool类型值
func (j *RawConf) GetBool(key string, def ...bool) (r bool) {
	if val := j.GetString(key); val != "" {
		if b, err := types.ParseBool(val); err == nil {
			return b
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

//GetJSON 获取section原始JSON串
func (j *RawConf) GetJSON(section string) (r []byte, version int32, err error) {
	if v, ok := j.data[section]; !ok || v == nil {
		err = fmt.Errorf("节点:%s不存在或值为空", section)
		return
	}
	val := j.data[section]
	buffer, err := json.Marshal(val)
	if err != nil {
		return nil, 0, err
	}
	return buffer, j.version, nil
}

//HasSection 是否存在节点
func (j *RawConf) HasSection(section string) bool {
	_, ok := j.data[section].(map[string]interface{})
	return ok
}

//GetSection 指定节点名称获取JSONConf
func (j *RawConf) GetSection(section string) (c *RawConf, err error) {
	if v, ok := j.data[section]; !ok || v == nil {
		err = fmt.Errorf("节点:%s不存在或值为空", section)
		return
	}
	if data, ok := j.data[section].(map[string]interface{}); ok {
		return NewRawConfByMap(data, j.version)
	}
	err = fmt.Errorf("节点:%s不是有效的json对象", section)
	return
}

//GetSignature 获取当前对象的唯一标识
func (j *RawConf) GetSignature() string {
	return j.signature
}
