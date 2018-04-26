package rpc

import (
	"golang.org/x/net/context"

	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/servers/rpc/pb"
)

type Processor struct {
	*dispatcher.Dispatcher
}

func NewProcessor() *Processor {
	return &Processor{
		Dispatcher: dispatcher.New(),
	}
}
func (r *Processor) Request(context context.Context, request *pb.RequestContext) (p *pb.ResponseContext, err error) {
	if request.Header == nil {
		request.Header = make(map[string]string)
	}
	response, err := r.Dispatcher.HandleRequest(request)
	if err != nil {
		return
	}
	p = &pb.ResponseContext{}
	p.Status = int32(response.Status())
	p.Result = string(response.Data())
	p.Header = make(map[string]string)
	for k, v := range response.Header() {
		p.Header[k] = v[0]
	}
	return p, nil
}
