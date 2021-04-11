package nfs

import (
	"time"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/logger"
)

func reqLog(ctx context.IContext, input ...interface{}) {
	p := make([]interface{}, 0, len(input)+1)
	p = append(p, "nfs.request:")
	p = append(p, ctx.Request().Path())
	for _, v := range ctx.Request().Keys() {
		p = append(p, v)
	}
	p = append(p, input...)
	ctx.Log().Debug(p...)
}

func rspnsLog(ctx context.IContext, input ...interface{}) {
	p := make([]interface{}, 0, len(input)+1)
	p = append(p, "nfs.response:")
	p = append(p, ctx.Request().Path())
	p = append(p, input...)
	ctx.Log().Debug(p...)
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
	l.log.Debug(input)
	return l
}
func (l log) end(input ...interface{}) {
	p := make([]interface{}, 0, len(input)+1)
	p = append(p, input...)
	p = append(p, time.Since(l.start))
	l.log.Debug(p...)
}
