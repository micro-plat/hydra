package creator

import "github.com/micro-plat/hydra/conf"

type ICronBinder interface {
	SetMain(*conf.CronServerConf)
	SetTasks(*conf.Tasks)
	IExtBinder
}

type CronBinder struct {
	*mainBinder
}

func newCronBinder(params map[string]string, inputs map[string]*Input) *CronBinder {
	return &CronBinder{
		mainBinder: newMainBinder(params, inputs),
	}
}
func (b *CronBinder) SetMain(c *conf.CronServerConf) {
	b.mainBinder.SetMainConf(c)
}
func (b *CronBinder) SetTasks(c *conf.Tasks) {
	b.mainBinder.SetSubConf("task", c)
}
