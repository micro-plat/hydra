package ras

import "strings"

type connectOption struct {

	//键值连接符
	KeyValue string `json:"kv,omitempty" valid:"ascii"`

	//每组键值连接符
	Chain string `json:"chain,omitempty" valid:"ascii"`

	//排序方式
	Sort string `json:"sort,omitempty" valid:"in(all|data|static)"`

	//参与签名验证的字段
	Fields string `json:"fields,omitempty" valid:"ascii"`

	//密钥连接方式
	Secret *SecretConnect `json:"secret,omitempty"`
}

//ConnectOption 配置选项
type ConnectOption func(*connectOption)

//WithConnectChar 设置字段拼接方式
func WithConnectChar(kv string, chain string) ConnectOption {
	return func(c *connectOption) {
		c.KeyValue = kv
		c.Chain = chain
	}
}

//WithConnectSortByData 只排序数据字段，不排序secrect
func WithConnectSortByData() ConnectOption {
	return func(c *connectOption) {
		c.Sort = "data"
	}
}

//WithConnectSortAll 排序所有字段，包括数据，secrect
func WithConnectSortAll() ConnectOption {
	return func(c *connectOption) {
		c.Sort = "all"
	}
}

//WithConnectSortStatic 使用指定的字段进行排序
func WithConnectSortStatic(fields ...string) ConnectOption {
	return func(c *connectOption) {
		c.Sort = "static"
		c.Fields = strings.Join(fields, "|")
	}
}
