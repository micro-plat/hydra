package lua

import (
	"fmt"
	"strings"
	"sync"

	"github.com/micro-plat/lib4go/types"
	"github.com/yuin/gopher-lua/parse"
	luar "layeh.com/gopher-luar"

	lua "github.com/yuin/gopher-lua"
)

const (

	//CodeBlockMode 代码块模式
	CodeBlockMode int = 0

	//MainFuncMode main函数调用模式
	MainFuncMode int = 1
)

//VM lua虚拟机
type VM struct {
	s         *lua.LState
	funcProto *lua.LFunction
	lock      sync.Mutex
	main      lua.LValue
	mode      int
}

//New 构建LUA虚拟机
func New(script string, opts ...Option) (*VM, error) {
	s := lua.NewState(lua.Options{IncludeGoStackTrace: true})
	vm := &VM{s: s}
	for _, opt := range opts {
		opt(vm)
	}
	if err := vm.s.DoString(script); err != nil {
		return nil, err
	}
	switch vm.mode {
	case MainFuncMode:
		if err := vm.s.DoString(script); err != nil {
			return nil, err
		}
		if vm.main = vm.s.GetGlobal("main"); vm.main == nil {
			return nil, fmt.Errorf("脚本中未包含main函数")
		}
	case CodeBlockMode:
		proto, err := vm.load(script)
		if err != nil {
			return nil, fmt.Errorf("脚本加载失败%w", err)
		}
		vm.main = vm.s.NewFunctionFromProto(proto)
	}
	return vm, nil
}

//Call 执行脚本
func (v *VM) Call(input ...interface{}) (string, error) {
	rets, err := v.Calls(input...)
	return types.GetStringByIndex(rets, 0, ""), err
}

//Calls 执行脚本
func (v *VM) Calls(input ...interface{}) ([]string, error) {
	v.lock.Lock()
	defer v.lock.Unlock()
	block := lua.P{
		Fn:      v.main,
		NRet:    1,
		Protect: true,
	}
	lvs := make([]lua.LValue, 0, len(input))
	for _, value := range input {
		lvs = append(lvs, luar.New(v.s, value))
	}
	if err := v.s.CallByParam(block, lvs...); err != nil {
		return nil, fmt.Errorf("调用main失败 %w", err)
	}
	return v.getRets()
}

//CallByMethod 调用脚本函数
func (v *VM) CallByMethod(method string, input ...interface{}) ([]string, error) {
	v.lock.Lock()
	defer v.lock.Unlock()
	block := lua.P{
		Fn:      v.s.GetGlobal(method),
		NRet:    1,
		Protect: true,
	}
	lvs := make([]lua.LValue, 0, len(input))
	for _, value := range input {
		lvs = append(lvs, luar.New(v.s, value))
	}
	if err := v.s.CallByParam(block, lvs...); err != nil {
		return nil, fmt.Errorf("调用%s失败 %w", method, err)
	}
	return v.getRets()
}

//GetValue 获取指定变量的名称
func (v *VM) GetValue(name string) string {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.s.GetGlobal(name).String()
}

//Shutdown 关闭虚拟机
func (v *VM) Shutdown() {
	if v.s != nil {
		v.s.Close()
	}
}

//Load 加载脚本
func (v *VM) load(source string) (*lua.FunctionProto, error) {
	reader := strings.NewReader(source)
	chunk, err := parse.Parse(reader, source)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, source)
	if err != nil {
		return nil, err
	}
	return proto, nil
}
func (v *VM) getRets() ([]string, error) {
	top := v.s.GetTop()
	if top <= 0 {
		return nil, nil
	}
	rets := make([]string, 0, top)
	for i := 1; i <= top; i++ {
		rets = append(rets, v.s.Get(i).String())
	}
	v.s.Pop(top)
	return rets, nil
}
