package skywalking

import (
	"context"
	"fmt"
	"time"

	"strings"

	"github.com/google/uuid"
 
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	v3 "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"

	"github.com/micro-plat/hydra/context/apm"
	"github.com/micro-plat/hydra/pkgs"

)

func NewTracer(service string, opts ...apm.TracerOption) (tracer apm.Tracer, err error) {
	innertracer := &innertracer{}
	for _, o := range opts {
		o(innertracer)
	}
	skyopts := []go2sky.TracerOption{}

	if innertracer.reporter != nil {
		rpter := innertracer.reporter.GetRealReporter().(go2sky.Reporter)
		skyopts = append(skyopts, go2sky.WithReporter(rpter))
	}

	if innertracer.instance == "" {
		innertracer.instance = getTraceID()
	}
	skyopts = append(skyopts, go2sky.WithInstance(innertracer.instance))

	skytracer, err := go2sky.NewTracer(service, skyopts...)
	if err != nil {
		err = fmt.Errorf("构建skywalking tracer 失败：%+v", err)
		return
	}
	innertracer.tracer = skytracer
	tracer = innertracer
	return
}

type innertracer struct {
	tracer   *go2sky.Tracer
	reporter apm.Reporter
	instance string
}

// WithReporter setup report pipeline for tracer
func WithReporter(reporter apm.Reporter) apm.TracerOption {
	return func(t apm.Tracer) {
		//fmt.Println("t:", t, ",reporter:", reporter)
		t.SetReporter(reporter)
	}
}

// WithInstance setup instance identify
func WithInstance(instance string) apm.TracerOption {
	return func(t apm.Tracer) {
		t.SetInstance(instance)
	}
}
func (t *innertracer) SetReporter(reporter apm.Reporter) {
	t.reporter = reporter
}

func (t *innertracer) SetInstance(instance string) {
	t.instance = instance
}

func (t *innertracer) CreateLocalSpan(ctx context.Context, opts ...apm.SpanOption) (s apm.Span, c context.Context, err error) {
	ospan, c, err := t.tracer.CreateLocalSpan(ctx)
	inner := &innerSpan{
		span: ospan,
	}
	for _, o := range opts {
		o(inner)
	}
	s = inner
	return
}
func (t *innertracer) CreateEntrySpan(ctx context.Context, operationName string, extractor apm.Extractor, opts ...apm.SpanOption) (s apm.Span, nCtx context.Context, err error) {
	ospan, nCtx, err := t.tracer.CreateEntrySpan(ctx, operationName, propagation.Extractor(extractor))
	inner := &innerSpan{
		span: ospan,
	}
	for _, o := range opts {
		o(inner)
	}
	s = inner
	return
}
func (t *innertracer) CreateExitSpan(ctx context.Context, operationName string, peer string, injector apm.Injector, opts ...apm.SpanOption) (s apm.Span, err error) {
	ospan, err := t.tracer.CreateExitSpan(ctx, operationName, peer, propagation.Injector(injector))
	inner := &innerSpan{
		span: ospan,
	}
	for _, o := range opts {
		o(inner)
	}
	s = inner
	return
}
func (t *innertracer) GetRealTracer() interface{} {
	return t.tracer
}

type innerSpan struct {
	span go2sky.Span
}

// For Span
func (ds *innerSpan) SetOperationName(name string) {
	ds.span.SetOperationName(name)
}

func (ds *innerSpan) GetOperationName() string {
	return ds.span.GetOperationName()
}

func (ds *innerSpan) SetPeer(peer string) {
	ds.span.SetPeer(peer)
}

func (ds *innerSpan) SetSpanLayer(layer int32) {
	ds.span.SetSpanLayer(v3.SpanLayer(layer))
}

func (ds *innerSpan) SetComponent(componentID int32) {
	ds.span.SetComponent(componentID)
}

func (ds *innerSpan) Tag(key string, value string) {
	ds.span.Tag(go2sky.Tag(key), value)
}

func (ds *innerSpan) Log(time time.Time, ll ...string) {
	ds.span.Log(time, ll...)
}

func (ds *innerSpan) Error(time time.Time, ll ...string) {
	ds.span.Error(time, ll...)
}

func (ds *innerSpan) End() {
	ds.span.End()
}

func (ds *innerSpan) IsEntry() bool {
	return ds.span.IsEntry()
}

func (ds *innerSpan) IsExit() bool {
	return ds.span.IsExit()
}

func (ds *innerSpan) GetRealSpan() interface{} {
	return ds.span
}

// UUID generate UUID
func getTraceID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		panic("获取github.com/google/uuid失败")
	}
	return strings.ReplaceAll(id.String(), "-", "") + "@" + pkgs.LocalIP()
}
