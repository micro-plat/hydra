package context

import (
	"sync"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/qxnw/lib4go/cache"
	"github.com/qxnw/lib4go/db"
	"github.com/qxnw/lib4go/influxdb"
	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lib4go/queue"
)

type IContainer interface {
	RPCInvoker

	conf.ISystemConf
	conf.IVarConf

	GetRegistry() registry.IRegistry
	GetCache(names ...string) (c cache.ICache, err error)
	GetCacheBy(tpName string, name string) (c cache.ICache, err error)
	SaveCacheObject(tpName string, name string, f func(c conf.IConf) (cache.ICache, error)) (bool, cache.ICache, error)

	GetDB(names ...string) (d db.IDB, err error)
	GetDBBy(tpName string, name string) (c db.IDB, err error)
	SaveDBObject(tpName string, name string, f func(c conf.IConf) (db.IDB, error)) (bool, db.IDB, error)

	GetInflux(names ...string) (d influxdb.IInfluxClient, err error)
	GetInfluxBy(tpName string, name string) (c influxdb.IInfluxClient, err error)
	SaveInfluxObject(tpName string, name string, f func(c conf.IConf) (influxdb.IInfluxClient, error)) (bool, influxdb.IInfluxClient, error)

	GetQueue(names ...string) (q queue.IQueue, err error)
	GetQueueBy(tpName string, name string) (c queue.IQueue, err error)
	SaveQueueObject(tpName string, name string, f func(c conf.IConf) (queue.IQueue, error)) (bool, queue.IQueue, error)

	GetGlobalObject(tpName string, name string) (c interface{}, err error)
	SaveGlobalObject(tpName string, name string, f func(c conf.IConf) (interface{}, error)) (bool, interface{}, error)
	Close() error
}

//Context 引擎执行上下文
type Context struct {
	Request   *Request
	Response  *Response
	RPC       *ContextRPC
	container IContainer
	Log       logger.ILogger
}

//GetContext 从缓存池中获取一个context
func GetContext(container IContainer, queryString IData, form IData, param IData, setting IData, ext map[string]interface{}, logger *logger.Logger) *Context {
	c := contextPool.Get().(*Context)
	c.Request.reset(queryString, form, param, setting, ext)
	c.Log = logger
	c.container = container
	return c
}

//GetContainer 获取当前容器
func (c *Context) GetContainer() IContainer {
	return c.container
}

//SetRPC 根据输入的context创建插件的上下文对象
func (c *Context) SetRPC(rpc RPCInvoker) {
	c.RPC.reset(c, rpc)
}

var contextPool *sync.Pool

func init() {
	contextPool = &sync.Pool{
		New: func() interface{} {
			return &Context{
				RPC:      &ContextRPC{},
				Request:  newRequest(),
				Response: NewResponse(),
			}
		},
	}
}

//Close 回收context
func (c *Context) Close() {
	c.Request.clear()
	c.Response.clear()
	c.container = nil
	contextPool.Put(c)
}
