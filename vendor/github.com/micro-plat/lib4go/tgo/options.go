package tgo

//Option 配置选项
type Option func(*VM)

//WithModule 添加module
func WithModule(m ...*Module) Option {
	return func(vm *VM) {
		vm.modules = append(vm.modules, m...)
	}
}

//WithVariable 添加变量
func WithVariable(m ...*Variable) Option {
	return func(vm *VM) {
		vm.variable = append(vm.variable, m...)
	}
}
