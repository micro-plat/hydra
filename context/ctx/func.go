package ctx

import (
	"sync"

	"github.com/micro-plat/hydra/context"
)

type tmplFuncs struct {
	once  sync.Once
	funcs map[string]interface{}
	ctx   *Ctx
}

func newTmplFunc(ctx *Ctx) *tmplFuncs {
	return &tmplFuncs{
		funcs: make(map[string]interface{}),
		ctx:   ctx,
	}
}

func (r *tmplFuncs) Instance() context.TFuncs {
	r.once.Do(func() {
		r.funcs["get_req"] = r.ctx.request.GetString
		r.funcs["get_req_string"] = r.ctx.request.GetString
		r.funcs["get_req_int"] = r.ctx.request.GetInt
		r.funcs["get_param"] = r.ctx.request.Param
		r.funcs["get_path"] = r.ctx.request.path.GetRequestPath
		r.funcs["get_router"] = r.ctx.request.path.GetRouter
		r.funcs["get_header"] = r.ctx.request.path.GetHeader
		r.funcs["get_cookie"] = r.ctx.request.path.getCookie
		r.funcs["get_status"] = r.ctx.response.getStatus
		r.funcs["get_content"] = r.ctx.response.getContent
		r.funcs["get_client_ip"] = r.ctx.user.GetClientIP
		r.funcs["get_request_id"] = r.ctx.user.GetRequestID
	})
	return r.funcs
}
