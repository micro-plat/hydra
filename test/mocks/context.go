package mocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Read struct {
	*strings.Reader
}

func (r Read) Close() error {
	return nil
}

type ErrRead struct {
}

func (r ErrRead) Read(b []byte) (n int, err error) {
	return 0, fmt.Errorf("读取出错")
}

func (r ErrRead) Close() error {
	return nil
}

type TestContxt struct {
	ClientIp        string
	Body            string
	Method          string
	URL             *url.URL
	HttpHeader      http.Header
	Cookie          []*http.Cookie
	Params          string
	Routerpath      string
	Form            url.Values
	HttpContentType string
}

func (t *TestContxt) ClientIP() string {
	return t.ClientIp
}
func (t *TestContxt) GetBody() io.ReadCloser {
	if t.Body == "TEST_BODY_READ_ERR" {
		return ErrRead{}
	}
	return Read{Reader: strings.NewReader(t.Body)}
}

func (t *TestContxt) GetMethod() string {
	return t.Method
}
func (t *TestContxt) GetURL() *url.URL {
	return t.URL
}
func (t *TestContxt) Header(string, string) {

}
func (t *TestContxt) GetHeaders() http.Header {
	return t.HttpHeader
}
func (t *TestContxt) GetCookies() []*http.Cookie {
	return t.Cookie
}
func (t *TestContxt) Param(string) string {
	return t.Params
}
func (t *TestContxt) GetRouterPath() string {
	return t.Routerpath
}

func (t *TestContxt) ShouldBind(v interface{}) error {
	Body := t.GetForm()["__body_"][0]
	return json.Unmarshal([]byte(fmt.Sprint(Body)), v)
}

func (t *TestContxt) GetForm() url.Values {
	return t.Form
}

func (t *TestContxt) GetQuery(string) (string, bool) {
	return "", false
}
func (t *TestContxt) GetFormValue(k string) (string, bool) {
	if v, ok := t.Form[k]; ok {
		return v[0], ok
	}
	return "", false
}
func (t *TestContxt) ContentType() string {
	return t.HttpContentType
}
func (t *TestContxt) Abort() {}

func (t *TestContxt) WStatus(int) {}
func (t *TestContxt) Status() int {
	return 0
}
func (t *TestContxt) Written() bool {
	return false
}
func (t *TestContxt) WHeader(k string) string {
	if v, ok := t.HttpHeader[k]; ok {
		return v[0]
	}
	return ""
}
func (t *TestContxt) File(string)              {}
func (t *TestContxt) Data(int, string, []byte) {}
func (t *TestContxt) Redirect(int, string)     {}
