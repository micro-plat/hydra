package mysql

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
		Provider:   "mysql",
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
func NewBy(uName string, pwd string, serverIP string, dbName string, opts ...Option) *Oracle {
	ora := &Oracle{
		Provider:   "mysql",
		ConnString: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", uName, pwd, serverIP, dbName),
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
