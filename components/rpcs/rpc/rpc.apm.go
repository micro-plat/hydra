package rpc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/micro-plat/hydra/context/apm"

	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"

	r "github.com/micro-plat/hydra/context"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func (c *Client) clientRequest(ctx context.Context, o *requestOption, form map[string]interface{}) (response *pb.ResponseContext, err error) {

	h, err := o.getData(o.headers)
	if err != nil {
		return nil, err
	}
	f, err := o.getData(form)
	if err != nil {
		return nil, err
	}
	remotecallback := func() (*pb.ResponseContext, error) {
		fmt.Println(o.method, o.service, string(h), string(f))
		return c.client.Request(ctx,
			&pb.RequestContext{
				Method:  o.method,
				Service: o.service,
				Header:  string(h),
				Input:   string(f),
			},
			grpc.FailFast(o.failFast))
	}

	hydractx := r.Current()
	apmCtx := hydractx.APMContext()
	if apmCtx == nil {
		return remotecallback()
	}
	apmConf := hydractx.ServerConf().GetAPMConf()
	if apmConf.Disable {
		return remotecallback()
	}
	response, err = c.execAPMRequest(apmCtx, o, remotecallback)
	return
}

func (c *Client) execAPMRequest(apmCtx r.IAPMContext, o *requestOption, callback func() (*pb.ResponseContext, error)) (res *pb.ResponseContext, err error) {
	tracer := apmCtx.GetTracer()
	rootCtx := apmCtx.GetRootCtx()

	span, err := tracer.CreateExitSpan(rootCtx, getOperationName(o), c.conn.Target(), func(header string) error {
		o.headers[apm.Header] = header
		return nil
	})
	if err != nil {
		return callback()
	}
	defer span.End()
	span.SetComponent(apm.ComponentIDGORpcClient)
	span.Tag("X-Request-Id", o.headers["X-Request-Id"])

	span.Tag(apm.TagRPCMethod, o.method) //span.Tag(apm.TagHTTPMethod, req.Method)
	span.Tag(apm.TagURL, c.address)      //span.Tag(apm.TagURL, req.URL.String())

	span.SetSpanLayer(apm.SpanLayer_RPCFramework)
	res, err = callback()
	if err != nil {
		span.Error(time.Now(), err.Error())
		return
	}

	span.Tag(apm.TagStatusCode, fmt.Sprintf("%d", res.Status))
	if res.Status >= http.StatusBadRequest {
		span.Error(time.Now(), "Errors on handling client")
	}

	return res, nil
}

func getOperationName(r *requestOption) string {
	return r.service
}
