package services

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
)

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

type testHandlerSuffix struct{}

func (h testHandlerSuffix) XxxHandle(ctx string) (r interface{}) {
	return "sucess"
}

type testHandler1 struct{}

func newTestHandler() (*testHandler, error) {
	return &testHandler{}, nil
}

type testHandler struct{}

func (t *testHandler) GetHandling(context.IContext) interface{} {
	return nil
}
func (t *testHandler) GetHandle(context.IContext) interface{} {
	return nil
}

func (t *testHandler) PostHandle(context.IContext) interface{} {
	return nil
}
func (t *testHandler) PutHandle(context.IContext) interface{} {
	return nil
}
func (t *testHandler) DeleteHandle(context.IContext) interface{} {
	return nil
}
func (t *testHandler) Handled(context.IContext) interface{} {
	return nil
}

func (t *testHandler) OrderHandle(context.IContext) interface{} {
	return nil
}
func (t *testHandler) OrderFallback(context.IContext) interface{} {
	return nil
}
func (t *testHandler) Order(string) interface{} {
	return nil
}

func (t *testHandler) Close() error {
	return nil
}

type testHandler2 struct{}

func (t *testHandler2) PostHandling(context.IContext) interface{} {
	return nil
}
func (t *testHandler2) Handling(context.IContext) interface{} {
	return nil
}

func (t *testHandler2) PostHandle(context.IContext) interface{} {
	return nil
}
func (t *testHandler2) PostHandled(context.IContext) interface{} {
	return nil
}
func (t *testHandler2) Handled(context.IContext) interface{} {
	return nil
}

func (t *testHandler2) Handle(context.IContext) interface{} {
	return nil
}

func (t *testHandler2) Fallback(context.IContext) interface{} {
	return nil
}
func (t *testHandler2) OrderHandle(context.IContext) interface{} {
	return nil
}

func (t *testHandler2) OrderFallback(context.IContext) interface{} {
	return nil
}

func (t *testHandler2) Close() error {
	return nil
}

type testHandler3 struct{}

func (t *testHandler3) PostHandle(context.IContext) interface{} {
	return nil
}
func (t *testHandler3) GetHandle(context.IContext) interface{} {
	return nil
}
func (t *testHandler3) Handle(context.IContext) interface{} {
	return nil
}

type testHandler4 struct{}

func (t *testHandler4) Handle(context.IContext) interface{} {
	return nil
}
func (t *testHandler4) Close() error {
	return fmt.Errorf("error")
}

type testHandler5 struct{}

func (t *testHandler5) PostHandle(context.IContext) interface{} {
	return nil
}
func (t *testHandler5) OrderHandling(context.IContext) interface{} {
	return nil
}
func (t *testHandler5) OrderHandle(context.IContext) interface{} {
	return nil
}

type testHandler6 struct{}

func (t *testHandler6) OrderHandling(context.IContext) interface{} {
	return nil
}

type testHandler7 struct{}

func (t testHandler7) PostHandle(context.IContext) interface{} {
	return nil
}
func (t testHandler7) OrderHandle(context.IContext) interface{} {
	return nil
}
func (t testHandler7) Handle(context.IContext) interface{} {
	return nil
}

type testHandler8 struct{}

func (t testHandler8) QueryHandling(context.IContext) interface{} {
	return nil
}
func (t testHandler8) QueryHandle(context.IContext) interface{} {
	return nil
}

func newTestHandler9() (*testHandler9, error) {
	return &testHandler9{}, nil
}

type testHandler9 struct{}

func (t *testHandler9) FileHandling(context.IContext) interface{} {
	return nil
}

func (t *testHandler9) FileHandle(context.IContext) interface{} {
	return nil
}

func SamePahtTestHandle(ctx context.IContext) interface{} {
	return nil
}
func SamePahtTestHandling(ctx context.IContext) interface{} {
	return nil
}
