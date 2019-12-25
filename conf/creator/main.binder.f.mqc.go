package creator

import "github.com/micro-plat/hydra/conf"

type IMQCBinder interface {
	SetMain(*conf.MQCServerConf)
	SetServer(conf.QueueConf)
	SetQueues(*conf.Queues)
	IExtBinder
}

type MQCBinder struct {
	*mainBinder
}

func NewMQCBinder(params map[string]string, inputs map[string]*Input) *MQCBinder {
	return &MQCBinder{
		mainBinder: newMainBinder(params, inputs),
	}
}

func (b *MQCBinder) SetMain(c *conf.MQCServerConf) {
	b.mainBinder.SetMainConf(c)
}
func (b *MQCBinder) SetServer(c conf.QueueConf) {
	b.mainBinder.SetSubConf("server", c)
}

func (b *MQCBinder) SetQueues(c *conf.Queues) {
	b.mainBinder.SetSubConf("queue", c)
}
