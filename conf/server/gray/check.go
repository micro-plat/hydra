package gray

import (
	"fmt"

	"github.com/micro-plat/hydra/pkgs/lua"
	"github.com/micro-plat/lib4go/types"
)

// var modules = lua.Modules{
// 	"request": map[string]lua.LGFunction{
// 		"getClientIP": func(ls *lua.LState) int {
// 			ls.Push(lua.LString("abc"))
// 			return 1
// 		},
// 	},
// }

//NeedGo2UpStream 检查当前是否需要转到上游服务器处理
func (g *Gray) NeedGo2UpStream(module lua.Modules) (bool, error) {
	vm, err := lua.New(g.Script, lua.WithMainFuncMode(), lua.WithModules(module))
	if err != nil {
		return false, err
	}
	defer vm.Shutdown()

	v, err := vm.CallByMethod(g.go2UpStreamMethod)
	if err != nil {
		return false, err
	}
	fmt.Println("check:", types.GetStringByIndex(v, 0, "false"))
	return types.GetStringByIndex(v, 0, "false") == "true", fmt.Errorf(types.GetStringByIndex(v, 0, "false"))
}
