package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/types"
)

const defHandling = "Handling"
const defHandler = "Handle"
const defHandled = "Handled"
const defFallback = "Fallback"
const defClose = "Close"

//IService 服务注册接口

// IMicroRegistry 微服务注册接口
type IMicroRegistry func(name string, h interface{}, r ...router.Option) IService
type IService interface {
	Group(name ...string) IService

	//GetGroup 获取服务的分组信息
	GetGroup(serverType string, service string, method ...string) string

	//Micro 注册为微服务，包括api,web,rpc
	Micro(name string, h interface{}, r ...router.Option) IService

	//Flow 注册为流程,包括mqc,cron
	Flow(name string, h interface{}) IService

	//API 注册为http api服务
	API(name string, h interface{}, r ...router.Option) IService

	//Web 注册为web服务
	Web(name string, h interface{}, r ...router.Option) IService

	//RPC 注册为RPC服务
	RPC(name string, h interface{}, r ...router.Option) IService

	//WSS 注册为websocket服务
	WSS(name string, h interface{}, r ...router.Option) IService

	//AIGW 注册为AI网关服务
	AIGW(name string, h interface{}, r ...router.Option) IService

	//MQC 注册为消息消费服务
	MQC(name string, h interface{}, queues ...string) IService

	//CRON 注册为定时任务服务
	CRON(name string, h interface{}, crons ...string) IService

	//custome 注册为自定义服务器的服务
	Custom(tp string, name string, h interface{}, ext ...interface{}) IService

	//MCP 注册为 MCP 工具（仅 MCP，不创建 HTTP 路由），opts 透传给 mcp 子包
	MCP(name string, h interface{}, opts ...interface{}) IService

	//WithMCP 将最近一次注册的服务同时登记为 MCP 工具，opts 透传给 mcp 子包
	WithMCP(opts ...interface{}) IService

	//Remove 移除已注册服务（重启服务器生效）
	Remove(path string, tp ...string)

	//Clear 清除所有注册服务
	Clear(tp ...string)

	//HookRemove 移除钩子函数
	HookRemove(hs ...interface{})

	//RegisterServer 注册新的服务器类型
	RegisterServer(tp string, f func(g *Unit, ext ...interface{}) error, remove func(path string))

	//OnSetup 服务器初始化勾子，服务器配置初始化前执行(首次启动前 或 收到配置变更后需要重启服务器前)
	OnSetup(h func(app.IAPPConf) error, tps ...string)

	//OnStarting 服务器启动前勾子，服务器启动前执行
	OnStarting(h func(app.IAPPConf) error, tps ...string)

	//OnStarted 服务器启动完成勾子，服务器启动后执行
	OnStarted(h func(app.IAPPConf) error, tps ...string)

	//OnClosing 服务器关闭勾子，服务器关闭后执行
	OnClosing(h func(app.IAPPConf) error, tps ...string)

	//OnHandleExecuting Handle勾子，Handle执行前执行
	OnHandleExecuting(h context.Handler, tps ...string)

	//OnHandleExecuted Handle勾子，Handle执行后执行
	OnHandleExecuted(h context.Handler, tps ...string)
}

// Def 服务注册管理
var Def = New()

// New 构建服务组件
func New() *regist {
	return &regist{
		servers: make(map[string]*serverServices),
		caches:  make(map[string]map[string]interface{}),
	}
}

// regist  本地服务
type regist struct {
	group          string
	servers        map[string]*serverServices
	caches         map[string]map[string]interface{}
	lastPath       string
	lastHandler    interface{}
	lastServerType string // 最近一次注册的服务器类型（错误提示用）
	lastHTTPCapable bool  // 最近一次注册是否为 HTTP 类型（Micro/API/Web），决定 .WithMCP() 是否有效
}

// Micro 注册为微服务包括api,web,rpc
func (s *regist) Micro(name string, h interface{}, r ...router.Option) IService {
	s.API(name, h, r...)
	s.Web(name, h, r...)
	s.RPC(name, h, r...)
	s.lastHTTPCapable = true // Micro 含 API/Web（HTTP），允许后续 .WithMCP()（覆盖末尾 RPC 的 false）
	return s
}
func (s *regist) Group(name ...string) IService {
	s.group = strings.Trim(types.GetStringByIndex(name, 0), "/")
	return s
}

// Flow 注册为流程服务，包括mqc,cron
func (s *regist) Flow(name string, h interface{}) IService {
	s.Custom(global.MQC, name, h)
	s.Custom(global.CRON, name, h)
	return s
}

// API 注册为API服务
func (s *regist) API(name string, h interface{}, ext ...router.Option) IService {
	v := make([]interface{}, 0, len(ext))
	for _, e := range ext {
		v = append(v, e)
	}
	return s.Custom(global.API, name, h, v...)
}

