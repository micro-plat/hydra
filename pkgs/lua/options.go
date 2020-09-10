package lua

import (
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

//Option 配置选项
type Option func(*VM)

//With 添加全局变量
func With(name string, value interface{}) Option {
	return func(L *VM) {
		L.s.SetGlobal(name, luar.New(L.s, value))
	}
}

//WithKV 批量添加一组变量
func WithKV(kv map[string]interface{}) Option {
	return func(L *VM) {
		for k, v := range kv {
			With(k, v)(L)
		}
	}
}

//WithType 添加全局类型
func WithType(tpName string, value interface{}) Option {
	return func(L *VM) {
		L.s.SetGlobal(tpName, luar.NewType(L.s, value))

	}
}

//WithMT 将对象转换为LuaTable
func WithMT(tpName string, value interface{}) Option {
	return func(L *VM) {
		L.s.SetGlobal(tpName, luar.MT(L.s, value))
	}
}

//WithModule 添加模块
func WithModule(module string, exports map[string]lua.LGFunction) Option {
	return func(L *VM) {
		L.s.PreloadModule(module, func(L *lua.LState) int {
			mod := L.SetFuncs(L.NewTable(), exports)
			L.Push(mod)
			return 1
		})
	}
}

//WithModules 添加多个模块
func WithModules(modules map[string]map[string]lua.LGFunction) Option {
	return func(L *VM) {
		for module, exports := range modules {
			L.s.PreloadModule(module, func(L *lua.LState) int {
				mod := L.SetFuncs(L.NewTable(), exports)
				L.Push(mod)
				return 1
			})
		}
	}
}

//WithCodeBlockMode 使用代码块模式，预编译脚本调用call时执行脚本
func WithCodeBlockMode() Option {
	return func(L *VM) {
		L.mode = CodeBlockMode
	}
}

//WithMainFuncMode 使用main函数调用模式，启动时执行脚本，调用call时执行main函数
func WithMainFuncMode() Option {
	return func(L *VM) {
		L.mode = MainFuncMode
	}
}
