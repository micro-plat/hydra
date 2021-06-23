package errs

import (
	"errors"
	"fmt"
)

var _ IError = &Result{}

type IResult interface {
	GetResult() interface{}
	IError
}
type Result struct {
	code   int
	result interface{}
}

func (a *Result) Error() string {
	return fmt.Sprintf("%v", a.result)
}
func (a *Result) GetError() error {
	return errors.New(a.Error())
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
