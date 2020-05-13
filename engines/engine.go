package engines

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/micro-plat/lib4go/concurrent/cmap"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/rpc"
	"github.com/micro-plat/lib4go/logger"
)

//IServiceEngine 服务引擎接口
type IServiceEngine interface {
	GetRegistry() registry.IRegistry
	GetServices() []string
	Fallback(c *context.Context) (rs interface{})
	Execute(ctx *context.Context) (rs interface{})
	Close() error
}

//ServiceEngine 服务引擎
type ServiceEngine struct {
	*component.StandardComponent
	cHandler component.IComponentHandler
	loggers  cmap.ConcurrentMap
	conf.IServerConf
	registryAddr string
	*rpc.Invoker
	logger   logger.ILogging
	registry registry.IRegistry
	component.IComponentCache
	component.IComponentDB
	component.IComponentInfluxDB
	component.IComponentQueue
	component.IComponentGlobalVarObject
}

//NewServiceEngine 构建服务引擎
func NewServiceEngine(conf conf.IServerConf, registryAddr string, log logger.ILogging) (e *ServiceEngine, err error) {
	e = &ServiceEngine{IServerConf: conf, registryAddr: registryAddr, logger: log}
	e.StandardComponent = component.NewStandardComponent("sys.engine", e)
	e.IComponentCache = component.NewStandardCache(e, "cache")
	e.IComponentDB = component.NewStandardDB(e, "db")
	e.IComponentInfluxDB = component.NewStandardInfluxDB(e, "influx")
	e.IComponentQueue = component.NewStandardQueue(e, "queue")
	e.IComponentGlobalVarObject = component.NewGlobalVarObjectCache(e)
	e.loggers = cmap.New(8)
	if e.registry, err = registry.NewRegistryWithAddress(registryAddr, log); err != nil {
		return
	}

	if err = e.loadEngineServices(); err != nil {
		return nil, err
	}
	e.StandardComponent.AddRPCProxy(e.RPCProxy())
	err = e.StandardComponent.LoadServices()
	return
}

//SetHandler 设置handler
func (r *ServiceEngine) SetHandler(h component.IComponentHandler) error {
	if h == nil {
		return nil
	}

	//初始化服务器
	r.cHandler = h
	funcs := h.GetInitializings()
	for _, f := range funcs {
		if err := f(r); err != nil {
			return err
		}
	}

	//初始化RPC调用证书
	tls := h.GetRPCTLS()
	opts := make([]rpc.InvokerOption, 0, 0)
	for k, v := range tls {
		opts = append(opts, rpc.WithRPCTLS(k, v))
	}

	//设置负载均衡器
	balancers := h.GetBalancer()
	for v, p := range balancers {
		opts = append(opts, rpc.WithBalancerMode(v, p.Mode, p.Param))
	}

	r.Invoker = rpc.NewInvoker(
		r.IServerConf.GetPlatName(),
		r.IServerConf.GetSysName(),
		r.registryAddr,
		opts...)

	//初始化服务注册
	svs := h.GetServices()

	for group, handlers := range svs {
		for name, handler := range handlers {
			r.StandardComponent.AddCustomerService(name, handler, group, h.GetTags(name)...)
		}
	}
	err := r.StandardComponent.LoadServices()
	return err
}

//UpdateVarConf 更新var配置参数
func (r *ServiceEngine) UpdateVarConf(conf conf.IServerConf) {
	r.SetVarConf(conf.GetVarConfClone())
	r.SetSubConf(conf.GetSubConfClone())
}

//GetServices 获取组件提供的所有服务
func (r *ServiceEngine) GetServices() map[string][]string {
	return r.GetRegistryNames(component.GetGroupName(r.GetServerType())...)

}

//GetTags 添加获和取tag接口
func (r *ServiceEngine) GetTags(name string) []string {
	return r.StandardComponent.GetTags(name)
}

//Execute 执行外部请求
func (r *ServiceEngine) Execute(ctx *context.Context) (rs interface{}) {
	id := fmt.Sprint(goid())
	r.loggers.Set(id, ctx.Log)
	defer r.loggers.Remove(id)
	if ctx.Request.CircuitBreaker.IsOpen() { //熔断开关打开，则自动降级
		rf := r.StandardComponent.Fallback(ctx)
		if r, ok := rf.(error); ok && r == component.ErrNotFoundFallbackService {
			ctx.Response.MustContent(ctx.Request.CircuitBreaker.GetDefStatus(),
				ctx.Request.CircuitBreaker.GetDefContent())
		}
		return rf
	}
	if strings.ToUpper(ctx.Request.GetMethod()) == "OPTIONS" {
		return
	}
	//当前引擎预处理
	if rs = r.Handling(ctx); ctx.Response.HasError(rs) {
		return rs
	}

	//当前服务器预处理
	if r.cHandler != nil && r.cHandler.GetHandlings() != nil {
		hds := r.cHandler.GetHandlings()
		for _, h := range hds {
			if rs = h(ctx); ctx.Response.HasError(rs) {
				return rs
			}
		}
	}
	if !ctx.Response.SkipHandle {
		//当前服务处理
		rs = r.Handle(ctx)
	}

	//当前服务器后处理
	if r.cHandler != nil && len(r.cHandler.GetHandleds()) > 0 {
		hdd := r.cHandler.GetHandleds()
		return hdd[0](ctx)
	}
	return rs
}

//Handling 每次handle执行前执行
func (r *ServiceEngine) Handling(ctx *context.Context) (rs interface{}) {
	ctx.SetRPC(r.Invoker)
	err := checkSignByFixedSecret(ctx)
	if err != nil {
		return err
	}
	return checkSignByRemoteSecret(ctx)
}

//GetRegistry 获取注册中心
func (r *ServiceEngine) GetRegistry() registry.IRegistry {
	return r.registry
}

//Close 关闭引擎
func (r *ServiceEngine) Close() error {
	if r.cHandler != nil {
		funcs := r.cHandler.GetClosings()
		for _, f := range funcs {
			if err := f(r); err != nil {
				r.logger.Error(err)
			}
		}
	}
	r.StandardComponent.Close()
	if r.Invoker != nil {
		r.Invoker.Close()
	}
	r.IComponentGlobalVarObject.Close()
	r.IComponentCache.Close()
	r.IComponentInfluxDB.Close()
	r.IComponentQueue.Close()
	r.IComponentDB.Close()
	r.loggers.Clear()
	return nil
}
func (r *ServiceEngine) GetLogger() logger.ILogging {
	id := fmt.Sprint(goid())
	if l, ok := r.loggers.Get(id); ok {
		return l.(logger.ILogging)
	}
	return r.logger
}

func appendEngines(engines []string, ext ...string) []string {
	addEngine := make([]string, 0, len(ext))
	for _, n := range ext {
		var b bool
		for _, en := range engines {
			if en == n {
				b = true
				continue
			}
		}
		if !b {
			addEngine = append(addEngine, n)
		}
	}
	return append(engines, addEngine...)
}
func goid() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		return 0
	}
	return id
}
