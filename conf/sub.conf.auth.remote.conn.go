package conf

import "strings"

//Connect 签名拼接串
type Connect struct {
	KeyValue string         `json:"kv,omitempty" valid:"ascii"`
	Chain    string         `json:"chain,omitempty" valid:"ascii"`
	Sort     string         `json:"sort,omitempty" valid:"in(all|data|static)"`
	Fields   string         `json:"fields,omitempty" valid:"ascii"`
	Secret   *SecretConnect `json:"secret,omitempty"`
	auth     *ServiceAuth   `json:"-" valid:"-"`
}

//SecretConnect secret拼接串
type SecretConnect struct {
	Name     string   `json:"name,omitempty" valid:"ascii"`
	KeyValue string   `json:"kv,omitempty" valid:"ascii"`
	Chain    string   `json:"chain,omitempty" valid:"ascii"`
	Mode     string   `json:"mode,omitempty" valid:"in(head|tail|headTail)"`
	connect  *Connect `json:"-" valid:"-"`
}

//Set 设置字段拼接方式
func (c *Connect) Set(kv string, chain string) *Connect {
	c.KeyValue = kv
	c.Chain = chain
	return c
}

//SortByData 只排序数据字段，不排序secrect
func (c *Connect) SortByData() *Connect {
	c.Sort = "data"
	return c
}

//SortAll 排序所有字段，包括数据，secrect
func (c *Connect) SortAll() *Connect {
	c.Sort = "all"
	return c
}

//SortStatic 使用指定的字段进行排序
func (c *Connect) SortStatic(fields ...string) *Connect {
	c.Sort = "static"
	c.Fields = strings.Join(fields, "|")
	return c
}

//Auth 返回Auth对象
func (c *Connect) Auth() *ServiceAuth {
	return c.auth
}

//SetSecretConnect 设置secrect拼接方式
func (c *Connect) SetSecretConnect() *SecretConnect {
	c.Secret = &SecretConnect{connect: c}
	return c.Secret
}

//SetName 设置secrect的键名称
func (c *SecretConnect) SetName(name string, kv string) *SecretConnect {
	c.Name = name
	c.KeyValue = kv
	return c
}

//SetChainWithHead 设置secrect与数据串之间的拼接方式,并将secret串拼接到数据串的头部
func (c *SecretConnect) SetChainWithHead(chain string) *SecretConnect {
	c.Chain = chain
	c.Mode = "head"
	return c
}

//SetChainWithTail 设置secrect与数据串之间的拼接方式，并将secret串拼接到数据串的尾部
func (c *SecretConnect) SetChainWithTail(chain string) *SecretConnect {
	c.Chain = chain
	c.Mode = "tail"
	return c
}

//SetChainWithHeadAndTail 设置secrect与数据串之间的拼接方式，并将secret串拼接到数据串的头部和尾部
func (c *SecretConnect) SetChainWithHeadAndTail(chain string) *SecretConnect {
	c.Chain = chain
	c.Mode = "headTail"
	return c
}

//Connect 获取父级拼接串
func (c *SecretConnect) Connect() *Connect {
	return c.connect
}
