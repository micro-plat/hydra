package context

import "fmt"

type IResult interface {
	GetResult() interface{}
	GetCode() int
}
type Result struct {
	code   int
	result interface{}
}

func (a *Result) GetCode() int {
	return a.code
}
func (a *Result) GetResult() interface{} {
	return a.result
}

//NewResultf 创建带状态码的返回对象
func NewResultf(code int, f string, args ...interface{}) *Result {
	return NewResult(code, fmt.Sprintf(f, args...))
}

//NewResult 创建
func NewResult(code int, content interface{}) *Result {
	return &Result{code: code, result: content}
}
