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
		cnf:  cnf,
		task: GetLoader(cnf, task.ConfHandler(task.GetConf).Handle),
	}
}

//GetCRONTaskConf 获取cron任务配置
func (s *cronSub) GetCRONTaskConf() *task.Tasks {
	return s.task.GetConf().(*task.Tasks)
}
