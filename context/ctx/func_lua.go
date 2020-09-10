package ctx

import (
	"github.com/micro-plat/hydra/context"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func (r *funcs) LuaFuncs() context.LuaModules {
	r.lonce.Do(func() {
		r.luaFuncs["request"] = map[string]lua.LGFunction{
			"get_string":     r.getLuaRquestString,
			"get_int":        r.getLuaRquestInt,
			"get_param":      r.getLuaRequestParam,
			"get_path":       r.getLuaRequestPath,
			"get_router":     r.getLuaRouter,
			"get_header":     r.getLuaHeader,
			"get_cookie":     r.getLuaCookie,
			"get_client_ip":  r.getLuaClientIP,
			"get_request_id": r.getLuaRquestID,
		}
		r.luaFuncs["response"] = map[string]lua.LGFunction{
			"get_status":  r.getLuaResponseStatus,
			"get_content": r.getLuaResponseContent,
		}
	})
	return r.luaFuncs
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
