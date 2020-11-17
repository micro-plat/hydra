package ras

import "strings"

const (
	//FiledSortAll 字段排序模式,排序所有字段，包括数据，secrect
	FiledSortAll = "all"
	//FiledSortData 密钥拼接模式,只排序数据字段，不排序secrect
	FiledSortData = "data"
	//FiledSortStatic 密钥拼接模式,使用指定的字段进行排序
	FiledSortStatic = "static"
)

//Connect 签名拼接串
type Connect struct {

	//键值连接符
	KeyValue string `json:"kv,omitempty" valid:"ascii" toml:"kv,omitempty"`

	//每组键值连接符
	Chain string `json:"chain,omitempty" valid:"ascii" toml:"chain,omitempty"`

	//排序方式
	Sort string `json:"sort,omitempty" valid:"in(all|data|static)" toml:"sort,omitempty"`

	//参与签名验证的字段
	Fields string `json:"fields,omitempty" valid:"ascii" toml:"fields,omitempty"`

	//密钥连接方式
	Secret *SecretConnect `json:"secret,omitempty" toml:"secret,omitempty"`
}

//ConnectOption 配置选项
type ConnectOption func(*Connect)

//WithConnectChar 设置字段拼接方式
func WithConnectChar(kv string, chain string) ConnectOption {
	return func(c *Connect) {
		c.KeyValue = kv
		c.Chain = chain
	}
}

//WithConnectSortByData 只排序数据字段，不排序secrect
func WithConnectSortByData() ConnectOption {
	return func(c *Connect) {
		c.Sort = FiledSortData
	}
}

//WithConnectSortAll 排序所有字段，包括数据，secrect
func WithConnectSortAll() ConnectOption {
	return func(c *Connect) {
		c.Sort = FiledSortAll
	}
}

//WithConnectSortStatic 使用指定的字段进行排序
func WithConnectSortStatic(fields ...string) ConnectOption {
	return func(c *Connect) {
		c.Sort = FiledSortStatic
		c.Fields = strings.Join(fields, "|")
	}
}

//WithSecretConnect 启用配置
func WithSecretConnect(opts ...SecretOption) ConnectOption {
	return func(a *Connect) {
		a.Secret = &SecretConnect{}
		for _, opt := range opts {
			opt(a.Secret)
		}
	}
}
