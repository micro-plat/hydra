package ctx

import (
	"sync"

	"github.com/micro-plat/hydra/context"
)

type funcs struct {
	tonce     sync.Once
	lonce     sync.Once
	tmplFuncs context.TFuncs
	luaFuncs  context.LuaModules
	ctx       *Ctx
}

func newTmplFunc(ctx *Ctx) *funcs {
	return &funcs{
		tmplFuncs: make(context.TFuncs),
		luaFuncs:  make(context.LuaModules),
		ctx:       ctx,
	}
}

func (r *funcs) TmplFuncs() context.TFuncs {
	r.tonce.Do(func() {
		r.tmplFuncs["get_req"] = r.ctx.request.GetString
		r.tmplFuncs["get_req_string"] = r.ctx.request.GetString
		r.tmplFuncs["get_req_int"] = r.ctx.request.GetInt
		r.tmplFuncs["get_param"] = r.ctx.request.Param
		r.tmplFuncs["get_path"] = r.ctx.request.path.GetRequestPath
		r.tmplFuncs["get_router"] = r.ctx.request.path.GetRouter
		r.tmplFuncs["get_header"] = r.ctx.request.path.GetHeader
		r.tmplFuncs["get_cookie"] = r.ctx.request.path.getCookie
		r.tmplFuncs["get_status"] = r.ctx.response.getStatus
		r.tmplFuncs["get_content"] = r.ctx.response.getContent
		r.tmplFuncs["get_client_ip"] = r.ctx.user.GetClientIP
		r.tmplFuncs["get_request_id"] = r.ctx.user.GetRequestID
	})
	return r.tmplFuncs
}
