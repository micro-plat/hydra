package dbs

import (
	"fmt"

	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/lib4go/db"
	"github.com/micro-plat/lib4go/types"

	"github.com/micro-plat/hydra/conf"
	xdb "github.com/micro-plat/hydra/conf/vars/db"
)

const (
	//typeNode DB在var配置中的类型名称
	dbTypeNode = "db"

	//nameNode DB名称在var配置中的末节点名称
	dbNameNode = "db"
)

//StandardDB db
type StandardDB struct {
	c container.IContainer
}

//NewStandardDB 创建DB
func NewStandardDB(c container.IContainer) *StandardDB {
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
	obj, err := s.c.GetOrCreate(dbTypeNode, name, func(conf *conf.RawConf, keys ...string) (obj interface{}, err error) {
		if conf.IsEmpty() {
			return nil, fmt.Errorf("节点/%s/%s未配置，或不可用", dbTypeNode, name)
		}
		var dbConf xdb.DB
		if err = conf.ToStruct(&dbConf); err != nil {
			return nil, fmt.Errorf("数据库[%s/%s]配置有误：%w", dbTypeNode, name, err)
		}
		return db.NewDB(dbConf.Provider, dbConf.ConnString, dbConf.MaxOpen, dbConf.MaxIdle, dbConf.LifeTime)
	})
	if err != nil {
		return nil, err
	}
	return obj.(IDB), nil
}
