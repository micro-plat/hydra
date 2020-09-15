package apm

import (
	"context"
)

type APMContext struct {
	tracer  Tracer
	rootCtx context.Context
}

func NewCtx(ctx context.Context, tracer Tracer) *APMContext {
	return &APMContext{
		tracer:  tracer,
		rootCtx: ctx,
	}
}

func (c *APMContext) GetTracer() Tracer {
	return c.tracer

}
func (c *APMContext) GetRootCtx() context.Context {
	return c.rootCtx
}
