package server

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
)

type MQCSub struct {
	cnf    conf.IServerConf
	queues *Loader
	server *Loader
}

func NewMQCSub(cnf conf.IServerConf) *MQCSub {
	return &MQCSub{
		cnf: cnf,
		queues: GetLoader(cnf, func(cnf conf.IServerConf) (interface{}, error) {
			return queue.GetConf(cnf)
		}),
		server: GetLoader(cnf, func(cnf conf.IServerConf) (interface{}, error) {
			return mqc.GetConf(cnf)
		}),
	}
}

//GetMQCMainConf MQC 服务器配置
func (s *MQCSub) GetMQCMainConf() (*mqc.Server, error) {
	mqcObj, err := s.server.GetConf()
	if err != nil {
		return nil, err
	}
	return mqcObj.(*mqc.Server), nil
}

//GetMQCQueueConf 获取MQC服务器的队列配置
func (s *MQCSub) GetMQCQueueConf() (*queue.Queues, error) {
	queuesObj, err := s.queues.GetConf()
	if err != nil {
		return nil, err
	}
	return queuesObj.(*queue.Queues), nil
}