// Web 注册为web服务
func (s *regist) Web(name string, h interface{}, ext ...router.Option) IService {
	v := make([]interface{}, 0, len(ext))
	for _, e := range ext {
		v = append(v, e)
	}
	return s.Custom(global.Web, name, h, v...)
}

// RPC 注册为rpc服务
func (s *regist) RPC(name string, h interface{}, ext ...router.Option) IService {
	v := make([]interface{}, 0, len(ext))
	for _, e := range ext {
		v = append(v, e)
	}
	return s.Custom(global.RPC, name, h, v...)
}

// WSS 注册为websocket服务
func (s *regist) WSS(name string, h interface{}, ext ...router.Option) IService {
	v := make([]interface{}, 0, len(ext))
	for _, e := range ext {
		v = append(v, e)
	}
	return s.Custom(global.WSSServer, name, h, v...)
}

// AIGW 注册为AI网关服务
func (s *regist) AIGW(name string, h interface{}, ext ...router.Option) IService {
	v := make([]interface{}, 0, len(ext))
	for _, e := range ext {
		v = append(v, e)
	}
	return s.Custom(global.AIGW, name, h, v...)
}

// MQC 注册为消息队列服务
func (s *regist) MQC(name string, h interface{}, queues ...string) IService {
	v := make([]interface{}, 0, len(queues))
	for _, e := range queues {
		v = append(v, e)
	}
	return s.Custom(global.MQC, name, h, v...)
}

// CRON 注册为定时任务服务
func (s *regist) CRON(name string, h interface{}, crons ...string) IService {
	v := make([]interface{}, 0, len(crons))
	for _, e := range crons {
		v = append(v, e)
	}
	return s.Custom(global.CRON, name, h, v...)
}

// Remove 移除已注册服务（重启服务生效）
func (s *regist) Remove(path string, tp ...string) {
	if len(tp) == 0 {
		for _, server := range s.servers {
			server.Remove(path)
		}
		return
	}
	for _, t := range tp {
		if server, ok := s.servers[t]; ok {
			server.Remove(path)
		}
	}

}

// Clear 清除所有注册服务
func (s *regist) Clear(tp ...string) {
	if len(tp) == 0 {
		for _, server := range s.servers {
			server.Clear()
		}
		return
	}
	for _, t := range tp {
		if server, ok := s.servers[t]; ok {
			server.Clear()
		}
	}

}

// HookRemove 移除钩子函数
func (s *regist) HookRemove(hs ...interface{}) {
	for _, server := range s.servers {
		server.HookRemove(hs...)
	}
}

// Custom 自定义服务注册
func (s *regist) Custom(tp string, name string, h interface{}, ext ...interface{}) IService {
	name = s.normalize(name)
	s.lastPath, s.lastHandler = name, h
	s.lastServerType = tp
	s.lastHTTPCapable = tp == global.Web || tp == global.API
	s.get(tp).Register(s.group, name, h, ext...)
	return s
}

// normalize 规范化服务路径：去两头'/'，按 group 补前缀
func (s *regist) normalize(name string) string {
	name = fmt.Sprintf("/%s", strings.Trim(name, "/"))
	if s.group != "" {
		name = fmt.Sprintf("/%s%s", s.group, name)
	}
	return name
}

// MCP 注册为 MCP 工具（仅 MCP，不创建 HTTP 路由）。opts 透传给 mcp 子包的注册钩子。
// 未 import mcp 子包（钩子为 nil）时为空操作，仅返回 s 以保持链式。
func (s *regist) MCP(name string, h interface{}, opts ...interface{}) IService {
	if MCPRegisterHook != nil {
		MCPRegisterHook(s.normalize(name), h, opts...)
	}
	return s
}

// WithMCP 将最近一次注册的 HTTP 服务（Micro/API/Web）同时登记为 MCP 工具。
// MCP 工具经 web/api 服务器的 /mcp 端点对外提供，仅 HTTP 类型注册有意义；
// 非 HTTP 类型（RPC/CRON/MQC/WSS/AIGW）的 handler 非 HTTP 请求处理器，登记无意义。
// 依赖 Custom 写入的 lastPath/lastHandler/lastHTTPCapable；opts 透传给 mcp 子包注册钩子。
// 未 import mcp 子包（钩子为 nil）或无最近注册时为空操作；最近注册非 HTTP 类型则启动期 panic（fail-fast）。
func (s *regist) WithMCP(opts ...interface{}) IService {
	if MCPRegisterHook == nil || s.lastHandler == nil {
		return s
	}
	if !s.lastHTTPCapable {
		panic(fmt.Sprintf("mcp: .WithMCP() 仅支持 Micro/API/Web 注册的服务，最近注册类型为 %q（非 HTTP），无法登记为 MCP 工具；请改用 Micro/API/Web 注册或去掉 .WithMCP()", s.lastServerType))
	}
	MCPRegisterHook(s.lastPath, s.lastHandler, opts...)
	return s
}

