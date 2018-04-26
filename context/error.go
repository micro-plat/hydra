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
}
type Error struct {
	code int
	error
}

func (a *Error) GetCode() int {
	return a.code
}
func (a *Error) GetError() error {
	return a
}
func NewError(code int, err interface{}) IError {
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
