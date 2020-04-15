package oracle

import "fmt"

//Oracle oracle数据库连接信息
type Oracle struct {
	Provider   string `json:"provider" valid:"required"`
	ConnString string `json:"connString" valid:"required"`
	*option
}

//New 构建oracle连接信息
func New(connString string, opts ...Option) *Oracle {
	ora := &Oracle{
		Provider:   "oracle",
		ConnString: connString,
		option: &option{
			MaxOpen:  10,
			MaxIdle:  3,
			LifeTime: 600,
		},
	}
	for _, opt := range opts {
		opt(ora.option)
	}
	return ora

}

//NewBy 构建oracle连接信息
func NewBy(uName string, pwd string, tnsName string, opts ...Option) *Oracle {
	ora := &Oracle{
		Provider:   "oracle",
		ConnString: fmt.Sprintf("%s/%s@%s", uName, pwd, tnsName),
		option: &option{
			MaxOpen:  10,
			MaxIdle:  3,
			LifeTime: 600,
		},
	}
	for _, opt := range opts {
		opt(ora.option)
	}
	return ora

}
