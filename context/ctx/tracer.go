package ctx

import (
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx/internal"
	"github.com/micro-plat/lib4go/logger"
)

type tracer struct {
	*internal.Tracer
	l logger.ILogger
}

func newTracer(path string, l logger.ILogger, c app.IAPPConf) *tracer {
	return &tracer{
		Tracer: internal.Empty,
		l:      l,
	}
}

//Root 根节点
func (t *tracer) Root() context.ITraceSpan {
	return t.Tracer.Root()
}