// MCPRegisterHook MCP 工具注册钩子，由 mcp 子包在 init 时通过 RegisterMCPHook 注入。
// 遵循开闭原则：services 不反向依赖 mcp，仅持有函数变量。未注入时为 nil，MCP/WithMCP 为空操作。
var MCPRegisterHook func(path string, h interface{}, opts ...interface{})

// RegisterMCPHook 注册 MCP 工具注册钩子，供 mcp 子包注入。
func RegisterMCPHook(fn func(path string, h interface{}, opts ...interface{})) {
	MCPRegisterHook = fn
}

// MCPHandlerProvider MCP /mcp 路由处理器提供者，由 mcp 子包注入。
// creator.httpBuilder.MCP() 启用时通过它获取 /mcp 的 JSON-RPC 处理器并注册到当前 http 服务器。
// 遵循开闭原则：services/creator 不反向依赖 mcp，仅持有函数变量。未注入时为 nil，.MCP() 为空操作。
//
// 注意：返回类型必须是未命名 func(context.IContext) interface{}，而非命名类型 context.Handler。
// 因注册经 Custom(interface{}) 装箱后，reflectHandle.swapFunc 用 .(func(context.IContext) interface{}) 做类型断言，
// 命名类型与未命名类型即使底层相同也互不 identical，断言会失败并误入 createObject 分支导致 panic。
var MCPHandlerProvider func() func(context.IContext) interface{}

// RegisterMCPHandlerProvider 注册 MCP /mcp 处理器提供者，供 mcp 子包注入。
func RegisterMCPHandlerProvider(fn func() func(context.IContext) interface{}) {
	MCPHandlerProvider = fn
}

// RegisterServer 注册服务器
func (s *regist) RegisterServer(tp string, f func(g *Unit, ext ...interface{}) error, v func(path string)) {
	if _, ok := s.servers[tp]; ok {
		panic(fmt.Errorf("服务%s已存在，不能重复注册", tp))
	}
	if f != nil {
		s.servers[tp] = newServerServices(f, v)
		return
	}
	s.servers[tp] = newServerServices(nil, nil)
}

// OnSetup 服务器初始化勾子，服务器配置初始化前执行(首次启动前或收到配置变更后需要重启服务器前)
func (s *regist) OnSetup(h func(app.IAPPConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddSetup(h); err != nil {
			panic(fmt.Errorf("%s OnSetup %v", typ, err))
		}
	}
}

// OnStarting 处理服务器启动前
func (s *regist) OnStarting(h func(app.IAPPConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddStarting(h); err != nil {
			panic(fmt.Errorf("%s OnStarting %v", typ, err))
		}
	}
}

// OnStarted 处理服务器启动后
func (s *regist) OnStarted(h func(app.IAPPConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddStarted(h); err != nil {
			panic(fmt.Errorf("%s OnStarted %v", typ, err))
		}
	}
}

// OnClosing 处理服务器关闭
func (s *regist) OnClosing(h func(app.IAPPConf) error, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddClosing(h); err != nil {
			panic(fmt.Errorf("%s OnClosing %v", typ, err))
		}
	}
}

// OnHandleExecuting 处理handling业务
func (s *regist) OnHandleExecuting(h context.Handler, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddHandleExecuting(h); err != nil {
			panic(fmt.Errorf("%s OnHandleExecuting %v", typ, err))
		}
	}
}

// Handled 处理Handled业务
func (s *regist) OnHandleExecuted(h context.Handler, tps ...string) {
	if len(tps) == 0 {
		tps = global.Def.ServerTypes
	}
	for _, typ := range tps {
		if err := s.get(typ).AddHandleExecuted(h); err != nil {
			panic(fmt.Errorf("%s OnHandleExecuted %v", typ, err))
		}
	}
}

// Has 服务器是否注册了某个服务
func (s *regist) Has(serverType string, service, method string) (ok bool) {
	if s.get(serverType).Has(service) {
		return true
	}
	return s.get(serverType).Has(fmt.Sprintf("%s$%s", service, method))
}

// GetHandleExecutings 获取handle预处理勾子
func (s *regist) GetHandleExecutings(serverType string) []context.IHandler {
	return s.get(serverType).GetHandleExecutings()
}

// GetHandleExecuted 获取handle后处理勾子
func (s *regist) GetHandleExecuted(serverType string) []context.IHandler {
	return s.get(serverType).GetHandleExecuteds()
}

// GetHandler 获取服务对应的处理函数
func (s *regist) GetHandler(serverType string, service string) (context.IHandler, bool) {
	return s.get(serverType).GetHandlers(service)
}

