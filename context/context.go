package context

import (
	"strings"
	"sync"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/cache"
	"github.com/micro-plat/lib4go/db"
	"github.com/micro-plat/lib4go/influxdb"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/queue"
)

type IContainer interface {
	RPCInvoker

	conf.ISystemConf
	conf.IVarConf
	conf.IMainConf
	GetRegistry() registry.IRegistry

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

//Context 引擎执行上下文
type Context struct {
	Request   *Request
	Response  *Response
	RPC       *ContextRPC
	container IContainer
	Component interface{}
	Meta      *Meta
	Log       logger.ILogger
	Name      string
	Engine    string
	Service   string
}

//GetContext 从缓存池中获取一个context
func GetContext(component interface{}, name string, engine string, service string, container IContainer, queryString IData, form IData, param IData, setting IData, ext map[string]interface{}, logger *logger.Logger) *Context {
	c := contextPool.Get().(*Context)
	c.Request.reset(c, queryString, form, param, setting, ext)
	c.Log = logger
	c.Component = component
	c.container = container
	c.Name = name
	c.Engine = engine
	c.Service = formatName(c.Request.Translate(service, false))
	c.Meta = NewMeta()
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
	c.RPC.clear()
	c.container = nil
	contextPool.Put(c)
}
func formatName(name string) string {
	text := "/" + strings.Trim(strings.Trim(name, " "), "/")
	return strings.ToLower(text)
}
