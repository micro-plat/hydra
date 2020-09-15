package ctx

import (
	"sync"

	"github.com/micro-plat/hydra/context"
	// "github.com/micro-plat/hydra/pkgs/lua"
)

type funcs struct {
	tonce     sync.Once
	lonce     sync.Once
	tmplFuncs context.TFuncs
	//luaFuncs  lua.Modules
	ctx *Ctx
}

func newFunc(ctx *Ctx) *funcs {
	return &funcs{
		tmplFuncs: make(context.TFuncs),
		//luaFuncs:  make(lua.Modules),
		ctx: ctx,
	}
}

func (r *funcs) TmplFuncs() context.TFuncs {
	r.tonce.Do(func() {
		r.tmplFuncs["getString"] = r.ctx.request.GetString
		r.tmplFuncs["getInt"] = r.ctx.request.GetInt
		r.tmplFuncs["getParam"] = r.ctx.request.Param
		r.tmplFuncs["getPath"] = r.ctx.request.path.GetRequestPath
		r.tmplFuncs["getRouter"] = r.ctx.request.path.GetRouter
		r.tmplFuncs["getHeader"] = r.ctx.request.path.GetHeader
		r.tmplFuncs["getCookie"] = r.ctx.request.path.getCookie
		r.tmplFuncs["getStatus"] = r.ctx.response.getStatus
		r.tmplFuncs["getContent"] = r.ctx.response.getContent
		r.tmplFuncs["getClientIP"] = r.ctx.user.GetClientIP
		r.tmplFuncs["getRequestID"] = r.ctx.user.GetRequestID
	})
	return r.tmplFuncs
}
