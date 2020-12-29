package ras

import (
	"encoding/json"

	"github.com/micro-plat/hydra/conf"
)

//Auth 远程认证服务
type Auth struct {

	//远程验证服务名
	Service string `json:"service,omitempty" valid:"required" toml:"service,omitempty" label:"远程验证服务名"`

	//指定需要远程验证的请求列表
	Requests []string `json:"requests,omitempty" valid:"required" toml:"requests,omitempty" label:"远程验证的请求列表"`

	//必须传入的字段列表
	Required []string `json:"required,omitempty" toml:"required,omitempty"`

	//远程验证的字段连接方式
	Connect *Connect `json:"connect,omitempty" toml:"connect,omitempty"`

	//字段别名表
	Alias map[string]string `json:"alias,omitempty" toml:"alias,omitempty"`

	//扩展参数表
	Params map[string]interface{} `json:"params,omitempty" toml:"params,omitempty"`

	//需要解密的字段
	Decrypt []string `json:"decrypt,omitempty" toml:"decrypt,omitempty"`

	//是否需要检查时间戳
	CheckTS bool `json:"checkTimestamp,omitempty" toml:"checkTimestamp,omitempty"`

	//配置是否禁用
	Disable bool `json:"disable,omitempty" toml:"disable,omitempty"`

	*conf.PathMatch `json:"-" toml:"-"`
}

//New 创建远程服务验证参数
func New(service string, opts ...AuthOption) *Auth {
	r := &Auth{
		Service:  service,
		Requests: []string{},
		Connect:  &Connect{},
		Params:   make(map[string]interface{}),
		Required: make([]string, 0, 1),
		Alias:    make(map[string]string),
		Decrypt:  make([]string, 0, 1),
	}
	for _, opt := range opts {
		opt(r)
	}
	if len(r.Requests) <= 0 {
		//默认通配 所有路径都要验证
		r.Requests = []string{"/**"}
	}
	r.PathMatch = conf.NewPathMatch(r.Requests...)
	return r
}

//String 获取签名串
func (a *Auth) String() (string, error) {
	buff, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

//AuthString 获取签名串
func (a *Auth) AuthString() (string, error) {
	b := *a
	b.Service = ""
	b.Requests = nil
	c := &b
	return c.String()
}
