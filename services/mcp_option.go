package services

// MCPConf MCP 服务端点配置，由 creator.httpBuilder.MCP 消费。
// MCP 复用 web/api 服务器的 /mcp 路由，复用完整中间件链（含 JwtAuth）——
// 即 MCP 客户端须与 HTTP 客户端一样携带有效登录态，不支持匿名访问。
type MCPConf struct {
	Enable bool   // 是否启用 MCP（默认 true）
	Path   string // MCP 端点路径（默认 "/mcp"）
}

// MCPOption MCP 配置选项，由 mcp 子包的 WithXxx 构造（如 mcp.WithDisable、mcp.WithPath）。
// 定义在 services 共享层，使 creator（配置构建）与 mcp（选项构造）均可引用而互不反向依赖。
type MCPOption func(*MCPConf)

// DefaultMCPConf 默认 MCP 配置：启用、路径 /mcp。
func DefaultMCPConf() *MCPConf {
	return &MCPConf{Enable: true, Path: "/mcp"}
}
