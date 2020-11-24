package context

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/encoding"
)

func Test_body_GetRawBody(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
		want        string
	}{
		{name: "1.1 content-type为application/xml,方法为POST", contentType: "application/xml", method: "POST", body: `body`, want: `body`},
		{name: "1.2 content-type为application/xml,方法为PUT", contentType: "application/xml", method: "PUT", body: `body`, want: `body`},
		{name: "1.3 content-type为application/xml,方法为DELETE", contentType: "application/xml", method: "DELETE", body: `body`, want: `body`},
		{name: "1.4 content-type为application/xml,方法为GET", contentType: "application/xml", method: "GET", body: `body`, want: `body`},
		{name: "1.5 content-type为application/xml,方法为PATCH", contentType: "application/xml", method: "PATCH", body: `body`, want: `body`},

		{name: "2.1 content-type为text/xml,方法为POST", contentType: "text/xml", method: "POST", body: `body`, want: `body`},
		{name: "2.2 content-type为text/xml,方法为PUT", contentType: "text/xml", method: "PUT", body: `body`, want: `body`},
		{name: "2.3 content-type为text/xml,方法为DELETE", contentType: "text/xml", method: "DELETE", body: `body`, want: `body`},
		{name: "2.4 content-type为text/xml,方法为GET", contentType: "text/xml", method: "GET", body: `body`, want: `body`},
		{name: "2.5 content-type为text/xml,方法为PATCH", contentType: "text/xml", method: "PATCH", body: `body`, want: `body`},

		{name: "3.1 content-type为application/json,方法为POST", contentType: "application/json", method: "POST", body: `body`, want: `body`},
		{name: "3.2 content-type为application/json,方法为PUT", contentType: "application/json", method: "PUT", body: `body`, want: `body`},
		{name: "3.3 content-type为application/json,方法为DELETE", contentType: "application/json", method: "DELETE", body: `body`, want: `body`},
		{name: "3.4 content-type为application/json,方法为GET", contentType: "application/json", method: "GET", body: `body`, want: `body`},
		{name: "3.5 content-type为application/json,方法为PATCH", contentType: "application/json", method: "PATCH", body: `body`, want: `body`},

		{name: "4.1 content-type为text/json,方法为POST", contentType: "text/json", method: "POST", body: `body`, want: `body`},
		{name: "4.2 content-type为text/json,方法为PUT", contentType: "text/json", method: "PUT", body: `body`, want: `body`},
		{name: "4.3 content-type为text/json,方法为DELETE", contentType: "json/json", method: "DELETE", body: `body`, want: `body`},
		{name: "4.4 content-type为text/json,方法为GET", contentType: "text/json", method: "GET", body: `body`, want: `body`},
		{name: "4.5 content-type为text/json,方法为PATCH", contentType: "text/json", method: "PATCH", body: `body`, want: `body`},

		{name: "5.1 content-type为application/x-yaml,方法为POST", contentType: "application/x-yaml", method: "POST", body: `body`, want: `body`},
		{name: "5.2 content-type为application/x-yaml,方法为PUT", contentType: "application/x-yaml", method: "PUT", body: `body`, want: `body`},
		{name: "5.3 content-type为application/x-yaml,方法为DELETE", contentType: "application/x-yaml", method: "DELETE", body: `body`, want: `body`},
		{name: "5.4 content-type为application/x-yaml,方法为GET", contentType: "application/x-yaml", method: "GET", body: `body`, want: `body`},
		{name: "5.5 content-type为application/x-yaml,方法为PATCH", contentType: "application/x-yaml", method: "PATCH", body: `body`, want: `body`},

		{name: "6.1 content-type为text/plain,方法为POST", contentType: "text/plain", method: "POST", body: `body`, want: `body`},
		{name: "6.2 content-type为text/plain,方法为PUT", contentType: "text/plain", method: "PUT", body: `body`, want: `body`},
		{name: "6.3 content-type为text/plain,方法为DELETE", contentType: "text/plain", method: "DELETE", body: `body`, want: `body`},
		{name: "6.4 content-type为text/plain,方法为GET", contentType: "text/plain", method: "GET", body: `body`, want: `body`},
		{name: "6.5 content-type为text/plain,方法为PATCH", contentType: "text/plain", method: "PATCH", body: `body`, want: `body`},

		{name: "7.1 content-type为application/x-www-form-urlencoded,方法为POST", contentType: "application/x-www-form-urlencoded", method: "POST", body: `body`, want: ``},
		{name: "7.2 content-type为application/x-www-form-urlencoded,方法为PUT", contentType: "application/x-www-form-urlencoded", method: "PUT", body: `body`, want: ``},
		{name: "7.3 content-type为application/x-www-form-urlencoded,方法为DELETE", contentType: "application/x-www-form-urlencoded", method: "DELETE", body: `body`, want: `body`},
		{name: "7.4 content-type为application/x-www-form-urlencoded,方法为GET", contentType: "application/x-www-form-urlencoded", method: "GET", body: `body`, want: `body`},
		{name: "7.5 content-type为application/x-www-form-urlencoded,方法为PATCH", contentType: "application/x-www-form-urlencoded", method: "PATCH", body: `body`, want: ``},

		{name: "8.1 content-type为multipart/form-data,方法为POST", contentType: "multipart/form-data", method: "POST", body: "body", want: "body"},
		{name: "8.2 content-type为multipart/form-data,方法为PUT", contentType: "multipart/form-data", method: "PUT", body: "body", want: "body"},
		{name: "8.3 content-type为multipart/form-data,方法为DELETE", contentType: "multipart/form-data", method: "DELETE", body: "body", want: "body"},
		{name: "8.4 content-type为multipart/form-data,方法为GET", contentType: "multipart/form-data", method: "GET", body: "body", want: "body"},
		{name: "8.5 content-type为multipart/form-data,方法为PATCH", contentType: "multipart/form-data", method: "PATCH", body: "body", want: "body"},

		{name: "9.1 content-type为text/html,方法为POST", contentType: "text/plain", method: "POST", body: `body`, want: `body`},
		{name: "9.2 content-type为text/html,方法为PUT", contentType: "text/plain", method: "PUT", body: `body`, want: `body`},
		{name: "9.3 content-type为text/html,方法为DELETE", contentType: "text/plain", method: "DELETE", body: `body`, want: `body`},
		{name: "9.4 content-type为text/html,方法为GET", contentType: "text/plain", method: "GET", body: `body`, want: `body`},
		{name: "9.4 content-type为text/html,方法为PATCH", contentType: "text/plain", method: "PATCH", body: `body`, want: `body`},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "/url", strings.NewReader(tt.body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", tt.contentType)

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), "utf-8")
		gotS, err := w.GetRawBody()
		str := string(gotS)
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, str, tt.name)
		gotS2, err := w.GetRawBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}

