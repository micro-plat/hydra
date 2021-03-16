package errs

import (
	"errors"
	"fmt"

	"github.com/micro-plat/lib4go/types"
)

//ErrNotExist 不存在
var ErrNotExist = errors.New("不存在")

//IError 包含错误码的error
type IError interface {
	Error() string
	GetError() error
	GetCode() int
}

//Error 错误信息
type Error struct {
	code int
	error
}

//GetCode 获取错误码
func (a *Error) GetCode() int {
	return a.code
}

//GetError 获取错误信息
func (a *Error) GetError() error {
	return a
}

//String 格式化错误信息
func (a *Error) String() string {
	return fmt.Sprintf("%d %s", a.code, a.Error())
}

func (a *Error) Is(e error) bool {
	return errors.Is(a.error, e)
}
func (a *Error) As(target interface{}) bool {
	return errors.As(a.error, target)
}

//NewErrorf 创建错误对象
func NewErrorf(code int, f string, args ...interface{}) *Error {
	return NewError(code, fmt.Errorf(f, args...))
}

//NewError 创建错误对象
func NewError(code int, err interface{}) *Error {
	r := &Error{code: code}
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

//GetCode 获取错误码
func GetCode(err interface{}, def ...int) int {
	switch v := err.(type) {
	case IError:
		return v.GetCode()
	default:
		return types.GetIntByIndex(def, 0, 0)
	}
}

//GetError 获取错误，当不包含错误时返回空
func GetError(r interface{}) IError {
	switch v := r.(type) {
	case IError:
		return v
	case error:
		return NewError(400, v)
	default:
		return nil
	}
}
