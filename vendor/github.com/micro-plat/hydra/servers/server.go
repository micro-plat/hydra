package servers

import (
	"fmt"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/logger"
)

var IsDebug = false
var (
	ST_RUNNING   = "running"
	ST_STOP      = "stop"
	ST_PAUSE     = "pause"
	SRV_TP_API   = "api"
	SRV_FILE_API = "file"
	SRV_TP_RPC   = "rpc"
	SRV_TP_CRON  = "cron"
	SRV_TP_MQ    = "mq"
	SRV_TP_WEB   = "web"
)

//IRegistryServer 基于注册中心的服务器
type IRegistryServer interface {
	Notify(conf.IServerConf) error
	Start() error
	GetAddress() string
	GetServices() []string
	Restarted() bool
	GetStatus() string
	Shutdown()
}

type IExecuter interface {
	Execute(ctx *context.Context) (rs interface{})
}

type IExecuteHandler func(ctx *context.Context) (rs interface{})

func (i IExecuteHandler) Execute(ctx *context.Context) (rs interface{}) {
	return i(ctx)
}

//IRegistryEngine 基于注册中心的执行引擎
type IRegistryEngine interface {
	context.IContainer
	IExecuter
	GetComponent() component.IComponent
	SetHandler(h component.IComponentHandler) error
	UpdateVarConf(conf conf.IServerConf)
	GetServices() []string
	Fallback(c *context.Context) (rs interface{})
}

//IServerResolver 服务器生成器
type IServerResolver interface {
	Resolve(registryAddr string, conf conf.IServerConf, log *logger.Logger) (IRegistryServer, error)
}
type IServerResolverHandler func(registryAddr string, conf conf.IServerConf, log *logger.Logger) (IRegistryServer, error)

//Resolve 创建服务器实例
func (i IServerResolverHandler) Resolve(registryAddr string, conf conf.IServerConf, log *logger.Logger) (IRegistryServer, error) {
	return i(registryAddr, conf, log)
}

var resolvers = make(map[string]IServerResolver)

//Register 注册服务器生成器
func Register(identifier string, resolver IServerResolver) {
	if _, ok := resolvers[identifier]; ok {
		panic("server: Register called twice for identifier: " + identifier)
	}
	resolvers[identifier] = resolver
}

//NewRegistryServer 根据服务标识创建服务器
func NewRegistryServer(identifier string, registryAddr string, conf conf.IServerConf, log *logger.Logger) (IRegistryServer, error) {
	if resolver, ok := resolvers[identifier]; ok {
		return resolver.Resolve(registryAddr, conf, log)
	}
	return nil, fmt.Errorf("server: unknown identifier name %q (forgotten import?)", identifier)
}

//Trace 打印跟踪信息
func Trace(print func(f string, args ...interface{}), args ...interface{}) {
	// if !IsDebug {
	// 	return
	// }
	print("%s", args)
}

//Tracef 根据格式打印跟踪信息
func Tracef(print func(f string, args ...interface{}), format string, args ...interface{}) {
	// if !IsDebug {
	// 	return
	// }
	print(format, args...)
}

//TraceIf 根据条件打印跟踪信息
func TraceIf(b bool, okPrint func(f string, args ...interface{}), print func(f string, args ...interface{}), args ...interface{}) {
	// if !IsDebug {
	// 	return
	// }
	if b {
		okPrint("%s", args)
		return
	}
	print("%s", args)
}
