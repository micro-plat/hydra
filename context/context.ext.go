package context

import "strings"

//IsMicroServer 是否是远程服务api,web,rpc
func (c *Context) IsMicroServer() bool {
	return c.IsAPIServer() || c.IsWebServer() || c.IsRPCServer() || c.IsWSServer()
}

//IsFlowServer 是否是流程服务mqc,cron
func (c *Context) IsFlowServer() bool {
	return c.IsMQCServer() || c.IsCRONServer()
}

//IsAPIServer 是否是http api服务
func (c *Context) IsAPIServer() bool {
	return strings.ToLower(c.GetContainer().GetServerType()) == "api"
}

//IsWSServer 是否是web socket服务
func (c *Context) IsWSServer() bool {
	return strings.ToLower(c.GetContainer().GetServerType()) == "ws"
}

//IsWebServer 是否是web服务
func (c *Context) IsWebServer() bool {
	return strings.ToLower(c.GetContainer().GetServerType()) == "web"
}

//IsRPCServer 是否是rpc服务
func (c *Context) IsRPCServer() bool {
	return strings.ToLower(c.GetContainer().GetServerType()) == "rpc"
}

//IsMQCServer 是否是mqc服务
func (c *Context) IsMQCServer() bool {
	return strings.ToLower(c.GetContainer().GetServerType()) == "mqc"
}

//IsCRONServer 是否是cron服务
func (c *Context) IsCRONServer() bool {
	return strings.ToLower(c.GetContainer().GetServerType()) == "cron"
}
