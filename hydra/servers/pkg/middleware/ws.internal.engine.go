package middleware

import (
	"strings"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

var wsInternalEngine *wsEngine

//InitWSInternalEngine 创建默认的WS引擎
func InitWSInternalEngine(routers ...*router.Router) {
	wsInternalEngine = newWSEngine(routers...)
}

type wsEngine struct {
	*dispatcher.Engine
}

func newWSEngine(routers ...*router.Router) *wsEngine {
	s := &wsEngine{Engine: dispatcher.New()}
	s.Engine.Use(Recovery().DispFunc(global.WS))
	s.Engine.Use(Logging().DispFunc()) //记录请求日志
	s.Engine.Use(Tag().DispFunc())
	s.Engine.Use(Trace().DispFunc()) //跟踪信息
	s.Engine.Use(Delay().DispFunc()) //
	s.Engine.Use(APIKeyAuth().DispFunc())
	s.Engine.Use(RASAuth().DispFunc())
	s.Engine.Use(JwtAuth().DispFunc())   //jwt安全认证
	s.Engine.Use(Render().DispFunc())    //响应渲染组件
	s.Engine.Use(JwtWriter().DispFunc()) //设置jwt回写
	// s.Engine.Use(s.metric.Handle().DispFunc()) //生成metric报表
	s.Engine.Use(Limit().DispFunc()) //限流处理
	s.addWSRouter(routers...)
	return s
}
func (s *wsEngine) addWSRouter(routers ...*router.Router) {
	for _, router := range routers {
		for _, method := range router.Action {
			s.Engine.Handle(strings.ToUpper(method), router.Path, ExecuteHandler(router.Service).DispFunc())
		}
	}
}
