package tgo

//UFunc 函数
type UFunc UserFunction

//NewUFunc 构建用户函数
func NewUFunc(name string, f CallableFunc) *UserFunction {
	return &UserFunction{Name: name, Value: f}
}
