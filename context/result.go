package context

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

//NewResult 创建
func NewResult(code int, content interface{}) *Result {
	return &Result{code: code, result: content}
}
