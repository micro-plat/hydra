package rpc

import (
	"context"
	"fmt"
	"net/http"
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
	p.Engine.Use(middleware.Recovery().DispFunc())
	p.Engine.Use(middleware.DispServiceExistsCheck(p.Engine).DispFunc())

	p.Engine.Use(middleware.Trace().DispFunc()) //跟踪信息
	p.Engine.Use(middleware.Delay().DispFunc())
	middleware.AddMiddlewareHook(rpcmiddlewares, func(item middleware.Handler) {
		p.Engine.Use(item.DispFunc())
	})
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

	//转换输入参数
	req, err := NewRequest(request)
	if err != nil {
		p = &pb.ResponseContext{}
		p.Status = int32(http.StatusNotAcceptable)
		p.Result = fmt.Sprintf("输入参数有误:%v", err)
		return p, nil
	}

	//发起本地处理
	w, err := s.Engine.HandleRequest(req)
	if err != nil {
		p = &pb.ResponseContext{}
		p.Status = int32(http.StatusInternalServerError)
		p.Result = fmt.Sprintf("处理请求有误%s", err.Error())
		return p, nil
	}

	//处理响应内容
	p = &pb.ResponseContext{}
	p.Status = int32(w.Status())
	p.Result = string(w.Data())
	h, err := jsons.Marshal(w.Header())
	if err != nil {
		p = &pb.ResponseContext{}
		p.Status = int32(http.StatusInternalServerError)
		p.Result = fmt.Sprintf("输换响应头失败 %s", err.Error())
		return p, nil
	}
	p.Header = string(h)
	return p, nil
}
