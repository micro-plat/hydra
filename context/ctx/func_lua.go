package ctx

import (
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func (c *Ctx) getLuaRquestString(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	def := lua.LVAsString(L.Get(2))
	ret := c.request.GetString(s, def)
	L.Push(lua.LString(ret))
	return 1
}
func (c *Ctx) getLuaRquestInt(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	def := lua.LVAsNumber(L.Get(2))
	ret := c.request.GetInt(s, int(float32(def)))
	L.Push(lua.LNumber(ret))
	return 1
}
func (c *Ctx) getLuaRequestParam(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	ret := c.request.Param(s)
	L.Push(lua.LString(ret))
	return 1
}
func (c *Ctx) getLuaRequestPath(L *lua.LState) int {
	ret := c.request.path.GetRequestPath()
	L.Push(lua.LString(ret))
	return 1
}
func (c *Ctx) getLuaRouter(L *lua.LState) int {
	ret := c.request.path.GetRouter()
	L.Push(luar.New(L, ret))
	return 1
}
func (c *Ctx) getLuaHeader(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	ret := c.request.path.GetHeader(s)
	L.Push(lua.LString(ret))
	return 1
}
func (c *Ctx) getLuaCookie(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	ret, ok := c.request.path.GetCookie(s)
	L.Push(lua.LString(ret))
	L.Push(lua.LBool(ok))
	return 2
}
func (c *Ctx) getLuaClientIP(L *lua.LState) int {
	ret := c.user.GetClientIP()
	L.Push(lua.LString(ret))
	return 1
}
func (c *Ctx) getLuaRquestID(L *lua.LState) int {
	ret := c.user.GetRequestID()
	L.Push(lua.LString(ret))
	return 1
}
func (c *Ctx) getLuaResponseStatus(L *lua.LState) int {
	ret := c.response.getStatus()
	L.Push(lua.LNumber(ret))
	return 1
}
func (c *Ctx) getLuaResponseContent(L *lua.LState) int {
	ret := c.response.getContent()
	L.Push(lua.LString(ret))
	return 1
}
