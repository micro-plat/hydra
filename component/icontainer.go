package component

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/cache"
	"github.com/micro-plat/lib4go/db"
	"github.com/micro-plat/lib4go/influxdb"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/queue"
)

type IContainer interface {
	context.RPCInvoker
	GetComponent() IComponent
	conf.ISystemConf
	conf.IVarConf
	conf.IMainConf
	GetRegistry() registry.IRegistry
	GetLooger() logger.ILogging

	GetRegularCache(names ...string) (c cache.ICache)
	GetCache(names ...string) (c cache.ICache, err error)
	GetCacheBy(tpName string, name string) (c cache.ICache, err error)
	SaveCacheObject(tpName string, name string, f func(c conf.IConf) (cache.ICache, error)) (bool, cache.ICache, error)

	GetRegularDB(names ...string) (d db.IDB)
	GetDB(names ...string) (d db.IDB, err error)
	GetDBBy(tpName string, name string) (c db.IDB, err error)
	SaveDBObject(tpName string, name string, f func(c conf.IConf) (db.IDB, error)) (bool, db.IDB, error)

	GetRegularInflux(names ...string) (c influxdb.IInfluxClient)
	GetInflux(names ...string) (d influxdb.IInfluxClient, err error)
	GetInfluxBy(tpName string, name string) (c influxdb.IInfluxClient, err error)
	SaveInfluxObject(tpName string, name string, f func(c conf.IConf) (influxdb.IInfluxClient, error)) (bool, influxdb.IInfluxClient, error)

	GetRegularQueue(names ...string) (c queue.IQueue)
	GetQueue(names ...string) (q queue.IQueue, err error)
	GetQueueBy(tpName string, name string) (c queue.IQueue, err error)
	SaveQueueObject(tpName string, name string, f func(c conf.IConf) (queue.IQueue, error)) (bool, queue.IQueue, error)

	GetGlobalObject(tpName string, name string) (c interface{}, err error)
	SaveGlobalObject(tpName string, name string, f func(c conf.IConf) (interface{}, error)) (bool, interface{}, error)
	Close() error
}
