package context

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/hydra/servers/http"
	h "github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNewCtx(t *testing.T) {
	if _, err := app.Cache.GetAPPConf(h.API); err != nil {
		confObj := mocks.NewConf()         //构建对象
		confObj.API(":8080")               //初始化参数
		serverConf := confObj.GetAPIConf() //获取配置
		_, _ = http.NewResponsive(serverConf)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("recover:NewCtx() = %v", r)
		}
	}()
	got := ctx.NewCtx(&mocks.TestContxt{}, h.API)
	assert.NotEqual(t, nil, got, "获取ctx对象")
}

func TestCtx_Close(t *testing.T) {
	if _, err := app.Cache.GetAPPConf(h.API); err != nil {
		confObj := mocks.NewConf()         //构建对象
		confObj.API(":8080")               //初始化参数
		serverConf := confObj.GetAPIConf() //获取配置
		_, _ = http.NewResponsive(serverConf)
	}

	c := ctx.NewCtx(&mocks.TestContxt{}, h.API)

	c.Close()

	//对ctx.funcs和ctx.context为空不能进行判断
	if !reflect.ValueOf(c.Response()).IsNil() {
		t.Errorf("Close():c.response is not nil")
		return
	}
	if c.APPConf() != nil {
		t.Errorf("Close():c.APPConf is not nil")
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
