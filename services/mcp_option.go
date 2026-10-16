package services

// MCPConf is the server endpoint configuration consumed by creator.httpBuilder.MCP.
// MCP reuses the web/api server middleware chain, including JwtAuth.
type MCPConf struct {
	Enable bool
	Path   string
	OAuth  *MCPOAuthConf
}

// MCPOAuthConf configures MCP OAuth authorization.
type MCPOAuthConf struct {
	Enable        bool
	AuthorizePath string
	LoginURL      string
	Scopes        []string
}

// MCPOption is built by the mcp package and consumed by creator.
type MCPOption func(*MCPConf)

// DefaultMCPConf returns the default MCP endpoint config.
func DefaultMCPConf() *MCPConf {
	return &MCPConf{Enable: true, Path: "/mcp"}
}

// DefaultMCPOAuthConf returns the default MCP OAuth config.
func DefaultMCPOAuthConf() *MCPOAuthConf {
	return &MCPOAuthConf{
		Enable:        true,
		AuthorizePath: "/oauth/authorize",
		Scopes:        []string{"mcp:tools:list", "mcp:tools:call"},
	}
}