var value = `中文~!@#$%^&*()_+{}:"<>?`

func Test_body_GetBody_MIMEXML(t *testing.T) {

	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		value       string
		body        string
		want        string
	}{
		{name: "1.1 content-type为application/xml,编码为UTF-8,方法为POST", contentType: "application/xml", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.2 content-type为application/xml,编码为UTF-8,方法为GET", contentType: "application/xml", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.3 content-type为application/xml,编码为UTF-8,方法为PUT", contentType: "application/xml", method: "PUT", encoding: "UTF-8", value: value},
		{name: "1.4 content-type为application/xml,编码为UTF-8,方法为DELETE", contentType: "application/xml", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "1.5 content-type为application/xml,编码为UTF-8,方法为PATCH", contentType: "application/xml", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "2.1 content-type为text/xml,编码为UTF-8,方法为POST", contentType: "text/xml", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.2 content-type为text/xml,编码为UTF-8,方法为GET", contentType: "text/xml", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.3 content-type为text/xml,编码为UTF-8,方法为PUT", contentType: "text/xml", method: "PUT", encoding: "UTF-8", value: value},
		{name: "2.4 content-type为text/xml,编码为UTF-8,方法为DELETE", contentType: "text/xml", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "2.5 content-type为text/xml,编码为UTF-8,方法为PATCH", contentType: "text/xml", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "3.1 content-type为application/xml,编码为GBK,方法为POST", contentType: "application/xml", method: "POST", encoding: "GBK", value: value},
		{name: "3.2 content-type为application/xml,编码为GBK,方法为GET", contentType: "application/xml", method: "POST", encoding: "GBK", value: value},
		{name: "3.3 content-type为application/xml,编码为GBK,方法为PUT", contentType: "application/xml", method: "PUT", encoding: "GBK", value: value},
		{name: "3.4 content-type为application/xml,编码为GBK,方法为DELETE", contentType: "application/xml", method: "DELETE", encoding: "GBK", value: value},
		{name: "3.5 content-type为application/xml,编码为GBK,方法为PATCH", contentType: "application/xml", method: "PATCH", encoding: "GBK", value: value},

		{name: "4.1 content-type为text/xml,编码为GBK,方法为POST", contentType: "text/xml", method: "POST", encoding: "GBK", value: value},
		{name: "4.2 content-type为text/xml,编码为GBK,方法为GET", contentType: "text/xml", method: "POST", encoding: "GBK", value: value},
		{name: "4.3 content-type为text/xml,编码为GBK,方法为PUT", contentType: "text/xml", method: "PUT", encoding: "GBK", value: value},
		{name: "4.4 content-type为text/xml,编码为GBK,方法为DELETE", contentType: "text/xml", method: "DELETE", encoding: "GBK", value: value},
		{name: "4.5 content-type为text/xml,编码为GBK,方法为PATCH", contentType: "text/xml", method: "PATCH", encoding: "GBK", value: value},
	}
	type xmlParams struct {
		XMLName xml.Name `xml:"xml"`
		Key     string   `xml:"key"`
	}
	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		data := []byte(tt.value)
		if strings.ToLower(tt.encoding) == encoding.GBK {
			data, _ = encoding.Encode(tt.value, encoding.GBK)
		}
		bodyRaw, _ := xml.Marshal(&xmlParams{
			Key: url.QueryEscape(string(data)),
		})
		s := bodyRaw
		body := string(s)

		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url", strings.NewReader(string(body)))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, "<xml><key>中文~!@#$%^&*()_+{}:\"<>?</key></xml>", gotS, tt.name)
		gotS2, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}

