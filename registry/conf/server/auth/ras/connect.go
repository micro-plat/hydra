package ras

//Connect 签名拼接串
type Connect struct {
	*connectOption
}

//SecretConnect secret拼接串
type SecretConnect struct {

	//为密钥指定的键名
	Name string `json:"name,omitempty" valid:"ascii"`

	//密钥键与密钥值连接符
	KeyValue string `json:"kv,omitempty" valid:"ascii"`

	//密钥与其它串的连接符
	Chain string `json:"chain,omitempty" valid:"ascii"`

	//密钥连接方式
	Mode string `json:"mode,omitempty" valid:"in(head|tail|headTail)"`
}

//WithSecretName 设置secrect的键名称
func (c *SecretConnect) WithSecretName(name string, kv string) *SecretConnect {
	c.Name = name
	c.KeyValue = kv
	return c
}

//WithSecretHeadMode 设置secrect与数据串之间的拼接方式,并将secret串拼接到数据串的头部
func (c *SecretConnect) WithSecretHeadMode(chain string) *SecretConnect {
	c.Chain = chain
	c.Mode = "head"
	return c
}

//WithSecretTailMode 设置secrect与数据串之间的拼接方式，并将secret串拼接到数据串的尾部
func (c *SecretConnect) WithSecretTailMode(chain string) *SecretConnect {
	c.Chain = chain
	c.Mode = "tail"
	return c
}

//WithSecretHeadAndTailMode 设置secrect与数据串之间的拼接方式，并将secret串拼接到数据串的头部和尾部
func (c *SecretConnect) WithSecretHeadAndTailMode(chain string) *SecretConnect {
	c.Chain = chain
	c.Mode = "headTail"
	return c
}
