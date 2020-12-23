package internal

// import (
// 	r "context"
// 	"sync"

// 	"github.com/SkyAPM/go2sky"
// 	"github.com/micro-plat/hydra/context"
// )

// //Span 事务处理跨度
// type Span struct {
// 	tracer    *go2sky.Tracer
// 	span      go2sky.Span
// 	ctx       r.Context
// 	operator  string
// 	subs      []*Span
// 	once      sync.Once
// 	avaliable bool
// }

// //New 创建一个处理
// func New(ctx r.Context, tracer *go2sky.Tracer, operator string) *Span {
// 	if tracer == nil {
// 		return &Span{operator: operator, ctx: ctx}
// 	}
// 	span := &Span{ctx: ctx, operator: operator, avaliable: true, tracer: tracer, subs: make([]*Span, 0, 1)}
// 	return span
// }

// //Start 启动任务
// func (s *Span) Start() context.IEnd {
// 	if !s.avaliable {
// 		return s
// 	}
// 	var err error
// 	s.span, s.ctx, err = s.tracer.CreateLocalSpan(s.ctx)
// 	if err != nil {
// 		return s
// 	}
// 	s.span.SetOperationName(s.operator)
// 	return s
// }

// //NewSpan 创建新的跟踪器
// func (s *Span) NewSpan(operator string) context.ITraceSpan {
// 	if !s.avaliable {
// 		return s
// 	}
// 	sub := New(s.ctx, s.tracer, operator)
// 	s.subs = append(s.subs, sub)
// 	return sub
// }

// //Available 是否可用
// func (s *Span) Available() bool {
// 	return s.avaliable
// }

// //End 处理完成
// func (s *Span) End() {
// 	s.once.Do(func() {
// 		if s.span != nil {
// 			// fmt.Println("end:", s.operator)
// 			s.span.End()
// 		}
// 		for _, v := range s.subs {
// 			v.End()
// 		}
// 	})
// }
