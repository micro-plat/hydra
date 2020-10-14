package tests

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	h "github.com/micro-plat/hydra/hydra/servers/http"
)

func TestNewCtx(t *testing.T) {
	startServer()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("recover:NewCtx() = %v", r)
		}
	}()
	if got := ctx.NewCtx(&TestContxt{}, h.API); got == nil {
		t.Errorf("NewCtx() got nil")
		return
	}
}

func TestCtx_Close(t *testing.T) {

	startServer()

	c := ctx.NewCtx(&TestContxt{}, h.API)

	c.Close()

	//对ctx.funcs和ctx.context为空不能进行判断
	if !reflect.ValueOf(c.Response()).IsNil() {
		t.Errorf("Close():c.response is not nil")
		return
	}
	if c.ServerConf() != nil {
		t.Errorf("Close():c.serverconf is not nil")
		return
	}
	if !reflect.ValueOf(c.User()).IsNil() {
		t.Errorf("Close():c.user is not nil")
		return
	}
	if c.Context() != nil {
		t.Errorf("Close():c.ctx is not nil")
		return
	}
	if !reflect.ValueOf(c.Request()).IsNil() {
		t.Errorf("Close():c.request is not nil")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			return
		}
		t.Errorf("context.Del(c.tid) doesn't run")
	}()

	context.Current()
}
