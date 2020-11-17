package rpc

import (
	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func (c *Client) clientRequest(ctx context.Context, o *requestOption, form string) (response *pb.ResponseContext, err error) {

	h, err := o.getData(o.headers)
	if err != nil {
		return nil, err
	}

	return c.client.Request(ctx,
		&pb.RequestContext{
			Method:  o.method,
			Service: o.service,
			Header:  string(h),
			Input:   form,
		},
		grpc.FailFast(o.failFast))

}
