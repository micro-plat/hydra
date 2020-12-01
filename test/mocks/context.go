package mocks

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/micro-plat/lib4go/types"
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
	params          map[string]interface{}
	Routerpath      string
	Form            url.Values
	HttpContentType string
	Doen            bool
	WrittenStatus   bool
	StatusCode      int
	Content         []byte
	FileStr         string
	Url             string
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
func (t *TestContxt) Header(k string, v string) {
	t.HttpHeader[k] = []string{v}
	return
}
func (t *TestContxt) GetHeaders() http.Header {
	return t.HttpHeader
}
func (t *TestContxt) GetCookies() []*http.Cookie {
	return t.Cookie
}
func (t *TestContxt) GetParams() map[string]interface{} {
	return t.params
}
func (t *TestContxt) GetRouterPath() string {
	return t.Routerpath
}

func (t *TestContxt) GetPostForm() url.Values {
	return t.Form
}

func (t *TestContxt) ContentType() string {
	if v, ok := t.HttpHeader["Content-Type"]; ok {
		return types.GetString(v)
	}
	return ""
}
func (t *TestContxt) Abort() {
	t.Doen = true
}

func (t *TestContxt) WStatus(s int) {
	t.StatusCode = s
	return
}
func (t *TestContxt) Status() int {
	return t.StatusCode
}
func (t *TestContxt) Written() bool {
	return t.WrittenStatus
}
func (t *TestContxt) WHeader(k string) string {
	if v, ok := t.HttpHeader[k]; ok {
		return fmt.Sprint(v)
	}
	return ""
}
func (t *TestContxt) File(s string) {
	t.WrittenStatus = true
	t.FileStr = s
}
func (t *TestContxt) Data(s int, ctp string, c []byte) {
	t.WrittenStatus = true
	t.StatusCode = s
	t.HttpHeader["Content-Type"] = []string{ctp}
	t.Content = c
	return
}
func (t *TestContxt) Redirect(s int, u string) {
	t.StatusCode = s
	t.Url = u
	t.WrittenStatus = true
	return
}

func (t *TestContxt) GetFile(fileKey string) (string, io.ReadCloser, int64, error) {
	return "", nil, 0, nil
}
