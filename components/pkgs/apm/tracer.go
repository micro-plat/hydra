package apm

import (
	"context"
	"time"
)

type TracerOption func(t Tracer)

type Tracer interface {
	CreateLocalSpan(ctx context.Context, opts ...SpanOption) (s Span, c context.Context, err error)
	CreateEntrySpan(ctx context.Context, operationName string, extractor Extractor, opts ...SpanOption) (s Span, nCtx context.Context, err error)
	CreateExitSpan(ctx context.Context, operationName string, peer string, injector Injector, opts ...SpanOption) (Span, error)
	GetRealTracer() interface{}
	SetReporter(reporter Reporter)
	SetInstance(instance string)
}

// Extractor is a tool specification which define how to
// extract trace parent context from propagation context
type Extractor func() (string, error)

// Injector is a tool specification which define how to
// inject trace context into propagation context
type Injector func(header string) error

// TracerOption allows for functional options to adjust behaviour
// of a Tracer to be created by NewTracer

type SpanLayer int32
type Tag string

type Span interface {
	SetOperationName(string)
	GetOperationName() string
	SetPeer(string)
	SetSpanLayer(int32)
	SetComponent(int32)
	Tag(string, string)
	Log(time.Time, ...string)
	Error(time.Time, ...string)
	End()
	IsEntry() bool
	IsExit() bool
	GetRealSpan() interface{}
}
