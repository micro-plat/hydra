package creator

import "github.com/micro-plat/hydra/conf"

type IWSBinder interface {
	imicroBinder
	SetMain(*conf.WSServerConf)
}

type WSBinder struct {
	*microBinder
}

func NewWSBinder(params map[string]string, inputs map[string]*Input) *WSBinder {
	return &WSBinder{
		microBinder: newMicroBinder(params, inputs),
	}
}

func (b *WSBinder) SetMain(c *conf.WSServerConf) {
	b.microBinder.SetMainConf(c)
}
