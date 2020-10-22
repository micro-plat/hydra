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
		cnf: cnf,
		queues: GetLoader(cnf, func(cnf conf.IMainConf) (interface{}, error) {
			return queue.GetConf(cnf)
		}),
		server: GetLoader(cnf, func(cnf conf.IMainConf) (interface{}, error) {
			return mqc.GetConf(cnf)
		}),
	}
}

//GetMQCMainConf MQC 服务器配置
func (s *mqcSub) GetMQCMainConf() (*mqc.Server, error) {
	mqcObj, err := s.server.GetConf()
	if err != nil {
		return nil, err
	}
	return mqcObj.(*mqc.Server), nil
}

//GetMQCQueueConf 获取MQC服务器的队列配置
func (s *mqcSub) GetMQCQueueConf() (*queue.Queues, error) {
	queuesObj, err := s.queues.GetConf()
	if err != nil {
		return nil, err
	}
	return queuesObj.(*queue.Queues), nil
}
