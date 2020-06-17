package server

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
)

type mqcSub struct {
	cnf    conf.IMainConf
	queues *Loader
	server *Loader
}

func newMQCSub(cnf conf.IMainConf) *mqcSub {
	return &mqcSub{
		cnf:    cnf,
		queues: GetLoader(cnf, queue.ConfHandler(queue.GetConf).Handle),
		server: GetLoader(cnf, mqc.ConfHandler(mqc.GetConf).Handle),
	}
}

//GetMQCMainConf MQC 服务器配置
func (s *mqcSub) GetMQCMainConf() *mqc.Server {
	return s.server.GetConf().(*mqc.Server)
}

//GetMQCQueueConf 获取MQC服务器的队列配置
func (s *mqcSub) GetMQCQueueConf() *queue.Queues {
	return s.queues.GetConf().(*queue.Queues)
}
