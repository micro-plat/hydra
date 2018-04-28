package rpc

import (
	"golang.org/x/net/context"

	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/servers/rpc/pb"
	"github.com/micro-plat/lib4go/jsons"
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

	response, err := r.Dispatcher.HandleRequest(&Request{RequestContext: request})
	if err != nil {
		return
	}
	p = &pb.ResponseContext{}
	p.Status = int32(response.Status())
	p.Result = string(response.Data())
	h, err := jsons.Marshal(response.Header())
	if err != nil {
		return p, err
	}
	p.Header = string(h)
	return p, nil
}
