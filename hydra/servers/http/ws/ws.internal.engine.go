package ws

import (
	"strings"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

var wsInternalEngine *wsEngine
var metric = middleware.NewMetric()

//InitWSEngine 创建默认的WS引擎
func InitWSEngine(routers ...*router.Router) {
	wsInternalEngine = newWSEngine(routers...)
}

type wsEngine struct {
	*dispatcher.Engine
	metric *middleware.Metric
}

func newWSEngine(routers ...*router.Router) *wsEngine {
	s := &wsEngine{Engine: dispatcher.New(), metric: metric}
	s.Engine.Use(middleware.Recovery().DispFunc(global.WS))
	s.Engine.Use(middleware.Logging().DispFunc()) //记录请求日志
	s.Engine.Use(middleware.Recovery().DispFunc())
	s.Engine.Use(middleware.Tag().DispFunc())
	s.Engine.Use(middleware.Trace().DispFunc()) //跟踪信息
	s.Engine.Use(middleware.Limit().DispFunc()) //限流处理
	s.Engine.Use(middleware.Delay().DispFunc()) //
	s.Engine.Use(middleware.APIKeyAuth().DispFunc())
	s.Engine.Use(middleware.RASAuth().DispFunc())
	s.Engine.Use(middleware.JwtAuth().DispFunc())   //jwt安全认证
	s.Engine.Use(middleware.Render().DispFunc())    //响应渲染组件
	s.Engine.Use(middleware.JwtWriter().DispFunc()) //设置jwt回写
	s.Engine.Use(middlewares.DispFunc()...)
	s.Engine.Use(s.metric.Handle().DispFunc()) //生成metric报表

	s.addWSRouter(routers...)
	return s
}
func (s *wsEngine) addWSRouter(routers ...*router.Router) {
	for _, router := range routers {
		for _, method := range router.Action {
			s.Engine.Handle(strings.ToUpper(method), router.Path, middleware.ExecuteHandler(router.Service).DispFunc())
		}
	}
}
