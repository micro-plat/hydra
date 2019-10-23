package creator

import "github.com/micro-plat/hydra/conf"

type IApiBinder interface {
	imicroBinder
	SetStatic(*conf.Static)
	SetMain(*conf.APIServerConf)
	SetResponse(*conf.Response)
	SetCrossDomain()
}

type ApiBinder struct {
	*microBinder
}

func NewApiBinder(params map[string]string, inputs map[string]*Input) *ApiBinder {
	return &ApiBinder{
		microBinder: newMicroBinder(params, inputs),
	}
}

func (b *ApiBinder) SetMain(c *conf.APIServerConf) {
	b.microBinder.SetMainConf(c)
}
func (b *ApiBinder) SetStatic(c *conf.Static) {
	b.microBinder.SetSubConf("static", c)
}
func (b *ApiBinder) SetCrossDomain() {
	b.microBinder.SetHeaders(conf.NewHeader().WithCrossDomain())
}
func (b *ApiBinder) SetResponse(h *conf.Response) {
	b.mainBinder.SetSubConf("response", h)
}
