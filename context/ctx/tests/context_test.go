package tests

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	h "github.com/micro-plat/hydra/hydra/servers/http"
)

type TestContxt struct {
	clientIP   string
	body       io.ReadCloser
	method     string
	url        *url.URL
	header     http.Header
	cookie     []*http.Cookie
	param      string
	routerpath string
	form       url.Values
}

func (t *TestContxt) ClientIP() string {
	return t.clientIP
}
func (t *TestContxt) GetBody() io.ReadCloser {
	return t.body
}
func (t *TestContxt) GetMethod() string {
	return t.method
}
func (t *TestContxt) GetURL() *url.URL {
	return t.url
}
func (t *TestContxt) Header(string, string) {

}
func (t *TestContxt) GetHeaders() http.Header {
	return t.header
}
func (t *TestContxt) GetCookies() []*http.Cookie {
	return t.cookie
}
func (t *TestContxt) Param(string) string {
	return t.param
}
func (t *TestContxt) GetRouterPath() string {
	return t.routerpath
}
func (t *TestContxt) ShouldBind(interface{}) error {
	return nil
}
func (t *TestContxt) GetForm() url.Values {
	return t.form
}

func (t *TestContxt) GetQuery(string) (string, bool) {
	return "", false
}
func (t *TestContxt) GetFormValue(k string) (string, bool) {
	return "", false
}
func (t *TestContxt) ContentType() string {
	return ""
}
func (t *TestContxt) Abort() {}

func (t *TestContxt) WStatus(int) {}
func (t *TestContxt) Status() int {
	return 0
}
func (t *TestContxt) Written() bool {
	return false
}
func (t *TestContxt) WHeader(string) string {
	return ""
}
func (t *TestContxt) File(string)              {}
func (t *TestContxt) Data(int, string, []byte) {}
func (t *TestContxt) Redirect(int, string)     {}

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
