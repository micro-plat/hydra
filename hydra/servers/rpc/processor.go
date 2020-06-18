package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/jsons"
)

//Processor cron管理程序，用于管理多个任务的执行，暂停，恢复，动态添加，移除
type Processor struct {
	*dispatcher.Engine
	done      bool
	closeChan chan struct{}
}

//NewProcessor 创建processor
func NewProcessor(routers ...*router.Router) (p *Processor) {
	p = &Processor{
		closeChan: make(chan struct{}),
	}
	p.Engine = dispatcher.New()
	p.Engine.Use(middleware.Recovery().DispFunc(RPC))
	p.Engine.Use(middleware.Logging().DispFunc())
	p.Engine.Use(middleware.BlackList().DispFunc()) //黑名单控制
	p.Engine.Use(middleware.WhiteList().DispFunc()) //白名单控制
	p.Engine.Use(middleware.Trace().DispFunc())     //跟踪信息
	p.Engine.Use(middleware.Delay().DispFunc())     //
	p.Engine.Use(middleware.Options().DispFunc())   //处理option响应
	p.Engine.Use(middleware.Static().DispFunc())    //处理静态文件
	p.Engine.Use(middleware.Header().DispFunc())    //设置请求头
	p.Engine.Use(middleware.BasicAuth().DispFunc()) //
	p.Engine.Use(middleware.APIKeyAuth().DispFunc())
	p.Engine.Use(middleware.RASAuth().DispFunc())
	p.Engine.Use(middleware.JwtAuth().DispFunc())   //jwt安全认证
	p.Engine.Use(middleware.Render().DispFunc())    //响应渲染组件
	p.Engine.Use(middleware.JwtWriter().DispFunc()) //设置jwt回写
	p.addRouter(routers...)
	return p
}

func (s *Processor) addRouter(routers ...*router.Router) {
	for _, router := range routers {
		for _, method := range router.Action {
			s.Engine.Handle(strings.ToUpper(method), router.Path, middleware.ExecuteHandler(router.Service).DispFunc())
		}
	}
}

//Request 处理业务请求
func (s *Processor) Request(context context.Context, request *pb.RequestContext) (p *pb.ResponseContext, err error) {
	req, err := NewRequest(request)
	if err != nil {
		return nil, fmt.Errorf("输入参数有误:%w", err)
	}
	w, err := s.Engine.HandleRequest(req)
	if err != nil {
		return nil, err
	}
	p = &pb.ResponseContext{}
	p.Status = int32(w.Status())
	p.Result = string(w.Data())
	h, err := jsons.Marshal(w.Header())
	if err != nil {
		return p, err
	}
	p.Header = string(h)
	return p, nil
}
