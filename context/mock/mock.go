package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/micro-plat/hydra/context"
	xhttp "github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/lib4go/types"
)

//mock 用于context的mock包
type mock struct {
	RHeaders   types.XMap
	wHeaders   types.XMap
	Cookies    types.XMap
	serverType string
	result     []byte
	Service    string
	status     int
	Body       string
	URL        string
}

//newMock 构建
func newMock(content string) *mock {
	mk := &mock{
		serverType: xhttp.API,
		RHeaders:   make(types.XMap),
		wHeaders:   make(types.XMap),
		Cookies:    make(types.XMap),
	}
	mk.RHeaders["Content-Type"], mk.Body = getContentType(content)
	return mk
}

//ClientIP 获取客户端ＩＰ
func (m *mock) ClientIP() string {
	return m.RHeaders.GetString("client_ip")
}

//GetMethod 获取请求方法
func (m *mock) GetMethod() string {
	return m.RHeaders.GetString("method")
}

//GetURL 获取请求的URL
func (m *mock) GetURL() *url.URL {
	url, _ := url.Parse(m.URL)
	return url
}

//Header 设置头信息
func (m *mock) Header(k string, v string) {
	m.wHeaders[k] = v
}

//GetHeaders 获取头信息
func (m *mock) GetHeaders() http.Header {
	hd := make(map[string][]string)
	for k, v := range m.RHeaders {
		hd[k] = []string{fmt.Sprint(v)}
	}
	return hd
}

//GetCookies 获取Cookie
func (m *mock) GetCookies() []*http.Cookie {
	cks := make([]*http.Cookie, 0, 1)
	for k, v := range m.Cookies {
		cks = append(cks, &http.Cookie{Name: k, Value: fmt.Sprint(v)})
	}
	return cks
}

//GetParams 获取URL参数
func (m *mock) GetParams() map[string]interface{} {
	return nil
}

//GetRouterPath 获取请求路径
func (m *mock) GetRouterPath() string {
	return m.URL
}

//GetPostForm 获取POST参数
func (m *mock) GetPostForm() url.Values {
	q, _ := url.ParseQuery(m.URL)
	return q
}

//GetRawForm 获取原始请求参数
func (m *mock) GetRawForm() map[string]interface{} {
	return nil
}

//ContentType 获取Content-Type
func (m *mock) ContentType() string {
	return m.RHeaders.GetString("Content-Type")
}

//GetBody 获取Body数据
func (m *mock) GetBody() io.ReadCloser {
	b := bytes.NewBufferString(m.Body)
	return &buffer{Buffer: b}
}

//Abort 中止当前请求
func (m *mock) Abort() {
	return
}

//WStatus 设置当前状态码
func (m *mock) WStatus(int) {
	return
}

//Status 获取当前状态码
func (m *mock) Status() int {
	return m.status
}

//Written 当前流是否已写入数据
func (m *mock) Written() bool {
	return false
}

//WHeaders　获取响应头信息
func (m *mock) WHeaders() http.Header {
	hd := make(map[string][]string)
	for k, v := range m.wHeaders {
		hd[k] = []string{fmt.Sprint(v)}
	}
	return hd
}

//WHeader 获取响应头的值
func (m *mock) WHeader(string) string {
	return ""
}

//File 写入文件
func (m *mock) File(string) {
	return
}

//Data　设置响应数据
func (m *mock) Data(s int, t string, c []byte) {
	m.status = s
	m.wHeaders["Content-Type"] = t
	m.result = c
	return
}

//Redirect 转跳路径
func (m *mock) Redirect(s int, p string) {
	m.status = s
	m.wHeaders["Location"] = p
	return
}

//GetService 获取服务名
func (m *mock) GetService() string {
	return m.Service
}

//GetFile　获取上传的文件信息
func (m *mock) GetFile(fileKey string) (string, io.ReadCloser, int64, error) {
	return "", nil, 0, nil
}

//GetHTTPReqResp 获取Http请求与响应
func (m *mock) GetHTTPReqResp() (*http.Request, http.ResponseWriter) {
	return nil, nil
}

type buffer struct {
	*bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}
func getContentType(text string) (string, string) {
	switch {
	case strings.HasPrefix(text, "<!DOCTYPE html"):
		return context.HTMLF, text
	case strings.HasPrefix(text, "<") && strings.HasSuffix(text, ">"):
		_, errx := mxj.BeautifyXml([]byte(text), "", "")
		if errx != nil {
			return context.PLAINF, text
		}
		return context.XMLF, text
	case json.Valid([]byte(text)) && (strings.HasPrefix(text, "{") ||
		strings.HasPrefix(text, "[")):
		return context.JSONF, text
	default:
		_, err := url.ParseQuery(text)
		if err == nil {
			return "application/x-www-form-urlencoded", text
		}
		return context.PLAINF, text
	}
}
