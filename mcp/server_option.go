package mcp

import "github.com/micro-plat/hydra/services"

// OAuthConfig configures MCP OAuth authorization.
type OAuthConfig = services.MCPOAuthConf

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

// WithOAuth enables MCP OAuth authorization. Empty config uses defaults.
func WithOAuth(config ...OAuthConfig) services.MCPOption {
	return func(c *services.MCPConf) {
		oauth := services.DefaultMCPOAuthConf()
		if len(config) > 0 {
			cfg := config[0]
			if cfg.AuthorizePath != "" {
				oauth.AuthorizePath = cfg.AuthorizePath
			}
			if cfg.LoginURL != "" {
				oauth.LoginURL = cfg.LoginURL
			}
			if len(cfg.Scopes) > 0 {
				oauth.Scopes = cfg.Scopes
			}
		}
		c.OAuth = oauth
	}
}
