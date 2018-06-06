package context

import (
	"errors"
	"fmt"
)

var ERR_DataNotExist = errors.New("查询的数据不存在")

type IError interface {
	GetError() error
	GetCode() int
	Error() string
	IsFatal() bool
}
type Error struct {
	code    int
	isFatal bool
	error
}

func (a *Error) GetCode() int {
	return a.code
}
func (a *Error) GetError() error {
	return a
}
func (a *Error) IsFatal() bool {
	return a.isFatal
}

//NewIgnoreError 当前一个可忽略的错误
func NewIgnoreError(code int, err interface{}) *Error {
	ex := NewError(code, err)
	ex.isFatal = false
	return ex
}

//NewError 创建一个致命的错误
func NewError(code int, err interface{}) *Error {
	r := &Error{code: code, isFatal: true}
	switch v := err.(type) {
	case string:
		r.error = errors.New(v)
	case error:
		r.error = v
	case IError:
		r.error = v.GetError()
	default:
		r.error = errors.New(fmt.Sprint(err))
	}
	return r
}
