package ctx

import (
	"fmt"

	xlua "github.com/micro-plat/hydra/pkgs/lua"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func (r *funcs) LuaFuncs() xlua.Modules {
	modules := make(xlua.Modules)
	modules["request"] = map[string]lua.LGFunction{
		"getString":    r.getLuaRquestString,
		"getInt":       r.getLuaRquestInt,
		"getParam":     r.getLuaRequestParam,
		"getPath":      r.getLuaRequestPath,
		"getRouter":    r.getLuaRouter,
		"getHeader":    r.getLuaHeader,
		"getCookie":    r.getLuaCookie,
		"getClientIP":  r.getLuaClientIP,
		"getRequestID": r.getLuaRquestID,
	}
	modules["response"] = map[string]lua.LGFunction{
		"getStatus":  r.getLuaResponseStatus,
		"getContent": r.getLuaResponseContent,
	}

	return modules
}

func (c *funcs) getLuaRquestString(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	def := lua.LVAsString(L.Get(2))
	ret := c.ctx.request.GetString(s, def)
	L.Push(lua.LString(ret))
	return 1
}
func (c *funcs) getLuaRquestInt(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	def := lua.LVAsNumber(L.Get(2))
	ret := c.ctx.request.GetInt(s, int(float32(def)))
	L.Push(lua.LNumber(ret))
	return 1
}
func (c *funcs) getLuaRequestParam(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	ret := c.ctx.request.Param(s)
	L.Push(lua.LString(ret))
	return 1
}
func (c *funcs) getLuaRequestPath(L *lua.LState) int {
	ret := c.ctx.request.path.GetRequestPath()
	L.Push(lua.LString(ret))
	return 1
}
func (c *funcs) getLuaRouter(L *lua.LState) int {
	ret := c.ctx.request.path.GetRouter()
	L.Push(luar.New(L, ret))
	return 1
}
func (c *funcs) getLuaHeader(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	ret := c.ctx.request.path.GetHeader(s)
	L.Push(lua.LString(ret))
	return 1
}
func (c *funcs) getLuaCookie(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	ret, ok := c.ctx.request.path.GetCookie(s)
	L.Push(lua.LString(ret))
	L.Push(lua.LBool(ok))
	return 2
}
func (c *funcs) getLuaClientIP(L *lua.LState) int {
	fmt.Println("abcef")
	ret := c.ctx.user.GetClientIP()
	L.Push(lua.LString(ret))
	return 1
}
func (c *funcs) getLuaRquestID(L *lua.LState) int {
	ret := c.ctx.user.GetRequestID()
	L.Push(lua.LString(ret))
	return 1
}
func (c *funcs) getLuaResponseStatus(L *lua.LState) int {
	ret := c.ctx.response.getStatus()
	L.Push(lua.LNumber(ret))
	return 1
}
func (c *funcs) getLuaResponseContent(L *lua.LState) int {
	ret := c.ctx.response.getContent()
	L.Push(lua.LString(ret))
	return 1
}
