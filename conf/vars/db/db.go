package db

import "github.com/micro-plat/hydra/conf"

//DB 数据库配置
type DB struct {
	Provider   string `json:"provider" valid:"required"`
	ConnString string `json:"connString" valid:"required"`
	MaxOpen    int    `json:"maxOpen" valid:"required"`
	MaxIdle    int    `json:"maxIdle" valid:"required"`
	LifeTime   int    `json:"lifeTime" valid:"required"`
}

//New 构建DB连接信息
func New(provider string, connString string, opts ...Option) *DB {
	db := &DB{
		Provider:   provider,
		ConnString: connString,
		MaxOpen:    10,
		MaxIdle:    3,
		LifeTime:   600,
	}
	for _, opt := range opts {
		opt(db)
	}
	return db
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IVarConf, tp string, name string) (s *DB, err error) {
	if _, err = cnf.GetObject(tp, name, &s); err != nil && err != conf.ErrNoSetting {
		return nil, err
	}
	return s, nil
}
