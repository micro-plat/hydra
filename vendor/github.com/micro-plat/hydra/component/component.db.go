package component

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/db"
)

//DBTypeNameInVar DB在var配置中的类型名称
const DBTypeNameInVar = "db"

//DBNameInVar DB名称在var配置中的末节点名称
const DBNameInVar = "db"

//IComponentDB Component DB
type IComponentDB interface {
	GetRegularDB(names ...string) (d db.IDB)
	GetDB(names ...string) (d db.IDB, err error)
	GetDBBy(tpName string, name string) (c db.IDB, err error)
	SaveDBObject(tpName string, name string, f func(c conf.IConf) (db.IDB, error)) (bool, db.IDB, error)
	Close() error
}

//StandardDB db
type StandardDB struct {
	IContainer
	name  string
	dbMap cmap.ConcurrentMap
}

//NewStandardDB 创建DB
func NewStandardDB(c IContainer, name ...string) *StandardDB {
	if len(name) > 0 {
		return &StandardDB{IContainer: c, name: name[0], dbMap: cmap.New(2)}
	}
	return &StandardDB{IContainer: c, name: DBNameInVar, dbMap: cmap.New(2)}
}

//GetRegularDB 获取正式的没有异常数据库实例
func (s *StandardDB) GetRegularDB(names ...string) (d db.IDB) {
	d, err := s.GetDB(names...)
	if err != nil {
		panic(err)
	}
	return d
}

//GetDB 获取数据库操作对象
func (s *StandardDB) GetDB(names ...string) (d db.IDB, err error) {
	name := s.name
	if len(names) > 0 {
		name = names[0]
	}
	return s.GetDBBy(DBTypeNameInVar, name)
}

//GetDBBy 根据类型获取缓存数据
func (s *StandardDB) GetDBBy(tpName string, name string) (c db.IDB, err error) {
	_, c, err = s.SaveDBObject(tpName, name, func(jConf conf.IConf) (db.IDB, error) {
		var dbConf conf.DBConf
		if err = jConf.Unmarshal(&dbConf); err != nil {
			return nil, err
		}
		if b, err := govalidator.ValidateStruct(&dbConf); !b {
			return nil, err
		}
		return db.NewDB(dbConf.Provider,
			dbConf.ConnString,
			dbConf.MaxOpen,
			dbConf.MaxIdle,
			dbConf.LefeTime)
	})
	return c, err
}

//SaveDBObject 缓存对象
func (s *StandardDB) SaveDBObject(tpName string, name string, f func(c conf.IConf) (db.IDB, error)) (bool, db.IDB, error) {
	cacheConf, err := s.IContainer.GetVarConf(tpName, name)
	if err != nil {
		return false, nil, fmt.Errorf("%s %v", registry.Join("/", s.GetPlatName(), "var", tpName, name), err)
	}
	key := fmt.Sprintf("%s/%s:%d", tpName, name, cacheConf.GetVersion())
	ok, ch, err := s.dbMap.SetIfAbsentCb(key, func(input ...interface{}) (c interface{}, err error) {
		return f(cacheConf)
	})
	if err != nil {
		err = fmt.Errorf("创建db失败:%s,err:%v", string(cacheConf.GetRaw()), err)
		return ok, nil, err
	}
	return ok, ch.(db.IDB), err
}

//Close 释放所有缓存配置
func (s *StandardDB) Close() error {
	s.dbMap.RemoveIterCb(func(k string, v interface{}) bool {
		v.(*db.DB).Close()
		return true
	})
	return nil
}
