package mqc

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
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
	s.AddRouters(engine)
	if err = engine.Consumes(); err != nil {
		return nil, err
	}
	return engine, nil
}
func (s *MqcServer) AddRouters(p *Processor) {
	for _, r := range p.queues {
		if _, ok := p.handles[r.Name]; !ok {
			s.Logger.Debugf("[订阅 队列(%s)消息]", r.Queue)
			handler := r.Handler.(dispatcher.HandlerFunc)
			p.handles[r.Name] = handler
			path := fmt.Sprintf("/%s", strings.TrimPrefix(r.Name, "/"))
			if !p.Dispatcher.Find(path) {
				p.Dispatcher.Handle(strings.ToUpper("GET"), path, handler)
			}

		}
	}
}
