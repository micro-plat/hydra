package dbs

import "github.com/micro-plat/lib4go/db"

//IDB 数据库接口
type IDB = db.IDB

//IComponentDB Component DB
type IComponentDB interface {
	GetRegularDB(names ...string) (d IDB)
	GetDB(names ...string) (d IDB, err error)
}
