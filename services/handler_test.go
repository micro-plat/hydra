package services

import "github.com/micro-plat/hydra/context"

type hander1 struct {
}

func (h hander1) Handle(ctx context.IContext) (r interface{}) {
	return "sucess"
}

type hander2 struct {
}

func (h hander2) Handle(ctx context.IContext) (r interface{}) {
	return "sucess"
}

type testHandler1 struct{}

type testHandler struct{}

func (t *testHandler) GetHandling(context.IContext) interface{} {
	return nil
}
func (t *testHandler) PostHandle(context.IContext) interface{} {
	return nil
}
func (t *testHandler) Handled(context.IContext) interface{} {
	return nil
}
func (t *testHandler) OrderFallback(context.IContext) interface{} {
	return nil
}

func (t *testHandler) Close() error {
	return nil
}
