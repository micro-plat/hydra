package db

import "github.com/micro-plat/hydra/conf"

//DB 数据库配置
type DB struct {
	Provider   string `json:"provider" valid:"required"`
	ConnString string `json:"connString" valid:"required"`
	*option
}

//New 构建DB连接信息
func New(provider string, connString string, opts ...Option) *DB {
	ora := &DB{
		Provider:   provider,
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

//GetConf 获取主配置信息
func GetConf(cnf conf.IVarConf, tp string, name string) (s *DB, err error) {
	if _, err = cnf.GetObject(tp, name, &s); err != nil && err != conf.ErrNoSetting {
		return nil, err
	}
	return s, nil
}
