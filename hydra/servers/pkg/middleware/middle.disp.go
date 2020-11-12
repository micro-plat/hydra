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
}

//
func (g *dispCtx) GetRouterPath() string {
	return g.Context.Request.GetService()
}
func (g *dispCtx) GetBody() io.ReadCloser {
	text := g.Request.GetForm()["__body_"]
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
func (g *dispCtx) GetMethod() string {
	return g.Context.Request.GetMethod()
}
func (g *dispCtx) GetURL() *url.URL {
	u, _ := url.ParseRequestURI(g.Context.Request.GetService())
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
func (g *dispCtx) GetQuery(k string) (string, bool) {
	v, ok := g.Context.Request.GetForm()[k]
	return fmt.Sprint(v), ok
}
func (g *dispCtx) GetFormValue(k string) (string, bool) {
	v, ok := g.Context.Request.GetForm()[k]
	return fmt.Sprint(v), ok
}

func (g *dispCtx) GetForm() url.Values {
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
		"__body_": body,
	})
}
func (g *dispCtx) ShouldBind(v interface{}) error {
	f := g.Context.Request.GetForm()
	if body, ok := f["__body_"]; ok && len(f) == 1 {
		switch msg := body.(type) {
		case json.RawMessage:
			return json.Unmarshal(msg, v)
		case []byte:
			return json.Unmarshal(msg, v)
		default:
			return json.Unmarshal([]byte(fmt.Sprint(msg)), v)
		}
	}
	js, err := json.Marshal(f)
	if err != nil {
		return fmt.Errorf("ShouldBind将输入的信息转换为JSON时失败 %w", err)
	}
	return json.Unmarshal(js, v)
}

func (g *dispCtx) GetFile(fileKey string) (string, io.ReadCloser, int64, error) {
	return "", nil, 0, nil
}
