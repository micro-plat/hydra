package conf

import (
	"encoding/json"

	"github.com/micro-plat/lib4go/types"
)

//ServiceAuth 服务认证配置
type ServiceAuth struct {
	Service  string                 `json:"service,omitempty" valid:"required"`
	Requests []string               `json:"requests,omitempty" valid:"required"`
	Required []string               `json:"required,omitempty"`
	Connect  *Connect               `json:"connect,omitempty"`
	Alias    map[string]string      `json:"alias,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
	Decrypt  []string               `json:"decrypt,omitempty"`
	CheckTS  bool                   `json:"check-timestamp,omitempty"`
	Disable  bool                   `json:"disable,omitempty"`
}

//WithServiceAuth 添加远程服务验证
func (a *Authes) WithServiceAuth(auth ...*ServiceAuth) *Authes {
	a.RemotingServiceAuths = append(a.RemotingServiceAuths, auth...)
	return a
}

//NewServiceAuth 创建远程服务验证参数
func NewServiceAuth(service string, request ...string) *ServiceAuth {
	requests := request
	if len(requests) == 0 {
		requests = []string{"*"}
	}
	return &ServiceAuth{
		Service:  service,
		Requests: requests,
		Connect:  &Connect{},
		Params:   make(map[string]interface{}),
		Required: make([]string, 0, 1),
		Alias:    make(map[string]string),
		Decrypt:  make([]string, 0, 1),
	}
}

//String 获取签名串
func (a *ServiceAuth) String() (string, error) {
	buff, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

//AuthString 获取签名串
func (a *ServiceAuth) AuthString() (string, error) {
	b := *a
	b.Service = ""
	b.Requests = nil
	c := &b
	return c.String()
}

//WithRequest 设置requests的请求服务路径
func (a *ServiceAuth) WithRequest(path ...string) *ServiceAuth {
	if len(path) > 0 {
		a.Requests = append(a.Requests, path...)
	}
	return a
}

//WithRequired 设置必须字段
func (a *ServiceAuth) WithRequired(fieldName ...string) *ServiceAuth {
	if len(fieldName) > 0 {
		a.Required = append(a.Required, fieldName...)
	}
	return a
}

//WithUIDAlias 设置用户id的字段名
func (a *ServiceAuth) WithUIDAlias(name string) *ServiceAuth {
	a.Alias["euid"] = name
	return a
}

//WithTimestampAlias 设置timestamp的字段名
func (a *ServiceAuth) WithTimestampAlias(name string) *ServiceAuth {
	a.Alias["timestamp"] = name
	return a
}

//WithSignAlias 设置sign的字段名
func (a *ServiceAuth) WithSignAlias(name string) *ServiceAuth {
	a.Alias["sign"] = name
	return a
}

//WithDecryptName 设置需要解密的字段名
func (a *ServiceAuth) WithDecryptName(name ...string) *ServiceAuth {
	a.Decrypt = append(a.Decrypt, name...)
	return a
}

//WithCheckTimestamp 设置需要检查时间戳
func (a *ServiceAuth) WithCheckTimestamp(e ...bool) *ServiceAuth {
	a.CheckTS = types.GetBoolByIndex(e, 0, true)
	return a
}

//WithParam 设置扩展参数
func (a *ServiceAuth) WithParam(key string, value interface{}) *ServiceAuth {
	a.Params[key] = value
	return a
}

// WithConnect 设置签名链接方式
func (a *ServiceAuth) WithConnect() *Connect {
	a.Connect = &Connect{auth: a}
	return a.Connect
}

//WithDisable 禁用配置
func (a *ServiceAuth) WithDisable() *ServiceAuth {
	a.Disable = true
	return a
}

//WithEnable 启用配置
func (a *ServiceAuth) WithEnable() *ServiceAuth {
	a.Disable = false
	return a
}

//ServiceAuths 远程服务验证组
type ServiceAuths []*ServiceAuth

//Contains 检查指定的路径是否允许签名
func (a ServiceAuths) Contains(p string) (bool, *ServiceAuth) {
	var last *ServiceAuth
	for _, auth := range a {
		for _, req := range auth.Requests {
			if req == p {
				return true, auth
			}
			if req == "*" {
				last = auth
			}
		}
	}
	return last != nil, last
}