func Test_body_GetBody_MIMEJSON(t *testing.T) {

	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		value       string
		body        string
		want        string
	}{
		{name: "1.1 content-type为application/json,编码为UTF-8,方法为POST", contentType: "application/json", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.2 content-type为application/json,编码为UTF-8,方法为GET", contentType: "application/json", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.3 content-type为application/json,编码为UTF-8,方法为PUT", contentType: "application/json", method: "PUT", encoding: "UTF-8", value: value},
		{name: "1.4 content-type为application/json,编码为UTF-8,方法为DELETE", contentType: "application/json", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "1.5 content-type为application/json,编码为UTF-8,方法为PATCH", contentType: "application/json", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "2.1 content-type为text/json,编码为UTF-8,方法为POST", contentType: "text/json", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.2 content-type为text/json,编码为UTF-8,方法为GET", contentType: "text/json", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.3 content-type为text/json,编码为UTF-8,方法为PUT", contentType: "text/json", method: "PUT", encoding: "UTF-8", value: value},
		{name: "2.4 content-type为text/json,编码为UTF-8,方法为DELETE", contentType: "text/json", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "2.5 content-type为text/json,编码为UTF-8,方法为PATCH", contentType: "text/json", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "3.1 content-type为application/json,编码为GBK,方法为POST", contentType: "application/json", method: "POST", encoding: "GBK", value: value},
		{name: "3.2 content-type为application/json,编码为GBK,方法为GET", contentType: "application/json", method: "POST", encoding: "GBK", value: value},
		{name: "3.3 content-type为application/json,编码为GBK,方法为PUT", contentType: "application/json", method: "PUT", encoding: "GBK", value: value},
		{name: "3.4 content-type为application/json,编码为GBK,方法为DELETE", contentType: "application/json", method: "DELETE", encoding: "GBK", value: value},
		{name: "3.5 content-type为application/json,编码为GBK,方法为PATCH", contentType: "application/json", method: "PATCH", encoding: "GBK", value: value},

		{name: "3.1 content-type为text/json,编码为GBK,方法为POST", contentType: "text/json", method: "POST", encoding: "GBK", value: value},
		{name: "3.2 content-type为text/json,编码为GBK,方法为GET", contentType: "text/json", method: "POST", encoding: "GBK", value: value},
		{name: "3.3 content-type为text/json,编码为GBK,方法为PUT", contentType: "text/json", method: "PUT", encoding: "GBK", value: value},
		{name: "3.4 content-type为text/json,编码为GBK,方法为DELETE", contentType: "text/json", method: "DELETE", encoding: "GBK", value: value},
		{name: "3.5 content-type为text/json,编码为GBK,方法为PATCH", contentType: "text/json", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		data := []byte(tt.value)
		if strings.ToLower(tt.encoding) == encoding.GBK {
			data, _ = encoding.Encode(tt.value, encoding.GBK)
		}
		bodyRaw, _ := json.Marshal(map[string]string{
			"key": url.QueryEscape(string(data)),
		})
		body := string(bodyRaw)
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url", strings.NewReader(body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		want := `{"key":"中文~!@#$%^&*()_+{}:"<>?"}`
		assert.Equal(t, want, gotS, tt.name)
		gotS2, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}

func Test_body_GetBody_MIMEYAML(t *testing.T) {

	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		value       string
		body        string
		want        string
	}{
		{name: "1.1 content-type为application/x-yaml,编码为UTF-8,方法为POST", contentType: "application/x-yaml", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.2 content-type为application/x-yaml,编码为UTF-8,方法为GET", contentType: "application/x-yaml", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.3 content-type为application/x-yaml,编码为UTF-8,方法为PUT", contentType: "application/x-yaml", method: "PUT", encoding: "UTF-8", value: value},
		{name: "1.4 content-type为application/x-yaml,编码为UTF-8,方法为DELETE", contentType: "application/x-yaml", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "1.5 content-type为application/x-yaml,编码为UTF-8,方法为PATCH", contentType: "application/x-yaml", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "3.1 content-type为application/x-yaml,编码为GBK,方法为POST", contentType: "application/x-yaml", method: "POST", encoding: "GBK", value: value},
		{name: "3.2 content-type为application/x-yaml,编码为GBK,方法为GET", contentType: "application/x-yaml", method: "POST", encoding: "GBK", value: value},
		{name: "3.3 content-type为application/x-yaml,编码为GBK,方法为PUT", contentType: "application/x-yaml", method: "PUT", encoding: "GBK", value: value},
		{name: "3.4 content-type为application/x-yaml,编码为GBK,方法为DELETE", contentType: "application/x-yaml", method: "DELETE", encoding: "GBK", value: value},
		{name: "3.5 content-type为application/x-yaml,编码为GBK,方法为PATCH", contentType: "application/x-yaml", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}
		bodyRaw := "key: " + tt.value
		body := bodyRaw
		if strings.ToLower(tt.encoding) == encoding.GBK {
			s, _ := encoding.Encode(body, encoding.GBK)
			body = string(s)
		}
		body = url.QueryEscape(body)
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url", strings.NewReader(body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, bodyRaw, gotS, tt.name)
		gotS2, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}
}

