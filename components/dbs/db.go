package dbs

import (
	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/db"
	"github.com/micro-plat/lib4go/types"
)

const (
	//typeNode DB在var配置中的类型名称
	dbTypeNode = "db"

	//nameNode DB名称在var配置中的末节点名称
	dbNameNode = "db"
)

//IDB 数据库接口
type IDB = db.IDB

//IComponentDB Component DB
type IComponentDB interface {
	GetRegularDB(names ...string) (d IDB)
	GetDB(names ...string) (d IDB, err error)
}

//StandardDB db
type StandardDB struct {
	c components.IComponents
}

//NewStandardDB 创建DB
func NewStandardDB(c components.IComponents) *StandardDB {
	return &StandardDB{c: c}
}

//GetRegularDB 获取正式的没有异常数据库实例
func (s *StandardDB) GetRegularDB(names ...string) (d IDB) {
	d, err := s.GetDB(names...)
	if err != nil {
		panic(err)
	}
	return d
}

//GetDB 获取数据库操作对象
func (s *StandardDB) GetDB(names ...string) (d IDB, err error) {
	name := types.GetStringByIndex(names, 0, dbNameNode)
	obj, err := s.c.GetOrCreate(dbTypeNode, name, func(c conf.IConf) (interface{}, error) {
		var dbConf conf.DBConf
		if err = c.Unmarshal(&dbConf); err != nil {
			return nil, err
		}
		if b, err := govalidator.ValidateStruct(&dbConf); !b {
			return nil, err
		}
		return db.NewDB(dbConf.Provider,
			dbConf.ConnString,
			dbConf.MaxOpen,
			dbConf.MaxIdle,
			dbConf.LifeTime)
	})
	if err != nil {
		return nil, err
	}
	return obj.(IDB), nil
}
