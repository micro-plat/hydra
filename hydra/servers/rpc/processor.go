package rpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/jsons"
)

//Processor cron管理程序，用于管理多个任务的执行，暂停，恢复，动态添加，移除
type Processor struct {
	//*dispatcher.Engine
	done          bool
	closeChan     chan struct{}
	metric        *middleware.Metric
	adapterEngine *adapter.Engine
}

//NewProcessor 创建processor
func NewProcessor(routers ...*router.Router) (p *Processor) {
	p = &Processor{
		closeChan: make(chan struct{}),
		metric:    middleware.NewMetric(),
	}
	p.adapterEngine = adapter.New(adapter.NewEngineWrapperDisp(dispatcher.New(), RPC))

	p.adapterEngine.Use(middleware.Recovery())
	p.adapterEngine.Use(middleware.Logging())
	p.adapterEngine.Use(middleware.Recovery())
	p.adapterEngine.Use(p.metric.Handle())

	p.adapterEngine.Use(middleware.Trace()) //跟踪信息
	p.adapterEngine.Use(middleware.Delay())
	p.adapterEngine.Use(middlewares...)

	p.addRouter(routers...)
	return p
}

func (s *Processor) addRouter(routers ...*router.Router) {
	adapterRouters := make([]adapter.IRouter, len(routers))
	for i := range routers {
		adapterRouters[i] = routers[i]
	}
	s.adapterEngine.Handle(adapterRouters...)
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
	w, err := s.adapterEngine.HandleRequest(req)
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

//Close 关闭处理程序
func (s *Processor) Close() {
	s.metric.Stop()
}
