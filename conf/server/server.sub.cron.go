package server

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/task"
)

type CronSub struct {
	cnf  conf.IServerConf
	task *Loader
}

func NewCronSub(cnf conf.IServerConf) *CronSub {
	return &CronSub{
		cnf: cnf,
		task: GetLoader(cnf,
			func(cnf conf.IServerConf) (interface{}, error) {
				return task.GetConf(cnf)
			}),
	}
}

//GetCRONTaskConf 获取cron任务配置
func (s *CronSub) GetCRONTaskConf() (*task.Tasks, error) {
	taskObj, err := s.task.GetConf()
	if err != nil {
		return nil, err
	}
	return taskObj.(*task.Tasks), nil
}
