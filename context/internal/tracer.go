package internal

// //Tracer 跟踪器
// type Tracer struct {
// 	reporter go2sky.Reporter
// 	*Span
// }

// var tracers = cmap.New(8)

// //Empty 空跟踪器
// var Empty = &Tracer{Span: New(context.Background(), nil, "")}

// //GetTracer 创建跟踪器
// func GetTracer(service string, c app.IAPPConf) (*Tracer, error) {
// 	conf, err := c.GetAPMConf()
// 	if err != nil || conf.Disable {
// 		return nil, err
// 	}
// 	reporter, err := reporter.NewGRPCReporter(conf.Address)
// 	if reporter == nil || err != nil {
// 		return &Tracer{Span: New(context.Background(), nil, service)}, err
// 	}
// 	tracer, err := go2sky.NewTracer(service, go2sky.WithReporter(reporter))
// 	if err != nil {
// 		return &Tracer{Span: New(context.Background(), nil, service)}, err
// 	}
// 	return &Tracer{reporter: reporter, Span: New(context.Background(), tracer, service)}, nil
// }

// //Root 根节点
// func (t *Tracer) Root() *Span {
// 	return t.Span
// }

// //End 结束跟踪
// func (t *Tracer) End() {
// 	t.reporter.Close()
// 	t.Span.End()
// }
