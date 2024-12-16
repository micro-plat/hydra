package ws

import (
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

var wsInternalEngine *wsEngine
var metric = middleware.NewMetric()

//InitWSEngine 创建默认的WS引擎
func InitWSEngine(routers ...*router.Router) {
	wsInternalEngine = newWSEngine(routers...)
}

type wsEngine struct {
	//*dispatcher.Engine
	metric *middleware.Metric
	engine *adapter.DispatcherEngine
}

func newWSEngine(routers ...*router.Router) *wsEngine {
	s := &wsEngine{
		metric: metric,
	}
	s.engine = adapter.NewDispatcherEngine(global.WS)

	s.engine.Use(middleware.Recovery())
	s.engine.Use(middleware.Logging()) //记录请求日志
	s.engine.Use(middleware.Recovery())
	s.engine.Use(middleware.Tag())
	s.engine.Use(middleware.Trace()) //跟踪信息
	s.engine.Use(middleware.Limit()) //限流处理
	s.engine.Use(middleware.Delay()) //
	s.engine.Use(middleware.APIKeyAuth())
	s.engine.Use(middleware.RASAuth())
	s.engine.Use(middleware.JwtAuth())   //jwt安全认证
	s.engine.Use(middleware.Render())    //响应渲染组件
	s.engine.Use(middleware.JwtWriter()) //设置jwt回写
	s.engine.Use(middlewares...)
	s.engine.Use(s.metric.Handle()) //生成metric报表

	s.addWSRouter(routers...)
	return s
}
func (s *wsEngine) addWSRouter(routers ...*router.Router) {
	for _, r := range routers {
		for _, action := range r.Action {
			s.engine.Handle(action, r.Path, middleware.ExecuteHandler())
		}
	}
}
