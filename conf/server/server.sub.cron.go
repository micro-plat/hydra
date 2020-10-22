package server

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/task"
)

type cronSub struct {
	cnf  conf.IMainConf
	task *Loader
}

func newCronSub(cnf conf.IMainConf) *cronSub {
	return &cronSub{
		cnf: cnf,
		task: GetLoader(cnf,
			func(cnf conf.IMainConf) (interface{}, error) {
				return task.GetConf(cnf)
			}),
	}
}

//GetCRONTaskConf 获取cron任务配置
func (s *cronSub) GetCRONTaskConf() (*task.Tasks, error) {
	taskObj, err := s.task.GetConf()
	if err != nil {
		return nil, err
	}
	return taskObj.(*task.Tasks), nil
}