func Test_body_GetBody_MIMEPlain(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		value       string
		body        string
		want        string
	}{
		{name: "2.1 content-type为text/plain,编码为UTF-8,方法为POST", contentType: "text/plain", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.2 content-type为text/plain,编码为UTF-8,方法为GET", contentType: "text/plain", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.3 content-type为text/plain,编码为UTF-8,方法为PUT", contentType: "text/plain", method: "PUT", encoding: "UTF-8", value: value},
		{name: "2.4 content-type为text/plain,编码为UTF-8,方法为DELETE", contentType: "text/plain", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "2.5 content-type为text/plain,编码为UTF-8,方法为PATCH", contentType: "text/plain", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "3.1 content-type为text/plain,编码为GBK,方法为POST", contentType: "text/plain", method: "POST", encoding: "GBK", value: value},
		{name: "3.2 content-type为text/plain,编码为GBK,方法为GET", contentType: "text/plain", method: "POST", encoding: "GBK", value: value},
		{name: "3.3 content-type为text/plain,编码为GBK,方法为PUT", contentType: "text/plain", method: "PUT", encoding: "GBK", value: value},
		{name: "3.4 content-type为text/plain,编码为GBK,方法为DELETE", contentType: "text/plain", method: "DELETE", encoding: "GBK", value: value},
		{name: "3.5 content-type为text/plain,编码为GBK,方法为PATCH", contentType: "text/plain", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		data := []byte(tt.value)
		if strings.ToLower(tt.encoding) == encoding.GBK {
			data, _ = encoding.Encode(tt.value, encoding.GBK)
		}
		body := url.QueryEscape(string(data))
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url", strings.NewReader(body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		want := value
		assert.Equal(t, want, gotS, tt.name)
		gotS2, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}

func Test_body_GetBody_MIMEHTML(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		value       string
		body        string
		want        string
	}{
		{name: "2.1 content-type为text/html,编码为UTF-8,方法为POST", contentType: "text/html", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.2 content-type为text/html,编码为UTF-8,方法为GET", contentType: "text/html", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.3 content-type为text/html,编码为UTF-8,方法为PUT", contentType: "text/html", method: "PUT", encoding: "UTF-8", value: value},
		{name: "2.4 content-type为text/html,编码为UTF-8,方法为DELETE", contentType: "text/html", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "2.5 content-type为text/html,编码为UTF-8,方法为PATCH", contentType: "text/html", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "3.1 content-type为text/html,编码为GBK,方法为POST", contentType: "text/html", method: "POST", encoding: "GBK", value: value},
		{name: "3.2 content-type为text/html,编码为GBK,方法为GET", contentType: "text/html", method: "POST", encoding: "GBK", value: value},
		{name: "3.3 content-type为text/html,编码为GBK,方法为PUT", contentType: "text/html", method: "PUT", encoding: "GBK", value: value},
		{name: "3.4 content-type为text/html,编码为GBK,方法为DELETE", contentType: "text/html", method: "DELETE", encoding: "GBK", value: value},
		{name: "3.5 content-type为text/html,编码为GBK,方法为PATCH", contentType: "text/html", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		data := []byte(tt.value)
		if strings.ToLower(tt.encoding) == encoding.GBK {
			data, _ = encoding.Encode(tt.value, encoding.GBK)
		}
		body := url.QueryEscape(string(data))
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url", strings.NewReader(body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		want := value
		assert.Equal(t, want, gotS, tt.name)
		gotS2, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}

func Test_body_GetBody_MIMEPOSTForm(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		isQuery     bool
		isBody      bool
		contentType string
		value       string
		body        string
		want        string
	}{
		{name: "1.1 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为POST,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.2 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为POST,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.3 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为POST,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.1 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为GET,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "UTF-8", value: value},
		{name: "2.2 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为GET,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "UTF-8", value: value},
		{name: "2.3 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为GET,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "UTF-8", value: value},
		{name: "3.1 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为DELETE,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "3.2 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为DELETE,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "3.3 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为DELETE,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "4.1 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PUT,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},
		{name: "4.2 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PUT,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},
		{name: "4.3 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PUT,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},
		{name: "5.1 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PATCH,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "UTF-8", value: value},
		{name: "5.2 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PATCH,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "UTF-8", value: value},
		{name: "5.3 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PATCH,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "GBK", value: value},

		{name: "1.4 content-type为application/x-www-form-urlencoded,编码为GBK,方法为POST,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "GBK", value: value},
		{name: "1.5 content-type为application/x-www-form-urlencoded,编码为GBK,方法为POST,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "GBK", value: value},
		{name: "1.6 content-type为application/x-www-form-urlencoded,编码为GBK,方法为POST,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "GBK", value: value},
		{name: "2.4 content-type为application/x-www-form-urlencoded,编码为GBK,方法为GET,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "GBK", value: value},
		{name: "2.5 content-type为application/x-www-form-urlencoded,编码为GBK,方法为GET,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "GBK", value: value},
		{name: "2.6 content-type为application/x-www-form-urlencoded,编码为GBK,方法为GET,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "GBK", value: value},
		{name: "3.4 content-type为application/x-www-form-urlencoded,编码为GBK,方法为DELETE,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "GBK", value: value},
		{name: "3.5 content-type为application/x-www-form-urlencoded,编码为GBK,方法为DELETE,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "GBK", value: value},
		{name: "3.6 content-type为application/x-www-form-urlencoded,编码为GBK,方法为DELETE,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "GBK", value: value},
		{name: "4.4 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PUT,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},
		{name: "4.5 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PUT,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},
		{name: "4.6 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PUT,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},
		{name: "5.4 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PATCH,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "GBK", value: value},
		{name: "5.5 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PATCH,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "GBK", value: value},
		{name: "5.6 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PATCH,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		data := []byte(tt.value)
		if strings.ToLower(tt.encoding) == encoding.GBK {
			data, _ = encoding.Encode(tt.value, encoding.GBK)
		}
		value := url.QueryEscape(string(data))

		queryRaw := ""
		if tt.isQuery {
			queryRaw = "a=1&query=" + value
		}
		body := ""
		if tt.isBody {
			body = "a=2&key=" + value
		}

		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+queryRaw, strings.NewReader(body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetBody()
		if tt.isQuery && !tt.isBody {
			assert.Equal(t, `a=1&query=中文~!@#$%^&*()_+{}:"<>?`, gotS, tt.name)
		}
		if tt.isBody && (tt.method == "GET" || tt.method == "DELETE") {
			assert.Equal(t, `a=2&key=中文~!@#$%^&*()_+{}:"<>?`, gotS, tt.name)
		}
		if tt.isBody && !tt.isQuery {
			assert.Equal(t, `a=2&key=中文~!@#$%^&*()_+{}:"<>?`, gotS, tt.name)
		}
		if tt.isQuery && tt.isBody && (tt.method != "GET" && tt.method != "DELETE") { //body和query的数据都有
			assert.Equal(t, `a=2&a=1&key=中文~!@#$%^&*()_+{}:"<>?&query=中文~!@#$%^&*()_+{}:"<>?`, gotS, tt.name)
		}
		gotS2, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}

func getTestMIMEMultipartPOSTForm() string {
	file, _ := os.Open("upload.test.txt")
	defer file.Close()
	body := &bytes.Buffer{}
	// 文件写入 body
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("upload", filepath.Base("upload.test.txt"))
	io.Copy(part, file)
	writer.Close()
	return body.String()
}

func getUploadBody() string {
	return "Content-Disposition: form-data; name=\"upload\"; filename=\"upload.test.txt\"\r\nContent-Type: application/octet-stream\r\n\r\nADASDASDASFHNOJM~!@#$%^&*"
}

func Test_body_GetBody_MIMEMultipartPOSTForm(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		value       string
		body        string
		want        string
	}{
		{name: "1.1 content-type为multipart/form-data,编码为UTF-8,方法为POST", contentType: "multipart/form-data", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.1 content-type为multipart/form-data,编码为UTF-8,方法为GET", contentType: "multipart/form-data", method: "GET", encoding: "UTF-8", value: value},
		{name: "3.1 content-type为multipart/form-data,编码为UTF-8,方法为DELETE", contentType: "multipart/form-data", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "4.1 content-type为multipart/form-data,编码为UTF-8,方法为PUT", contentType: "multipart/form-data", method: "PUT", encoding: "UTF-8", value: value},
		{name: "5.1 content-type为multipart/form-data,编码为UTF-8,方法为PATCH", contentType: "multipart/form-data", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "1.2 content-type为multipart/form-data,编码为GBK,方法为POST", contentType: "multipart/form-data", method: "POST", encoding: "GBK", value: value},
		{name: "2.2 content-type为multipart/form-data,编码为GBK,方法为GET", contentType: "multipart/form-data", method: "GET", encoding: "GBK", value: value},
		{name: "3.2 content-type为multipart/form-data,编码为GBK,方法为DELETE", contentType: "multipart/form-data", method: "DELETE", encoding: "GBK", value: value},
		{name: "4.2 content-type为multipart/form-data,编码为GBK,方法为PUT", contentType: "multipart/form-data", method: "PUT", encoding: "UTF-8", value: value},
		{name: "5.2 content-type为multipart/form-data,编码为GBK,方法为PATCH", contentType: "multipart/form-data", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		data := []byte(getTestMIMEMultipartPOSTForm())
		if strings.ToLower(tt.encoding) == encoding.GBK {
			data, _ = encoding.Encode(getTestMIMEMultipartPOSTForm(), encoding.GBK)
		}

		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url", strings.NewReader(url.QueryEscape(string(data))))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetBody()

		assert.Equal(t, true, strings.Contains(gotS, getUploadBody()), tt.name)

		gotS2, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}
