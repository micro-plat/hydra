package tgo

import (
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/micro-plat/lib4go/types"
)

const (

	//CodeBlockMode 代码块模式
	CodeBlockMode int = 0

	//MainFuncMode main函数调用模式
	MainFuncMode int = 1
)

//VM lua虚拟机
type VM struct {
	script   *tengo.Script
	compiled *tengo.Compiled
	modules  []*Module
	variable []*Variable
	mode     int
}

//New 构建虚拟机
func New(scope string, opts ...Option) (*VM, error) {
	vm := &VM{}

	for _, opt := range opts {
		opt(vm)
	}

	//加载脚本
	script := tengo.NewScript([]byte(scope))

	//添加变量
	for _, fn := range vm.variable {
		if err := script.Add(fn.Name(), fn.Object()); err != nil {
			return nil, err
		}

	}

	//加载模块
	modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	for _, v := range vm.modules {
		modules.AddBuiltinModule(v.name, v.Objects())
	}
	script.SetImports(modules)

	//编译脚本
	compiled, err := script.Compile()
	if err != nil {
		return nil, err
	}
	vm.compiled = compiled
	return vm, nil
}

//Run 执行脚本
func (v *VM) Run(vars ...*Variable) (types.XMap, error) {
	script := v.compiled.Clone()
	//添加变量
	for _, fn := range vars {
		if err := script.Set(fn.Name(), fn.Object()); err != nil {
			return nil, err
		}

	}
	if err := script.Run(); err != nil {
		return nil, err
	}
	return var2Mpa(script), nil
}
