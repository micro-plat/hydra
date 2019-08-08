package conf

import "fmt"

//DBConf 数据库配置
type DBConf struct {
	Provider   string `json:"provider" valid:"required"`
	ConnString string `json:"connString" valid:"required"`
	MaxOpen    int    `json:"maxOpen" valid:"required"`
	MaxIdle    int    `json:"maxIdle" valid:"required"`
	LifeTime   int    `json:"lifeTime" valid:"required"`
}

//NewMysqlConf 创建mysql数据库
func NewMysqlConf(uName string, pwd string, serverIP string, dbName string) *DBConf {
	return NewDBConf("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", uName, pwd, serverIP, dbName), 20, 10, 600)
}

//NewMysqlConfForProd 创建prod mysql数据库#connectStr
func NewMysqlConfForProd(name ...string) *DBConf {
	kn := "#connectStr"
	if len(name) > 0 {
		kn = name[0]
	}
	return NewDBConf("mysql", kn, 20, 10, 600)
}

//NewOracleConf 创建oracle数据库
func NewOracleConf(uName string, pwd string, tnsName string) *DBConf {
	return NewDBConf("ora", fmt.Sprintf("%s/%s@%s", uName, pwd, tnsName), 200, 100, 600)
}

//NewOracleConfForProd 创建prod oracle数据库#connectStr
func NewOracleConfForProd(name ...string) *DBConf {
	kn := "#connectStr"
	if len(name) > 0 {
		kn = name[0]
	}
	return NewDBConf("ora", kn, 200, 100, 600)
}

//WithConnect 设置连接数与超时时间
func (d *DBConf) WithConnect(maxOpen int, maxIdle int, lifeTime int) *DBConf {
	d.MaxOpen = maxOpen
	d.MaxIdle = maxIdle
	d.LifeTime = lifeTime
	return d
}

//NewDBConf 构建数据库配置对象
func NewDBConf(provider string, connString string, maxOpen int, maxIdle int, lifeTime int) *DBConf {
	return &DBConf{
		Provider:   provider,
		ConnString: connString,
		MaxOpen:    maxOpen,
		MaxIdle:    maxIdle,
		LifeTime:   lifeTime,
	}
}
