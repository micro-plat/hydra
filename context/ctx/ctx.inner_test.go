package ctx

import (
	"io"
	"net/http"
	"net/url"
)

type TestContxt struct {
	clientIP    string
	body        io.ReadCloser
	method      string
	url         *url.URL
	header      http.Header
	cookie      []*http.Cookie
	param       string
	routerpath  string
	form        url.Values
	contentType string
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
	return t.contentType
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
