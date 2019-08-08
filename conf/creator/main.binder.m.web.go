package creator

import "github.com/micro-plat/hydra/conf"

type IWebBinder interface {
	imicroBinder
	SetStatic(*conf.Static)
	SetMain(*conf.WebServerConf)
}

type WebBinder struct {
	*microBinder
}

func NewWebBinder(params map[string]string, inputs map[string]*Input) *WebBinder {
	return &WebBinder{
		microBinder: newMicroBinder(params, inputs),
	}
}

func (b *WebBinder) SetMain(c *conf.WebServerConf) {
	b.microBinder.SetMainConf(c)
}
func (b *WebBinder) SetStatic(c *conf.Static) {
	b.microBinder.SetSubConf("static", c)
}
