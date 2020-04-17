package cron

import (
	"fmt"

	"github.com/micro-plat/hydra/servers/pkg/middleware"
)

func (s *CronServer) getProcessor() (engine *Processor, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("%v", err1)
		}
	}()
	if engine, err = NewProcessor(s.engine); err != nil {
		return nil, err
	}
	engine.Use(middleware.Logging(s.conf)) //记录请求日志
	engine.Use(middleware.Recovery())
	engine.Use(s.option.metric.Handle()) //生成metric报表
	engine.Use(middleware.NoResponse(s.conf))
	return engine, nil
}
