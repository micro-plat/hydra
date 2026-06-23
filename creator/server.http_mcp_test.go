package creator

import (
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

// Test_httpBuilder_MCP 验证 httpBuilder.MCP() 的启用/禁用语义与 provider 注入的协作：
//   - provider 未注入（未 import mcp）时 .MCP() 为空操作，不 panic；
//   - provider 注入后 .MCP()（默认启用）调用 provider 并返回 *httpBuilder 保持链式；
//   - .MCP(禁用选项) 显式禁用，不调用 provider。
//
// 说明：creator 不可反向 import mcp（开闭原则），故禁用选项用内联 services.MCPOption 表达，
// 等价于业务侧的 mcp.WithDisable()。
func Test_httpBuilder_MCP(t *testing.T) {
	// 1. provider 未注入 → 空操作不 panic
	services.RegisterMCPHandlerProvider(nil)
	b0 := newHTTP(global.Web, ":0")
	b0.MCP()
	b0.MCP(func(c *services.MCPConf) { c.Enable = false })

	// 2. 注入 provider（返回未命名 func(context.IContext) interface{}，
	//    与真实 NewJSONRPCHandler 同签名，确保 Custom→swapFunc 类型断言命中而非误入 createObject）
	invoked := 0
	services.RegisterMCPHandlerProvider(func() func(context.IContext) interface{} {
		invoked++
		return func(context.IContext) interface{} { return nil }
	})
	defer services.RegisterMCPHandlerProvider(nil)

	b := newHTTP(global.Web, ":0")
	ret := b.MCP()
	if invoked != 1 {
		t.Fatalf("MCP() 默认启用应调用 provider 1 次, 实际 %d", invoked)
	}
	if ret != b {
		t.Fatal("MCP() 应返回 *httpBuilder 以保持链式")
	}

	// 3. 显式禁用
	b.MCP(func(c *services.MCPConf) { c.Enable = false })
	if invoked != 1 {
		t.Fatalf("MCP(禁用) 不应调用 provider, 实际 %d", invoked)
	}
}
