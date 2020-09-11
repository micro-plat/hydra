package gray

import (
	"github.com/micro-plat/hydra/pkgs/lua"
	"github.com/micro-plat/lib4go/types"
)

//NeedGo2UpStream 检查当前是否需要转到上游服务器处理
func (g *Gray) NeedGo2UpStream(module lua.Modules) (bool, error) {
	vm, err := lua.New(g.Script, lua.WithMainFuncMode(), lua.WithModules(module))
	if err != nil {
		return false, err
	}
	defer vm.Shutdown()

	v, err := vm.CallByMethod(g.Go2UpStreamMethod)
	if err != nil {
		return false, err
	}
	return types.GetStringByIndex(v, 0, "false") == "true", nil
}
