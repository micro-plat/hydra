package nfs

import (
	"time"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/logger"
)

type reqLog struct {
	start time.Time
	ctx   context.IContext
}

func req(ctx context.IContext, input ...interface{}) reqLog {
	r := reqLog{start: time.Now(), ctx: ctx}
	p := make([]interface{}, 0, len(input)+1)
	p = append(p, "nfs.request:")
	p = append(p, ctx.Request().Path())
	for _, v := range ctx.Request().Keys() {
		p = append(p, v)
	}
	p = append(p, input...)
	// ctx.Log().Debug(p...)
	return r
}

func (r reqLog) rspns(input ...interface{}) {
	p := make([]interface{}, 0, len(input)+1)
	p = append(p, "nfs.response:")
	p = append(p, r.ctx.Request().Path())
	p = append(p, input...)
	s, _, _ := r.ctx.Response().GetFinalResponse()
	p = append(p, s)
	p = append(p, time.Since(r.start))
	// r.ctx.Log().Debug(p...)
}

type log struct {
	start time.Time
	log   *logger.Logger
}

func start(input ...interface{}) log {
	l := log{
		start: time.Now(),
		log:   logger.New("nfs"),
	}
	p := make([]interface{}, 0, len(input)+1)
	p = append(p, "event.start > ")
	p = append(p, input...)
	l.log.Debug(p...)
	return l
}
func (l log) end(input ...interface{}) {
	p := make([]interface{}, 0, len(input)+2)
	p = append(p, "event.end > ")
	p = append(p, input...)
	p = append(p, time.Since(l.start))
	l.log.Debug(p...)
}
