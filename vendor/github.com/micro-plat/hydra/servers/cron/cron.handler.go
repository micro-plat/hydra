package cron

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/middleware"
)

func (s *CronServer) getProcessor(redisSetting string, tasks []*conf.Task) (engine *Processor, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("%v", err1)
		}
	}()
	engine, err = NewProcessor(redisSetting, fmt.Sprintf("%s:tasks:%s", s.conf.Name, "%s"))
	if err != nil {
		return nil, err
	}
	engine.Use(middleware.Logging(s.conf)) //记录请求日志
	engine.Use(middleware.Recovery())
	engine.Use(s.option.metric.Handle()) //生成metric报表
	engine.Use(middleware.NoResponse(s.conf))
	for _, task := range tasks {
		ct, err := newCronTask(task)
		if err != nil {
			return nil, err
		}
		if _, _, err = engine.Add(ct, true); err != nil {
			return nil, err
		}

	}
	return engine, nil
}
