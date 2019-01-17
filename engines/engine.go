package engines

import (
	"strings"

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
func NewServiceEngine(conf conf.IServerConf, registryAddr string, logger logger.ILogging) (e *ServiceEngine, err error) {
	e = &ServiceEngine{IServerConf: conf, registryAddr: registryAddr, logger: logger}
	e.StandardComponent = component.NewStandardComponent("sys.engine", e)
	e.Invoker = rpc.NewInvoker(conf.GetPlatName(), conf.GetSysName(), registryAddr)
	e.IComponentCache = component.NewStandardCache(e, "cache")
	e.IComponentDB = component.NewStandardDB(e, "db")
	e.IComponentInfluxDB = component.NewStandardInfluxDB(e, "influx")
	e.IComponentQueue = component.NewStandardQueue(e, "queue")
	e.IComponentGlobalVarObject = component.NewGlobalVarObjectCache(e)
	if e.registry, err = registry.NewRegistryWithAddress(registryAddr, logger); err != nil {
		return
	}

	if err = e.loadEngineServices(); err != nil {
		return nil, err
	}
	// if err = e.LoadComponents(fmt.Sprintf("./%s.so", conf.GetPlatName()),
	// 	fmt.Sprintf("./%s.so", conf.GetSysName()),
	// 	fmt.Sprintf("./%s_%s.so", conf.GetPlatName(), conf.GetSysName())); err != nil {
	// 	return
	// }
	e.StandardComponent.AddRPCProxy(e.RPCProxy())
	err = e.StandardComponent.LoadServices()
	return
}

//SetHandler 设置handler
func (r *ServiceEngine) SetHandler(h component.IComponentHandler) error {
	if h == nil {
		return nil
	}
	funcs := h.GetInitializings()
	for _, f := range funcs {
		if err := f(r); err != nil {
			return err
		}
	}
	r.cHandler = h
	svs := h.GetServices()
	for group, handlers := range svs {
		for name, handler := range handlers {
			r.StandardComponent.AddCustomerService(name, handler, group, h.GetTags(name)...)
		}
	}
	return r.StandardComponent.LoadServices()
}

//UpdateVarConf 更新var配置参数
func (r *ServiceEngine) UpdateVarConf(conf conf.IServerConf) {
	r.SetVarConf(conf.GetVarConfClone())
	r.SetSubConf(conf.GetSubConfClone())
}

//GetServices 获取组件提供的所有服务
func (r *ServiceEngine) GetServices() []string {
	return r.GetGroupServices(component.GetGroupName(r.GetServerType())...)
}

//Execute 执行外部请求
func (r *ServiceEngine) Execute(ctx *context.Context) (rs interface{}) {
	if ctx.Request.CircuitBreaker.IsOpen() { //熔断开关打开，则自动降级
		rf := r.StandardComponent.Fallback(ctx)
		if r, ok := rf.(error); ok && r == component.ErrNotFoundService {
			ctx.Response.MustContent(ctx.Request.CircuitBreaker.GetDefStatus(), ctx.Request.CircuitBreaker.GetDefContent())
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
		if rs = r.Handle(ctx); ctx.Response.HasError(rs) {
			return rs
		}
	}

	//当前服务器后处理
	if r.cHandler != nil && r.cHandler.GetHandleds() != nil {
		hdd := r.cHandler.GetHandleds()
		for _, h := range hdd {
			if rh := h(ctx); ctx.Response.HasError(rh) {
				return rh
			}
		}
	}

	//当前引擎后处理
	if rd := r.Handled(ctx); ctx.Response.HasError(rd) {
		return rd
	}
	return rs
}

//Handling 每次handle执行前执行
func (r *ServiceEngine) Handling(c *context.Context) (rs interface{}) {
	c.SetRPC(r.Invoker)
	return nil
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
	r.Invoker.Close()
	r.IComponentGlobalVarObject.Close()
	r.IComponentCache.Close()
	r.IComponentInfluxDB.Close()
	r.IComponentQueue.Close()
	r.IComponentDB.Close()
	return nil
}
func (r *ServiceEngine) GetLooger() logger.ILogging {
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
