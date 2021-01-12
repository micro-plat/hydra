package pkgs

import (
	"github.com/micro-plat/hydra/global"
)

//FnInvoker 调用函数，用于指定调用本地服务的函数签名格式
type FnInvoker func(s string) *Rspns

//Invoker 本地调用配置
type Invoker struct {
	allow   bool
	service string
	addr    string
}

//NewInvoker 构建本地调用配置
func NewInvoker(p string) *Invoker {
	addr, allow := global.IsProto(p, global.ProtoInvoker, global.ProtoRPC)
	return &Invoker{
		allow:   allow,
		service: p,
		addr:    addr,
	}
}

//Allow 检查是否允许本地调用，允许则返回调用地址
func (i *Invoker) Allow() bool {
	return i.allow
}

//Invoke 调用指定的函数并检查返回结果，返回结果不包含error则认为成功
func (i *Invoker) Invoke(call FnInvoker) *Rspns {
	return call(i.service)
}

//CheckAndInvoke 检查是否可以调用服务，可以则直接调用并返回结果
func (i *Invoker) CheckAndInvoke(call FnInvoker) (bool, *Rspns) {
	if !i.Allow() {
		return false, nil
	}
	r := i.Invoke(call)
	return true, r
}
