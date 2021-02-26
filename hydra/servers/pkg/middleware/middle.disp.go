package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/lib4go/types"
)

type buffer struct {
	*bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}

func NewDispCtx() *dispCtx {
	return &dispCtx{Context: &dispatcher.Context{}}
}

type dispCtx struct {
	*dispatcher.Context
	service       string
	needClearAuth bool
	servicePrefix string
}

func (g *dispCtx) GetParams() map[string]interface{} {
	params := make(map[string]interface{})
	for _, v := range g.Context.Params {
		params[v.Key] = v.Value
	}
	return params
}
func (g *dispCtx) GetBody() io.ReadCloser {
	text := g.Request.GetForm()["__body__"]
	switch v := text.(type) {
	case json.RawMessage:
		b := bytes.NewBuffer([]byte(v))
		return &buffer{Buffer: b}
	case []byte:
		b := bytes.NewBuffer(v)
		return &buffer{Buffer: b}
	default:
		b := bytes.NewBufferString(types.GetString(text))
		return &buffer{Buffer: b}
	}
}
func (g *dispCtx) GetService() string {
	return g.Context.Request.GetService()
}
func (g *dispCtx) Service(service string) {
	g.service = service
}
func (g *dispCtx) GetMethod() string {
	return g.Context.Request.GetMethod()
}
func (g *dispCtx) GetURL() *url.URL {
	u, err := url.ParseRequestURI(g.Context.Request.GetService())
	if err != nil {
		global.Def.Log().Error("service不是有效的路径，转换为URL失败", err)
		return &url.URL{}
	}
	return u
}
func (g *dispCtx) GetHeaders() http.Header {
	hd := http.Header{}
	for k, v := range g.Context.Request.GetHeader() {
		hd[k] = []string{v}
	}
	return hd
}
func (g *dispCtx) GetCookies() []*http.Cookie {
	return nil
}

func (g *dispCtx) GetRawForm() map[string]interface{} {
	return g.Context.Request.GetForm()
}
func (g *dispCtx) GetPostForm() url.Values {
	values := url.Values{}
	for k, v := range g.Context.Request.GetForm() {
		values.Set(k, fmt.Sprint(v))
	}
	return values
}

func (g *dispCtx) WStatus(s int) {
	g.Context.Writer.WriteHeader(s)
}
func (g *dispCtx) Status() int {
	return g.Context.Writer.Status()
}
func (g *dispCtx) Written() bool {
	return g.Context.Writer.Written()
}
func (g *dispCtx) WHeaders() http.Header {
	return g.Context.Writer.Header()
}
func (g *dispCtx) WHeader(k string) string {
	return g.Context.Writer.Header().Get(k)
}
func (g *dispCtx) ClientIP() string {
	if ip := g.GetHeader("Client-IP"); ip != "" {
		return ip
	}
	return g.Context.GetClientIP()
}
func (g *dispCtx) ContentType() string {
	return g.Context.GetHeader("Content-Type")
}
func (g *dispCtx) File(name string) {
	ff, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	body := base64.StdEncoding.EncodeToString(ff)
	g.Context.Header("file", name)
	g.Context.JSON(200, map[string]interface{}{
		"__body__": body,
	})
}

func (g *dispCtx) GetFile(fileKey string) (string, io.ReadCloser, int64, error) {
	return "", nil, 0, nil
}

//GetHTTPReqResp 获取http请求与响应对象
func (g *dispCtx) GetHTTPReqResp() (*http.Request, http.ResponseWriter) {
	return nil, nil
}
func (g *dispCtx) ClearAuth(c ...bool) bool {
	if len(c) == 0 {
		return g.needClearAuth
	}
	g.needClearAuth = types.GetBoolByIndex(c, 0, false)
	return g.needClearAuth
}

func (g *dispCtx) ServeContent(filepath string, fs http.FileSystem) int {
	return http.StatusOK
}

//
func (g *dispCtx) GetRouterPath() string {
	return g.Context.Request.GetService()
}

func (g *dispCtx) ServicePrefix(prefix string) {
	g.servicePrefix = prefix
}