// GetGroup 获取服务的分组信息
func (s *regist) GetGroup(serverType string, service string, method ...string) string {
	if len(method) == 0 {
		return s.get(serverType).GetGroup(service)
	}
	return s.get(serverType).GetGroup(registry.Join(service, "$"+strings.ToLower(types.GetStringByIndex(method, 0))))
}

// GetRawPathAndTag 获取服务原始注册路径与方法名(restful服务的tag值为空)
func (s *regist) GetRawPathAndTag(serverType string, service string) (path string, tagName string, ok bool) {
	return s.get(serverType).GetRawPathAndTag(service)
}

// GetHandling 获取预处理函数
func (s *regist) GetHandlings(serverType string, service string) []context.IHandler {
	return s.get(serverType).GetHandlings(service)
}

// GetHandling 获取后处理函数
func (s *regist) GetHandleds(serverType string, service string) []context.IHandler {
	return s.get(serverType).GetHandleds(service)
}

// GetFallback 获取服务对应的降级函数
func (s *regist) GetFallback(serverType string, service string) (context.IHandler, bool) {
	return s.get(serverType).GetFallback(service)
}

func (s *regist) get(tp string) *serverServices {
	if v, ok := s.servers[tp]; ok {
		return v
	}
	panic(fmt.Sprintf("不支持的服务器类型:%s", tp))
}

// DoStarting 执行服务启动函数
func (s *regist) DoStarting(c app.IAPPConf) error {
	return s.get(c.GetServerConf().GetServerType()).DoStarting(c)

}

// DoStarted 执行服务启动函数
func (s *regist) DoStarted(c app.IAPPConf) error {
	return s.get(c.GetServerConf().GetServerType()).DoStarted(c)

}

// DoSetup 执行服务启动函数
func (s *regist) DoSetup(c app.IAPPConf) error {
	return s.get(c.GetServerConf().GetServerType()).DoSetup(c)

}

// DoClosing 执行服务关闭函数
func (s *regist) DoClosing(c app.IAPPConf) error {
	return s.get(c.GetServerConf().GetServerType()).DoClosing(c)
}

// GetClosers 获取资源释放函数
func (s *regist) Close() error {
	var sb strings.Builder
	for _, r := range s.servers {
		for _, cs := range r.GetClosingHandlers() {
			if err := cs(); err != nil {
				sb.WriteString(err.Error())
				sb.WriteString("\n")
			}
		}
	}
	if sb.Len() == 0 {
		return nil
	}
	return errors.New(strings.Trim(sb.String(), "\n"))
}

//-----------------------注册缓存-------------------------------------------

// init 处理服务初始化及特殊注册函数
func init() {
	Def.servers[global.API] = newServerServices(func(g *Unit, ext ...interface{}) error {
		return API.Add(g.Path, g.Service, g.Actions, ext...)
	}, API.Remove)
	Def.servers[global.Web] = newServerServices(func(g *Unit, ext ...interface{}) error {
		return WEB.Add(g.Path, g.Service, g.Actions, ext...)
	}, WEB.Remove)
	Def.servers[global.RPC] = newServerServices(func(g *Unit, ext ...interface{}) error {
		return RPC.Add(g.Path, g.Service, g.Actions, ext...)
	}, RPC.Remove)

	wssServices := newServerServices(func(g *Unit, ext ...interface{}) error {
		return WSS.Add(g.Path, g.Service, g.Actions, ext...)
	}, WSS.Remove)
	Def.servers[global.WSSServer] = wssServices
	Def.servers[global.WSSClient] = wssServices
	Def.servers[global.AIGW] = newServerServices(func(g *Unit, ext ...interface{}) error {
		return AIGW.Add(g.Path, g.Service, g.Actions, ext...)
	}, AIGW.Remove)
	Def.servers[global.CRON] = newServerServices(func(g *Unit, ext ...interface{}) error {
		for _, t := range ext {
			CRON.Add(t.(string), g.Service)
		}

		routerCRON.Add(g.Path, g.Service, g.Actions)
		return nil
	}, routerCRON.Remove)
	Def.servers[global.MQC] = newServerServices(func(g *Unit, ext ...interface{}) error {
		for _, t := range ext {
			qname, concurrency := extractParts(t.(string), 10)
			MQC.Add(qname, g.Service, concurrency)
		}
		routerMQC.Add(g.Path, g.Service, g.Actions)
		return nil
	}, routerMQC.Remove)
}
func extractParts(input string, defNumber int) (prefix string, number int) {
	// 定义正则表达式，匹配 { 前的部分和 {} 中的数字
	re := regexp.MustCompile(`^(.*)\{(\d+)\}$`)

	// 查找匹配项
	matches := re.FindStringSubmatch(input)
	if len(matches) < 3 {
		return input, defNumber
	}

	// 返回 { 前的部分和 {} 中的数字部分
	return matches[1], types.GetInt(matches[2], defNumber)
}
