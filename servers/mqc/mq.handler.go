package mqc

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/middleware"
)

func (s *MqcServer) getProcessor(addr string, raw string, queues []*conf.Queue) (engine *Processor, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("%v", err1)
		}
	}()
	if engine, err = NewProcessor(addr, raw, queues); err != nil {
		return nil, err
	}
	engine.Use(middleware.Logging(s.conf)) //记录请求日志
	engine.Use(middleware.Recovery())
	engine.Use(s.option.metric.Handle()) //生成metric报表
	engine.Use(middleware.NoResponse(s.conf))
	engine.AddRouters()
	if err = engine.Consumes(); err != nil {
		return nil, err
	}
	return engine, nil
}
