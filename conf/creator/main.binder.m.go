package creator

import "github.com/micro-plat/hydra/conf"

type imicroBinder interface {
	SetAuthes(*conf.Authes)
	SetRouters(*conf.Routers)
	SetCircuitBreaker(*conf.CircuitBreaker)
	SetHeaders(conf.Headers)
	IExtBinder
}

type microBinder struct {
	*mainBinder
}

func newMicroBinder(params map[string]string, inputs map[string]*Input) *microBinder {
	return &microBinder{
		mainBinder: newMainBinder(params, inputs),
	}
}

func (b *microBinder) SetAuthes(a *conf.Authes) {
	b.mainBinder.SetSubConf("auth", a)
}
func (b *microBinder) SetRouters(r *conf.Routers) {
	b.mainBinder.SetSubConf("router", r)
}
func (b *microBinder) SetCircuitBreaker(c *conf.CircuitBreaker) {
	b.mainBinder.SetSubConf("circuit", c)
}
func (b *microBinder) SetHeaders(h conf.Headers) {
	b.mainBinder.SetSubConf("header", h)
}
