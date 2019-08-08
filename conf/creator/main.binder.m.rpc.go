package creator

import "github.com/micro-plat/hydra/conf"

type IRPCBinder interface {
	SetMain(c *conf.RPCServerConf)
	imicroBinder
}

type RpcBinder struct {
	*microBinder
}

func NewRpcBinder(params map[string]string, inputs map[string]*Input) *RpcBinder {
	return &RpcBinder{
		microBinder: newMicroBinder(params, inputs),
	}
}
func (b *RpcBinder) SetMain(c *conf.RPCServerConf) {
	b.microBinder.SetMainConf(c)
}
