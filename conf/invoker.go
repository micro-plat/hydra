package conf

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/errs"
)

//Invoker 本地调用配置
type Invoker struct {
	allow bool
	addr  string
}

//NewInvoker 构建本地调用配置
func NewInvoker(p string) *Invoker {
	addr, allow := global.IsProto(p, global.ProtoInvoker)
	return &Invoker{
		allow: allow,
		addr:  addr,
	}
}

//Allow 检查是否允许本地调用，允许则返回调用地址
func (i *Invoker) Allow() bool {
	return i.allow
}

//Invoke 调用指定的函数并检查返回结果，返回结果不包含error则认为成功
func (i *Invoker) Invoke(call func(s string) interface{}) error {
	result := call(i.addr)
	if err := errs.GetError(result); err != nil {
		return err
	}
	return nil
}
