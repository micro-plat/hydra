package mcp

import "github.com/micro-plat/hydra/services"

// 本文件提供 MCP **服务器端点**配置选项（用于 hydra.Conf.Web/API(...).MCP(opts...)），
// 与 option.go 中的 **工具登记**选项（Name/Desc/...，用于 .WithMCP）相区分：
//   - 工具登记选项：mcp.Option = func(*Config)，由 services.MCPRegisterHook 消费；
//   - 服务器端点选项：services.MCPOption = func(*MCPConf)，由 creator.httpBuilder.MCP 消费。

// WithDisable 禁用 MCP 服务。用于 hydra.Conf.Web/API(...).MCP(mcp.WithDisable())。
func WithDisable() services.MCPOption {
	return func(c *services.MCPConf) { c.Enable = false }
}

// WithPath 自定义 MCP 端点路径（默认 "/mcp"）。用于 hydra.Conf.Web/API(...).MCP(mcp.WithPath("/v1/mcp"))。
func WithPath(path string) services.MCPOption {
	return func(c *services.MCPConf) { c.Path = path }
}
