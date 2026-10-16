package mcp

import "github.com/micro-plat/hydra/services"

func init() {
	services.RegisterMCPAuthorizeHandlerProvider(NewOAuthAuthorizeHandler)
}
