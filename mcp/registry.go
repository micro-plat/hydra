package mcp

import (
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/services"
)

// Tool 一个 MCP 工具的注册信息。
type Tool struct {
	Path     string          // 源 HTTP 路径（构造 tool 执行上下文时回填 path）
	Name     string          // tool 名
	Handler  context.Handler // 包装后的业务处理函数（tools/call 调用）
	ReqType  reflect.Type    // 输入类型（保留，供校验/调试）
	RespType reflect.Type    // 输出类型
	Config   *Config         // 最终配置（含 inputSchema/outputSchema/desc）
}

var (
	toolsMu sync.RWMutex
	tools   = map[string]*Tool{} // toolName → Tool
)

// registerTool services.MCPRegisterHook 的注入目标：把 path+h+opts 登记为 MCP tool(s)。
// h 必须是 typed 函数 func(ctx context.IContext, req)(resp, error)，或含 typed *Handle 方法的对象。
// 非 typed（如旧签名 func(ctx) interface{}）无法作为 MCP tool：显式 .WithMCP()/.MCP() 登记它属编程错误，
// 启动期 panic 暴露（与 buildConfig 配置错误同处理，fail-fast，避免静默漏登记）。
// opts 见 buildConfig（string/Option/Config）。
func registerTool(path string, h interface{}, opts ...interface{}) {
	items, ok := resolve(h, path)
	if !ok {
		panic(fmt.Sprintf("mcp: 无法将 %q 登记为 MCP 工具：handler 非 typed 签名。要求 func(ctx context.IContext, req)(resp, error) 或含 typed *Handle 方法的对象；旧签名（如 func(ctx) interface{}）不支持，请改 typed 签名或去掉 .WithMCP()", path))
	}
	cfg, err := buildConfig(opts...)
	if err != nil {
		panic(err)
	}
	objConfigs := objectConfigs(h)

	for _, it := range items {
		registerOne(buildTool(it, cfg, objConfigs))
	}
}

// buildTool 合并优先级：opts（显式）> 对象服务声明 > 默认。默认名取自源路径，默认描述为 "调用 {name}"。
func buildTool(it resolved, opts *Config, objConfigs map[string]*Config) *Tool {
	cfg := &Config{
		Name:         toolName(it.sourcePath),
		Desc:         "调用 " + toolName(it.sourcePath),
		InputSchema:  buildInputSchema(it.reqType),
		OutputSchema: buildOutputSchema(it.respType),
	}
	// 对象服务声明（按方法名）覆盖默认
	if obj, ok := objConfigs[it.methodName]; ok {
		mergeConfig(cfg, obj)
	}
	// 显式 opts 优先级最高
	mergeConfig(cfg, opts)

	return &Tool{
		Path:     it.sourcePath,
		Name:     cfg.Name,
		Handler:  it.handler,
		ReqType:  it.reqType,
		RespType: it.respType,
		Config:   cfg,
	}
}

// objectConfigs 若 h 是嵌入 mcp.Service 的对象服务，返回其方法级配置；否则 nil。
func objectConfigs(h interface{}) map[string]*Config {
	if s, ok := h.(IMCPService); ok {
		return s.MCPConfigs()
	}
	return nil
}

func registerOne(t *Tool) {
	toolsMu.Lock()
	tools[t.Name] = t
	toolsMu.Unlock()
}

// List 返回所有已注册 tool（tools/list 用），按 name 排序保证稳定输出。
func List() []*Tool {
	toolsMu.RLock()
	defer toolsMu.RUnlock()
	list := make([]*Tool, 0, len(tools))
	for _, t := range tools {
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
	return list
}

// Get 按 tool 名查询单个 tool（tools/call 用）。
func Get(name string) (*Tool, bool) {
	toolsMu.RLock()
	defer toolsMu.RUnlock()
	t, ok := tools[name]
	return t, ok
}

func init() {
	services.RegisterMCPHook(registerTool)
	// 注入 /mcp 路由处理器提供者：creator.httpBuilder.MCP() 启用时据此自动注册 /mcp。
	// NewJSONRPCHandler 返回未命名 func(context.IContext) interface{}，经 Custom 装箱后
	// swapFunc 的类型断言方可命中（命名类型 context.Handler 会断言失败，见 services.MCPHandlerProvider 注释）。
	services.RegisterMCPHandlerProvider(NewJSONRPCHandler)
}
