package context

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/encoding"
)

func Test_body_GetRawBody(t *testing.T) {

	content, contentType := getTestMIMEMultipartPOSTForm(url.Values{"body": []string{"123456"}})
	data, _ := encoding.Encode(content, "utf-8")
	body := string(data)

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

		{name: "7.1 content-type为application/x-www-form-urlencoded,方法为POST", contentType: "application/x-www-form-urlencoded", method: "POST", body: `body`, want: `body`},
		{name: "7.2 content-type为application/x-www-form-urlencoded,方法为PUT", contentType: "application/x-www-form-urlencoded", method: "PUT", body: `body`, want: `body`},
		{name: "7.3 content-type为application/x-www-form-urlencoded,方法为DELETE", contentType: "application/x-www-form-urlencoded", method: "DELETE", body: `body`, want: `body`},
		{name: "7.4 content-type为application/x-www-form-urlencoded,方法为GET", contentType: "application/x-www-form-urlencoded", method: "GET", body: `body`, want: `body`},
		{name: "7.5 content-type为application/x-www-form-urlencoded,方法为PATCH", contentType: "application/x-www-form-urlencoded", method: "PATCH", body: `body`, want: `body`},

		{name: "8.1 content-type为multipart/form-data,方法为POST", contentType: contentType, method: "POST", body: body, want: "body=123456"},
		{name: "8.2 content-type为multipart/form-data,方法为PUT", contentType: contentType, method: "PUT", body: body, want: "body=123456"},
		{name: "8.3 content-type为multipart/form-data,方法为DELETE", contentType: contentType, method: "DELETE", body: body, want: "body=123456"},
		{name: "8.4 content-type为multipart/form-data,方法为GET", contentType: contentType, method: "GET", body: body, want: "body=123456"},
		{name: "8.5 content-type为multipart/form-data,方法为PATCH", contentType: contentType, method: "PATCH", body: body, want: "body=123456"},

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
		gotS, err := w.GetBody()
		str := string(gotS)
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, str, tt.name)
		gotS2, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}

var value = `中文~!@#$%^&*()_+{}:"<>?=`

func Test_body_GetFullRaw_MIMEXML(t *testing.T) {

	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		queryRaw    string
		value       string
		body        string
		want        string
	}{
		{name: "1.1 content-type为application/xml,编码为UTF-8,方法为POST", contentType: "application/xml", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.2 content-type为application/xml,编码为UTF-8,方法为GET", contentType: "application/xml", method: "GET", encoding: "UTF-8", value: value},
		{name: "1.3 content-type为application/xml,编码为UTF-8,方法为PUT", contentType: "application/xml", method: "PUT", encoding: "UTF-8", value: value},
		{name: "1.4 content-type为application/xml,编码为UTF-8,方法为DELETE", contentType: "application/xml", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "1.5 content-type为application/xml,编码为UTF-8,方法为PATCH", contentType: "application/xml", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "2.1 content-type为text/xml,编码为UTF-8,方法为POST", contentType: "text/xml", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.2 content-type为text/xml,编码为UTF-8,方法为GET", contentType: "text/xml", method: "GET", encoding: "UTF-8", value: value},
		{name: "2.3 content-type为text/xml,编码为UTF-8,方法为PUT", contentType: "text/xml", method: "PUT", encoding: "UTF-8", value: value},
		{name: "2.4 content-type为text/xml,编码为UTF-8,方法为DELETE", contentType: "text/xml", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "2.5 content-type为text/xml,编码为UTF-8,方法为PATCH", contentType: "text/xml", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "3.1 content-type为application/xml,编码为GBK,方法为POST", contentType: "application/xml", method: "POST", encoding: "GBK", value: value},
		{name: "3.2 content-type为application/xml,编码为GBK,方法为GET", contentType: "application/xml", method: "GET", encoding: "GBK", value: value},
		{name: "3.3 content-type为application/xml,编码为GBK,方法为PUT", contentType: "application/xml", method: "PUT", encoding: "GBK", value: value},
		{name: "3.4 content-type为application/xml,编码为GBK,方法为DELETE", contentType: "application/xml", method: "DELETE", encoding: "GBK", value: value},
		{name: "3.5 content-type为application/xml,编码为GBK,方法为PATCH", contentType: "application/xml", method: "PATCH", encoding: "GBK", value: value},

		{name: "4.1 content-type为text/xml,编码为GBK,方法为POST", contentType: "text/xml", method: "POST", encoding: "GBK", value: value},
		{name: "4.2 content-type为text/xml,编码为GBK,方法为GET", contentType: "text/xml", method: "GET", encoding: "GBK", value: value},
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

		bodyRaw, _ := xml.Marshal(&xmlParams{
			Key: tt.value,
		})

		data, _ := encoding.EncodeBytes(bodyRaw, tt.encoding)
		queryRaw := getTestQueryRaw(value, tt.encoding)
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+queryRaw, bytes.NewReader(data))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, query, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, string(data), string(gotS), tt.name)
		assert.Equal(t, queryRaw, query, tt.name)
		gotS2, query2, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
		assert.Equal(t, query, query2, tt.name+"再次读取body")
	}

}

func Test_body_GetFullRaw_MIMEJSON(t *testing.T) {

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
		{name: "1.2 content-type为application/json,编码为UTF-8,方法为GET", contentType: "application/json", method: "GET", encoding: "UTF-8", value: value},
		{name: "1.3 content-type为application/json,编码为UTF-8,方法为PUT", contentType: "application/json", method: "PUT", encoding: "UTF-8", value: value},
		{name: "1.4 content-type为application/json,编码为UTF-8,方法为DELETE", contentType: "application/json", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "1.5 content-type为application/json,编码为UTF-8,方法为PATCH", contentType: "application/json", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "1.6 content-type为application/json,编码为GBK,方法为POST", contentType: "application/json", method: "POST", encoding: "GBK", value: value},
		{name: "1.7 content-type为application/json,编码为GBK,方法为GET", contentType: "application/json", method: "GET", encoding: "GBK", value: value},
		{name: "1.8 content-type为application/json,编码为GBK,方法为PUT", contentType: "application/json", method: "PUT", encoding: "GBK", value: value},
		{name: "1.9 content-type为application/json,编码为GBK,方法为DELETE", contentType: "application/json", method: "DELETE", encoding: "GBK", value: value},
		{name: "1.10 content-type为application/json,编码为GBK,方法为PATCH", contentType: "application/json", method: "PATCH", encoding: "GBK", value: value},

		{name: "2.1 content-type为text/json,编码为UTF-8,方法为POST", contentType: "text/json", method: "POST", encoding: "UTF-8", value: value},
		{name: "2.2 content-type为text/json,编码为UTF-8,方法为GET", contentType: "text/json", method: "GET", encoding: "UTF-8", value: value},
		{name: "2.3 content-type为text/json,编码为UTF-8,方法为PUT", contentType: "text/json", method: "PUT", encoding: "UTF-8", value: value},
		{name: "2.4 content-type为text/json,编码为UTF-8,方法为DELETE", contentType: "text/json", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "2.5 content-type为text/json,编码为UTF-8,方法为PATCH", contentType: "text/json", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "2.6 content-type为text/json,编码为GBK,方法为POST", contentType: "text/json", method: "POST", encoding: "GBK", value: value},
		{name: "2.7 content-type为text/json,编码为GBK,方法为GET", contentType: "text/json", method: "GET", encoding: "GBK", value: value},
		{name: "2.8 content-type为text/json,编码为GBK,方法为PUT", contentType: "text/json", method: "PUT", encoding: "GBK", value: value},
		{name: "2.9 content-type为text/json,编码为GBK,方法为DELETE", contentType: "text/json", method: "DELETE", encoding: "GBK", value: value},
		{name: "2.10 content-type为text/json,编码为GBK,方法为PATCH", contentType: "text/json", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		data := []byte(tt.value)
		if strings.ToLower(tt.encoding) == encoding.GBK {
			data, _ = encoding.Encode(tt.value, encoding.GBK)
		}
		bodyRaw, _ := json.Marshal(map[string]string{
			"key": string(data),
		})
		queryRaw := getTestQueryRaw(value, tt.encoding)
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+queryRaw, bytes.NewReader(bodyRaw))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, query, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, string(bodyRaw), string(gotS), tt.name)
		assert.Equal(t, queryRaw, query, tt.name)
		gotS2, query2, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
		assert.Equal(t, query, query2, tt.name+"再次读取body")
	}

}

func Test_body_GetFullRaw_MIMEYAML(t *testing.T) {

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
		{name: "1.2 content-type为application/x-yaml,编码为UTF-8,方法为GET", contentType: "application/x-yaml", method: "GET", encoding: "UTF-8", value: value},
		{name: "1.3 content-type为application/x-yaml,编码为UTF-8,方法为PUT", contentType: "application/x-yaml", method: "PUT", encoding: "UTF-8", value: value},
		{name: "1.4 content-type为application/x-yaml,编码为UTF-8,方法为DELETE", contentType: "application/x-yaml", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "1.5 content-type为application/x-yaml,编码为UTF-8,方法为PATCH", contentType: "application/x-yaml", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "2.1 content-type为application/x-yaml,编码为GBK,方法为POST", contentType: "application/x-yaml", method: "POST", encoding: "GBK", value: value},
		{name: "2.2 content-type为application/x-yaml,编码为GBK,方法为GET", contentType: "application/x-yaml", method: "GET", encoding: "GBK", value: value},
		{name: "2.3 content-type为application/x-yaml,编码为GBK,方法为PUT", contentType: "application/x-yaml", method: "PUT", encoding: "GBK", value: value},
		{name: "2.4 content-type为application/x-yaml,编码为GBK,方法为DELETE", contentType: "application/x-yaml", method: "DELETE", encoding: "GBK", value: value},
		{name: "2.5 content-type为application/x-yaml,编码为GBK,方法为PATCH", contentType: "application/x-yaml", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}
		bodyRaw := []byte("key: " + tt.value)
		bodyRaw, _ = encoding.EncodeBytes(bodyRaw, tt.encoding)

		queryRaw := getTestQueryRaw(value, tt.encoding)
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+queryRaw, bytes.NewReader(bodyRaw))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, query, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, string(bodyRaw), string(gotS), tt.name)
		assert.Equal(t, queryRaw, query, tt.name)
		gotS2, query2, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
		assert.Equal(t, query, query2, tt.name+"再次读取body")
	}
}

func Test_body_GetFullRaw_MIMEPlain(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		value       string
		body        string
		want        string
	}{
		{name: "1.1 content-type为text/plain,编码为UTF-8,方法为POST", contentType: "text/plain", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.2 content-type为text/plain,编码为UTF-8,方法为GET", contentType: "text/plain", method: "GET", encoding: "UTF-8", value: value},
		{name: "1.3 content-type为text/plain,编码为UTF-8,方法为PUT", contentType: "text/plain", method: "PUT", encoding: "UTF-8", value: value},
		{name: "1.4 content-type为text/plain,编码为UTF-8,方法为DELETE", contentType: "text/plain", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "1.5 content-type为text/plain,编码为UTF-8,方法为PATCH", contentType: "text/plain", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "2.1 content-type为text/plain,编码为GBK,方法为POST", contentType: "text/plain", method: "POST", encoding: "GBK", value: value},
		{name: "2.2 content-type为text/plain,编码为GBK,方法为GET", contentType: "text/plain", method: "GET", encoding: "GBK", value: value},
		{name: "2.3 content-type为text/plain,编码为GBK,方法为PUT", contentType: "text/plain", method: "PUT", encoding: "GBK", value: value},
		{name: "2.4 content-type为text/plain,编码为GBK,方法为DELETE", contentType: "text/plain", method: "DELETE", encoding: "GBK", value: value},
		{name: "2.5 content-type为text/plain,编码为GBK,方法为PATCH", contentType: "text/plain", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		data := []byte(tt.value)
		data, _ = encoding.EncodeBytes(data, tt.encoding)
		queryRaw := getTestQueryRaw(value, tt.encoding)
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+queryRaw, bytes.NewReader(data))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, query, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, string(data), string(gotS), tt.name)
		assert.Equal(t, queryRaw, query, tt.name)
		gotS2, query2, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
		assert.Equal(t, query, query2, tt.name+"再次读取body")
	}

}

func Test_body_GetFullRaw_MIMEHTML(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		value       string
		body        string
		want        string
	}{
		{name: "1.1 content-type为text/html,编码为UTF-8,方法为POST", contentType: "text/html", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.2 content-type为text/html,编码为UTF-8,方法为GET", contentType: "text/html", method: "GET", encoding: "UTF-8", value: value},
		{name: "1.3 content-type为text/html,编码为UTF-8,方法为PUT", contentType: "text/html", method: "PUT", encoding: "UTF-8", value: value},
		{name: "1.4 content-type为text/html,编码为UTF-8,方法为DELETE", contentType: "text/html", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "1.5 content-type为text/html,编码为UTF-8,方法为PATCH", contentType: "text/html", method: "PATCH", encoding: "UTF-8", value: value},

		{name: "2.1 content-type为text/html,编码为GBK,方法为POST", contentType: "text/html", method: "POST", encoding: "GBK", value: value},
		{name: "2.2 content-type为text/html,编码为GBK,方法为GET", contentType: "text/html", method: "GET", encoding: "GBK", value: value},
		{name: "2.3 content-type为text/html,编码为GBK,方法为PUT", contentType: "text/html", method: "PUT", encoding: "GBK", value: value},
		{name: "2.4 content-type为text/html,编码为GBK,方法为DELETE", contentType: "text/html", method: "DELETE", encoding: "GBK", value: value},
		{name: "2.5 content-type为text/html,编码为GBK,方法为PATCH", contentType: "text/html", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		data := []byte(tt.value)
		data, _ = encoding.EncodeBytes(data, tt.encoding)
		queryRaw := getTestQueryRaw(value, tt.encoding)
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+queryRaw, bytes.NewReader(data))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, query, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, string(data), string(gotS), tt.name)
		assert.Equal(t, queryRaw, query, tt.name)
		gotS2, query2, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
		assert.Equal(t, query, query2, tt.name+"再次读取body")
	}

}

func Test_body_GetFullRaw_MIMEPOSTForm(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		isQuery     bool
		isBody      bool
		contentType string
		value       string
		want        string
	}{
		{name: "1.1 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为POST,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.2 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为POST,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "UTF-8", value: value},
		{name: "1.3 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为POST,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "UTF-8", value: value},

		{name: "1.4 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为GET,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "UTF-8", value: value},
		{name: "1.5 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为GET,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "UTF-8", value: value},
		{name: "1.6 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为GET,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "UTF-8", value: value},

		{name: "1.7 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为DELETE,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "1.8 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为DELETE,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "UTF-8", value: value},
		{name: "1.9 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为DELETE,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "UTF-8", value: value},

		{name: "1.10 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PUT,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},
		{name: "1.11 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PUT,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},
		{name: "1.12 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PUT,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},

		{name: "1.13 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PATCH,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "UTF-8", value: value},
		{name: "1.14 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PATCH,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "UTF-8", value: value},
		{name: "1.15 content-type为application/x-www-form-urlencoded,编码为UTF-8,方法为PATCH,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "GBK", value: value},

		{name: "2.1 content-type为application/x-www-form-urlencoded,编码为GBK,方法为POST,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "GBK", value: value},
		{name: "2.2 content-type为application/x-www-form-urlencoded,编码为GBK,方法为POST,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "GBK", value: value},
		{name: "2.3 content-type为application/x-www-form-urlencoded,编码为GBK,方法为POST,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "POST", encoding: "GBK", value: value},

		{name: "2.4 content-type为application/x-www-form-urlencoded,编码为GBK,方法为GET,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "GBK", value: value},
		{name: "2.5 content-type为application/x-www-form-urlencoded,编码为GBK,方法为GET,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "GBK", value: value},
		{name: "2.6 content-type为application/x-www-form-urlencoded,编码为GBK,方法为GET,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "GET", encoding: "GBK", value: value},

		{name: "2.7 content-type为application/x-www-form-urlencoded,编码为GBK,方法为DELETE,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "GBK", value: value},
		{name: "2.8 content-type为application/x-www-form-urlencoded,编码为GBK,方法为DELETE,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "GBK", value: value},
		{name: "2.9 content-type为application/x-www-form-urlencoded,编码为GBK,方法为DELETE,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "DELETE", encoding: "GBK", value: value},

		{name: "2.10 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PUT,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},
		{name: "2.11 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PUT,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},
		{name: "2.12 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PUT,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "PUT", encoding: "UTF-8", value: value},

		{name: "2.13 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PATCH,body为空,query不为空", isQuery: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "GBK", value: value},
		{name: "2.14 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PATCH,body不为空,query为空", isBody: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "GBK", value: value},
		{name: "2.15 content-type为application/x-www-form-urlencoded,编码为GBK,方法为PATCH,body和query不为空", isQuery: true, isBody: true, contentType: "application/x-www-form-urlencoded", method: "PATCH", encoding: "GBK", value: value},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		data, _ := encoding.EncodeBytes([]byte(tt.value), tt.encoding)

		queryRaw := ""
		if tt.isQuery {
			queryRaw = fmt.Sprintf("query=%s", url.QueryEscape(string(data)))
		}

		body := ""
		if tt.isBody {
			body = fmt.Sprintf("key=%s", url.QueryEscape(string(data)))
		}

		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+queryRaw, bytes.NewReader([]byte(body)))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, query, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, body, string(gotS), tt.name)
		assert.Equal(t, queryRaw, query, tt.name)
		gotS2, query2, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
		assert.Equal(t, query, query2, tt.name+"再次读取body")

		// gotS2, err := w.GetBody()
		// assert.Equal(t, nil, err, tt.name)
		// assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}

func Test_body_GetFullRaw_MIMEMultipartPOSTForm(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		encoding string
		value    url.Values
		body     string
		want     string
	}{
		{name: "1.1 content-type为multipart/form-data,编码为UTF-8,方法为POST", method: "POST", value: url.Values{"key": []string{value}}, encoding: "UTF-8"},
		{name: "1.2 content-type为multipart/form-data,编码为UTF-8,方法为GET", method: "GET", encoding: "UTF-8", value: url.Values{"key": []string{value}}},
		{name: "1.3 content-type为multipart/form-data,编码为UTF-8,方法为DELETE", method: "DELETE", encoding: "UTF-8", value: url.Values{"key": []string{value}}},
		{name: "1.4 content-type为multipart/form-data,编码为UTF-8,方法为PUT", method: "PUT", encoding: "UTF-8", value: url.Values{"key": []string{value}}},
		{name: "1.5 content-type为multipart/form-data,编码为UTF-8,方法为PATCH", method: "PATCH", encoding: "UTF-8", value: url.Values{"key": []string{value}}},

		{name: "2.1 content-type为multipart/form-data,编码为GBK,方法为POST", method: "POST", value: url.Values{"key": []string{value}}, encoding: "GBK"},
		{name: "2.2 content-type为multipart/form-data,编码为GBK,方法为GET", method: "GET", encoding: "GBK", value: url.Values{"key": []string{value}}},
		{name: "2.3 content-type为multipart/form-data,编码为GBK,方法为DELETE", method: "DELETE", encoding: "GBK", value: url.Values{"key": []string{value}}},
		{name: "2.4 content-type为multipart/form-data,编码为GBK,方法为PUT", method: "PUT", encoding: "GBK", value: url.Values{"key": []string{value}}},
		{name: "2.5 content-type为multipart/form-data,编码为GBK,方法为PATCH", method: "PATCH", encoding: "GBK", value: url.Values{"key": []string{value}}},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		//写入文件和设置参数
		content, contentType := getTestMIMEMultipartPOSTForm(tt.value)
		data, _ := encoding.Encode(content, tt.encoding)
		body := string(data)
		queryRaw := getTestQueryRaw(value, tt.encoding)
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+queryRaw, bytes.NewReader([]byte(body)))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, query, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		s, _ := encoding.Encode(value, tt.encoding)
		assert.Equal(t, "key="+url.QueryEscape(string(s)), string(gotS), tt.name)
		assert.Equal(t, queryRaw, query, tt.name)
		gotS2, query2, err := w.GetFullRaw()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
		assert.Equal(t, query, query2, tt.name+"再次读取body")
	}

}

type xmlParams struct {
	XMLName xml.Name `xml:"xml"`
	Key     string   `xml:"key"`
}

func getTestBody(value, e, ctp string) []byte {
	var bodyRaw = []byte("")
	switch ctp {
	case "xml":
		bodyRaw, _ = xml.Marshal(&xmlParams{
			Key: value,
		})
	case "json":
		bodyRaw, _ = json.Marshal(map[string]string{
			"key": value,
		})
	case "yaml":
		bodyRaw = []byte("key: " + value)
	case "form":
		buff, _ := encoding.Encode(value, e)
		v := url.QueryEscape(string(buff))
		bodyRaw = []byte("key=" + v)
		return bodyRaw
	}
	bodyRaw, _ = encoding.EncodeBytes(bodyRaw, e)
	return bodyRaw
}

func getTestQueryRaw(value, e string) string {
	t, _ := encoding.Encode(value, e)
	value = "key=" + url.QueryEscape(string(t))
	return value
}

func Test_body_GetMap__MIMEXML(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		queryRaw    string
		body        []byte
		errStr      string
		want        map[string]interface{}
	}{
		{name: "1.1 content-type为xml,编码为UTF-8,POST,body非正确的xml", method: "POST", encoding: "UTF-8", contentType: "application/xml", body: []byte("<xml>"), errStr: "将<xml>转换为map失败:xml.Decoder.Token() - XML syntax error on line 1: unexpected EOF"},
		{name: "1.2 content-type为xml,编码为UTF-8,POST,body为xml,url参数为空", method: "POST", encoding: "UTF-8", contentType: "application/xml", body: getTestBody(value, "UTF-8", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}}},
		{name: "1.3 content-type为xml,编码为UTF-8,POST,body为空,url带有参数", method: "POST", encoding: "UTF-8", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.4 content-type为xml,编码为UTF-8,POST,body为xml,url带有参数", method: "POST", encoding: "UTF-8", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}, "key": value}},

		{name: "1.5 content-type为xml,编码为UTF-8,GET,body非正确的xml", method: "GET", encoding: "UTF-8", contentType: "application/xml", body: []byte("<xml>"), errStr: "将<xml>转换为map失败:xml.Decoder.Token() - XML syntax error on line 1: unexpected EOF"},
		{name: "1.6 content-type为xml,编码为UTF-8,GET,body为xml,url参数为空", method: "GET", encoding: "UTF-8", contentType: "application/xml", body: getTestBody(value, "UTF-8", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}}},
		{name: "1.7 content-type为xml,编码为UTF-8,GET,body为空,url带有参数", method: "GET", encoding: "UTF-8", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.8 content-type为xml,编码为UTF-8,GET,body为xml,url带有参数", method: "GET", encoding: "UTF-8", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}, "key": value}},

		{name: "1.9 content-type为xml,编码为UTF-8,DELETE,body非正确的xml", method: "DELETE", encoding: "UTF-8", contentType: "application/xml", body: []byte("<xml>"), errStr: "将<xml>转换为map失败:xml.Decoder.Token() - XML syntax error on line 1: unexpected EOF"},
		{name: "1.10 content-type为xml,编码为UTF-8,DELETE,body为xml,url参数为空", method: "DELETE", encoding: "UTF-8", contentType: "application/xml", body: getTestBody(value, "UTF-8", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}}},
		{name: "1.11 content-type为xml,编码为UTF-8,DELETE,body为空,url带有参数", method: "DELETE", encoding: "UTF-8", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.12 content-type为xml,编码为UTF-8,DELETE,body为xml,url带有参数", method: "DELETE", encoding: "UTF-8", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}, "key": value}},

		{name: "1.13 content-type为xml,编码为UTF-8,PUT,body非正确的xml", method: "PUT", encoding: "UTF-8", contentType: "application/xml", body: []byte("<xml>"), errStr: "将<xml>转换为map失败:xml.Decoder.Token() - XML syntax error on line 1: unexpected EOF"},
		{name: "1.14 content-type为xml,编码为UTF-8,PUT,body为xml,url参数为空", method: "PUT", encoding: "UTF-8", contentType: "application/xml", body: getTestBody(value, "UTF-8", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}}},
		{name: "1.15 content-type为xml,编码为UTF-8,PUT,body为空,url带有参数", method: "PUT", encoding: "UTF-8", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.16 content-type为xml,编码为UTF-8,PUT,body为xml,url带有参数", method: "PUT", encoding: "UTF-8", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}, "key": value}},

		{name: "1.17 content-type为xml,编码为UTF-8,PATCH,body非正确的xml", method: "PATCH", encoding: "UTF-8", contentType: "application/xml", body: []byte("<xml>"), errStr: "将<xml>转换为map失败:xml.Decoder.Token() - XML syntax error on line 1: unexpected EOF"},
		{name: "1.18 content-type为xml,编码为UTF-8,PATCH,body为xml,url参数为空", method: "PATCH", encoding: "UTF-8", contentType: "application/xml", body: getTestBody(value, "UTF-8", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}}},
		{name: "1.19 content-type为xml,编码为UTF-8,PATCH,body为空,url带有参数", method: "PATCH", encoding: "UTF-8", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.20 content-type为xml,编码为UTF-8,PATCH,body为xml,url带有参数", method: "PATCH", encoding: "UTF-8", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}, "key": value}},

		{name: "2.1 content-type为xml,编码为GBK,POST,body非正确的xml", method: "POST", encoding: "GBK", contentType: "application/xml", body: []byte("<xml>"), errStr: "将<xml>转换为map失败:xml.Decoder.Token() - XML syntax error on line 1: unexpected EOF"},
		{name: "2.2 content-type为xml,编码为GBK,POST,body为xml,url参数为空", method: "POST", encoding: "GBK", contentType: "application/xml", body: getTestBody(value, "GBK", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}}},
		{name: "2.3 content-type为xml,编码为GBK,POST,body为空,url带有参数", method: "POST", encoding: "GBK", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.4 content-type为xml,编码为GBK,POST,body为xml,url带有参数", method: "POST", encoding: "GBK", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}, "key": value}},

		{name: "2.5 content-type为xml,编码为GBK,GET,body非正确的xml", method: "GET", encoding: "GBK", contentType: "application/xml", body: []byte("<xml>"), errStr: "将<xml>转换为map失败:xml.Decoder.Token() - XML syntax error on line 1: unexpected EOF"},
		{name: "2.6 content-type为xml,编码为GBK,GET,body为xml,url参数为空", method: "GET", encoding: "GBK", contentType: "application/xml", body: getTestBody(value, "GBK", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}}},
		{name: "2.7 content-type为xml,编码为GBK,GET,body为空,url带有参数", method: "GET", encoding: "GBK", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.8 content-type为xml,编码为GBK,GET,body为xml,url带有参数", method: "GET", encoding: "GBK", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}, "key": value}},

		{name: "2.9 content-type为xml,编码为GBK,DELETE,body非正确的xml", method: "DELETE", encoding: "GBK", contentType: "application/xml", body: []byte("<xml>"), errStr: "将<xml>转换为map失败:xml.Decoder.Token() - XML syntax error on line 1: unexpected EOF"},
		{name: "2.10 content-type为xml,编码为GBK,DELETE,body为xml,url参数为空", method: "DELETE", encoding: "GBK", contentType: "application/xml", body: getTestBody(value, "GBK", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}}},
		{name: "2.11 content-type为xml,编码为GBK,DELETE,body为空,url带有参数", method: "DELETE", encoding: "GBK", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.12 content-type为xml,编码为GBK,DELETE,body为xml,url带有参数", method: "DELETE", encoding: "GBK", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}, "key": value}},

		{name: "2.13 content-type为xml,编码为GBK,PUT,body非正确的xml", method: "PUT", encoding: "GBK", contentType: "application/xml", body: []byte("<xml>"), errStr: "将<xml>转换为map失败:xml.Decoder.Token() - XML syntax error on line 1: unexpected EOF"},
		{name: "2.14 content-type为xml,编码为GBK,PUT,body为xml,url参数为空", method: "PUT", encoding: "GBK", contentType: "application/xml", body: getTestBody(value, "GBK", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}}},
		{name: "2.15 content-type为xml,编码为GBK,PUT,body为空,url带有参数", method: "PUT", encoding: "GBK", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.16 content-type为xml,编码为GBK,PUT,body为xml,url带有参数", method: "PUT", encoding: "GBK", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}, "key": value}},

		{name: "1.17 content-type为xml,编码为GBK,PATCH,body非正确的xml", method: "PATCH", encoding: "GBK", contentType: "application/xml", body: []byte("<xml>"), errStr: "将<xml>转换为map失败:xml.Decoder.Token() - XML syntax error on line 1: unexpected EOF"},
		{name: "1.18 content-type为xml,编码为GBK,PATCH,body为xml,url参数为空", method: "PATCH", encoding: "GBK", contentType: "application/xml", body: getTestBody(value, "GBK", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}}},
		{name: "1.19 content-type为xml,编码为GBK,PATCH,body为空,url带有参数", method: "PATCH", encoding: "GBK", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "1.20 content-type为xml,编码为GBK,PATCH,body为xml,url带有参数", method: "PATCH", encoding: "GBK", contentType: "application/xml", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "xml"), want: map[string]interface{}{"xml": map[string]interface{}{"key": value}, "key": value}},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+tt.queryRaw, bytes.NewReader(tt.body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetMap()
		if tt.errStr != "" {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, gotS, tt.name)
		gotS2, err := w.GetMap()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}
}

func Test_body_GetMap__MIMEJSON(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		queryRaw    string
		body        []byte
		errStr      string
		want        map[string]interface{}
	}{
		{name: "1.1 content-type为json,编码为UTF-8,POST,body非正确的json", method: "POST", encoding: "UTF-8", contentType: "application/json", body: []byte("json"), errStr: "将json转换为map失败:invalid character 'j' looking for beginning of value"},
		{name: "1.2 content-type为json,编码为UTF-8,POST,body为json,url参数为空", method: "POST", encoding: "UTF-8", contentType: "application/json", body: getTestBody(value, "UTF-8", "json"), want: map[string]interface{}{"key": value}},
		{name: "1.3 content-type为json,编码为UTF-8,POST,body为空,url带有参数", method: "POST", encoding: "UTF-8", contentType: "application/json", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.4 content-type为json,编码为UTF-8,POST,body为json,url带有参数", method: "POST", encoding: "UTF-8", contentType: "application/json", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "json"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "1.5 content-type为json,编码为UTF-8,GET,body非正确的json", method: "GET", encoding: "UTF-8", contentType: "application/json", body: []byte("json"), errStr: "将json转换为map失败:invalid character 'j' looking for beginning of value"},
		{name: "1.6 content-type为json,编码为UTF-8,GET,body为json,url参数为空", method: "GET", encoding: "UTF-8", contentType: "application/json", body: getTestBody(value, "UTF-8", "json"), want: map[string]interface{}{"key": value}},
		{name: "1.7 content-type为json,编码为UTF-8,GET,body为空,url带有参数", method: "GET", encoding: "UTF-8", contentType: "application/json", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.8 content-type为json,编码为UTF-8,GET,body为json,url带有参数", method: "GET", encoding: "UTF-8", contentType: "application/json", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "json"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "1.9 content-type为json,编码为UTF-8,DELETE,body非正确的json", method: "DELETE", encoding: "UTF-8", contentType: "application/json", body: []byte("json"), errStr: "将json转换为map失败:invalid character 'j' looking for beginning of value"},
		{name: "1.10 content-type为json,编码为UTF-8,DELETE,body为json,url参数为空", method: "DELETE", encoding: "UTF-8", contentType: "application/json", body: getTestBody(value, "UTF-8", "json"), want: map[string]interface{}{"key": value}},
		{name: "1.11 content-type为json,编码为UTF-8,DELETE,body为空,url带有参数", method: "DELETE", encoding: "UTF-8", contentType: "application/json", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.12 content-type为json,编码为UTF-8,DELETE,body为json,url带有参数", method: "DELETE", encoding: "UTF-8", contentType: "application/json", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "json"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "1.13 content-type为json,编码为UTF-8,PUT,body非正确的json", method: "PUT", encoding: "UTF-8", contentType: "application/json", body: []byte("json"), errStr: "将json转换为map失败:invalid character 'j' looking for beginning of value"},
		{name: "1.14 content-type为json,编码为UTF-8,PUT,body为json,url参数为空", method: "PUT", encoding: "UTF-8", contentType: "application/json", body: getTestBody(value, "UTF-8", "json"), want: map[string]interface{}{"key": value}},
		{name: "1.15 content-type为json,编码为UTF-8,PUT,body为空,url带有参数", method: "PUT", encoding: "UTF-8", contentType: "application/json", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.16 content-type为json,编码为UTF-8,PUT,body为json,url带有参数", method: "PUT", encoding: "UTF-8", contentType: "application/json", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "json"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "1.17 content-type为json,编码为UTF-8,PATCH,body非正确的json", method: "PATCH", encoding: "UTF-8", contentType: "application/json", body: []byte("json"), errStr: "将json转换为map失败:invalid character 'j' looking for beginning of value"},
		{name: "1.18 content-type为json,编码为UTF-8,PATCH,body为json,url参数为空", method: "PATCH", encoding: "UTF-8", contentType: "application/json", body: getTestBody(value, "UTF-8", "json"), want: map[string]interface{}{"key": value}},
		{name: "1.19 content-type为json,编码为UTF-8,PATCH,body为空,url带有参数", method: "PATCH", encoding: "UTF-8", contentType: "application/json", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.20 content-type为json,编码为UTF-8,PATCH,body为json,url带有参数", method: "PATCH", encoding: "UTF-8", contentType: "application/json", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "json"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "2.1 content-type为json,编码为GBK,POST,body非正确的json", method: "POST", encoding: "GBK", contentType: "application/json", body: []byte("json"), errStr: "将json转换为map失败:invalid character 'j' looking for beginning of value"},
		{name: "2.2 content-type为json,编码为GBK,POST,body为json,url参数为空", method: "POST", encoding: "GBK", contentType: "application/json", body: getTestBody(value, "GBK", "json"), want: map[string]interface{}{"key": value}},
		{name: "2.3 content-type为json,编码为GBK,POST,body为空,url带有参数", method: "POST", encoding: "GBK", contentType: "application/json", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.4 content-type为json,编码为GBK,POST,body为json,url带有参数", method: "POST", encoding: "GBK", contentType: "application/json", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "json"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "2.5 content-type为json,编码为GBK,GET,body非正确的json", method: "GET", encoding: "GBK", contentType: "application/json", body: []byte("json"), errStr: "将json转换为map失败:invalid character 'j' looking for beginning of value"},
		{name: "2.6 content-type为json,编码为GBK,GET,body为json,url参数为空", method: "GET", encoding: "GBK", contentType: "application/json", body: getTestBody(value, "GBK", "json"), want: map[string]interface{}{"key": value}},
		{name: "2.7 content-type为json,编码为GBK,GET,body为空,url带有参数", method: "GET", encoding: "GBK", contentType: "application/json", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.8 content-type为json,编码为GBK,GET,body为json,url带有参数", method: "GET", encoding: "GBK", contentType: "application/json", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "json"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "2.9 content-type为json,编码为GBK,DELETE,body非正确的json", method: "DELETE", encoding: "GBK", contentType: "application/json", body: []byte("json"), errStr: "将json转换为map失败:invalid character 'j' looking for beginning of value"},
		{name: "2.10 content-type为json,编码为GBK,DELETE,body为json,url参数为空", method: "DELETE", encoding: "GBK", contentType: "application/json", body: getTestBody(value, "GBK", "json"), want: map[string]interface{}{"key": value}},
		{name: "2.11 content-type为json,编码为GBK,DELETE,body为空,url带有参数", method: "DELETE", encoding: "GBK", contentType: "application/json", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.12 content-type为json,编码为GBK,DELETE,body为json,url带有参数", method: "DELETE", encoding: "GBK", contentType: "application/json", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "json"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "2.13 content-type为json,编码为GBK,PUT,body非正确的json", method: "PUT", encoding: "GBK", contentType: "application/json", body: []byte("json"), errStr: "将json转换为map失败:invalid character 'j' looking for beginning of value"},
		{name: "2.14 content-type为json,编码为GBK,PUT,body为json,url参数为空", method: "PUT", encoding: "GBK", contentType: "application/json", body: getTestBody(value, "GBK", "json"), want: map[string]interface{}{"key": value}},
		{name: "2.15 content-type为json,编码为GBK,PUT,body为空,url带有参数", method: "PUT", encoding: "GBK", contentType: "application/json", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.16 content-type为json,编码为GBK,PUT,body为json,url带有参数", method: "PUT", encoding: "GBK", contentType: "application/json", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "json"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "2.17 content-type为json,编码为GBK,PATCH,body非正确的json", method: "PATCH", encoding: "GBK", contentType: "application/json", body: []byte("json"), errStr: "将json转换为map失败:invalid character 'j' looking for beginning of value"},
		{name: "2.18 content-type为json,编码为GBK,PATCH,body为json,url参数为空", method: "PATCH", encoding: "GBK", contentType: "application/json", body: getTestBody(value, "GBK", "json"), want: map[string]interface{}{"key": value}},
		{name: "2.19 content-type为json,编码为GBK,PATCH,body为空,url带有参数", method: "PATCH", encoding: "GBK", contentType: "application/json", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.20 content-type为json,编码为GBK,PATCH,body为json,url带有参数", method: "PATCH", encoding: "GBK", contentType: "application/json", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "json"), want: map[string]interface{}{"key": value + "," + value}},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+tt.queryRaw, bytes.NewReader(tt.body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetMap()
		if tt.errStr != "" {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, gotS, tt.name)
		gotS2, err := w.GetMap()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}
}

func Test_body_GetMap__MIMEYAML(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		queryRaw    string
		body        []byte
		errStr      string
		want        map[string]interface{}
	}{
		{name: "1.2 content-type为yaml,编码为UTF-8,POST,body为yaml,url参数为空", method: "POST", encoding: "UTF-8", contentType: "application/x-yaml", body: getTestBody(value, "UTF-8", "yaml"), want: map[string]interface{}{"key": value}},
		{name: "1.3 content-type为yaml,编码为UTF-8,POST,body为空,url带有参数", method: "POST", encoding: "UTF-8", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.4 content-type为yaml,编码为UTF-8,POST,body为yaml,url带有参数", method: "POST", encoding: "UTF-8", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "yaml"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "1.6 content-type为yaml,编码为UTF-8,GET,body为yaml,url参数为空", method: "GET", encoding: "UTF-8", contentType: "application/x-yaml", body: getTestBody(value, "UTF-8", "yaml"), want: map[string]interface{}{"key": value}},
		{name: "1.7 content-type为yaml,编码为UTF-8,GET,body为空,url带有参数", method: "GET", encoding: "UTF-8", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.8 content-type为yaml,编码为UTF-8,GET,body为yaml,url带有参数", method: "GET", encoding: "UTF-8", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "yaml"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "1.10 content-type为yaml,编码为UTF-8,DELETE,body为yaml,url参数为空", method: "DELETE", encoding: "UTF-8", contentType: "application/x-yaml", body: getTestBody(value, "UTF-8", "yaml"), want: map[string]interface{}{"key": value}},
		{name: "1.11 content-type为yaml,编码为UTF-8,DELETE,body为空,url带有参数", method: "DELETE", encoding: "UTF-8", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.12 content-type为yaml,编码为UTF-8,DELETE,body为yaml,url带有参数", method: "DELETE", encoding: "UTF-8", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "yaml"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "1.14 content-type为yaml,编码为UTF-8,PUT,body为yaml,url参数为空", method: "PUT", encoding: "UTF-8", contentType: "application/x-yaml", body: getTestBody(value, "UTF-8", "yaml"), want: map[string]interface{}{"key": value}},
		{name: "1.15 content-type为yaml,编码为UTF-8,PUT,body为空,url带有参数", method: "PUT", encoding: "UTF-8", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.16 content-type为yaml,编码为UTF-8,PUT,body为yaml,url带有参数", method: "PUT", encoding: "UTF-8", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "yaml"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "1.18 content-type为yaml,编码为UTF-8,PATCH,body为yaml,url参数为空", method: "PATCH", encoding: "UTF-8", contentType: "application/x-yaml", body: getTestBody(value, "UTF-8", "yaml"), want: map[string]interface{}{"key": value}},
		{name: "1.19 content-type为yaml,编码为UTF-8,PATCH,body为空,url带有参数", method: "PATCH", encoding: "UTF-8", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.20 content-type为yaml,编码为UTF-8,PATCH,body为yaml,url带有参数", method: "PATCH", encoding: "UTF-8", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "yaml"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "2.2 content-type为yaml,编码为GBK,POST,body为yaml,url参数为空", method: "POST", encoding: "GBK", contentType: "application/x-yaml", body: getTestBody(value, "GBK", "yaml"), want: map[string]interface{}{"key": value}},
		{name: "2.3 content-type为yaml,编码为GBK,POST,body为空,url带有参数", method: "POST", encoding: "GBK", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.4 content-type为yaml,编码为GBK,POST,body为yaml,url带有参数", method: "POST", encoding: "GBK", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "yaml"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "2.6 content-type为yaml,编码为GBK,GET,body为yaml,url参数为空", method: "GET", encoding: "GBK", contentType: "application/x-yaml", body: getTestBody(value, "GBK", "yaml"), want: map[string]interface{}{"key": value}},
		{name: "2.7 content-type为yaml,编码为GBK,GET,body为空,url带有参数", method: "GET", encoding: "GBK", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.8 content-type为yaml,编码为GBK,GET,body为yaml,url带有参数", method: "GET", encoding: "GBK", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "yaml"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "2.10 content-type为yaml,编码为GBK,DELETE,body为yaml,url参数为空", method: "DELETE", encoding: "GBK", contentType: "application/x-yaml", body: getTestBody(value, "GBK", "yaml"), want: map[string]interface{}{"key": value}},
		{name: "2.11 content-type为yaml,编码为GBK,DELETE,body为空,url带有参数", method: "DELETE", encoding: "GBK", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.12 content-type为yaml,编码为GBK,DELETE,body为yaml,url带有参数", method: "DELETE", encoding: "GBK", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "yaml"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "2.14 content-type为yaml,编码为GBK,PUT,body为yaml,url参数为空", method: "PUT", encoding: "GBK", contentType: "application/x-yaml", body: getTestBody(value, "GBK", "yaml"), want: map[string]interface{}{"key": value}},
		{name: "2.15 content-type为yaml,编码为GBK,PUT,body为空,url带有参数", method: "PUT", encoding: "GBK", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.16 content-type为yaml,编码为GBK,PUT,body为yaml,url带有参数", method: "PUT", encoding: "GBK", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "yaml"), want: map[string]interface{}{"key": value + "," + value}},

		{name: "2.18 content-type为yaml,编码为GBK,PATCH,body为yaml,url参数为空", method: "PATCH", encoding: "GBK", contentType: "application/x-yaml", body: getTestBody(value, "GBK", "yaml"), want: map[string]interface{}{"key": value}},
		{name: "2.19 content-type为yaml,编码为GBK,PATCH,body为空,url带有参数", method: "PATCH", encoding: "GBK", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.20 content-type为yaml,编码为GBK,PATCH,body为yaml,url带有参数", method: "PATCH", encoding: "GBK", contentType: "application/x-yaml", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "yaml"), want: map[string]interface{}{"key": value + "," + value}},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+tt.queryRaw, bytes.NewReader(tt.body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r

		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetMap()
		if tt.errStr != "" {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, gotS, tt.name)
		gotS2, err := w.GetMap()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}
}

func Test_body_GetMap__MIMEPOSTForm(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		queryRaw    string
		body        []byte
		errStr      string
		want        map[string]interface{}
	}{
		{name: "1.1 content-type为form,编码为UTF-8,GET,body为空,url为空", method: "POST", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", body: nil, want: map[string]interface{}{}},
		{name: "1.2 content-type为form,编码为UTF-8,POST,body为form,url参数为空", method: "POST", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", body: getTestBody(value, "UTF-8", "form"), want: map[string]interface{}{"key": value}},
		{name: "1.3 content-type为form,编码为UTF-8,POST,body为空,url带有参数", method: "POST", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.4 content-type为form,编码为UTF-8,POST,body为form,url带有参数", method: "POST", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "form"), want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "1.5 content-type为form,编码为UTF-8,GET,body为空,url为空", method: "GET", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", body: nil, want: map[string]interface{}{}},
		{name: "1.6 content-type为form,编码为UTF-8,GET,body为form,url参数为空", method: "GET", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", body: getTestBody(value, "UTF-8", "form"), want: map[string]interface{}{"key": value}},
		{name: "1.7 content-type为form,编码为UTF-8,GET,body为空,url带有参数", method: "GET", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.8 content-type为form,编码为UTF-8,GET,body为form,url带有参数", method: "GET", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "form"), want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "1.9 content-type为form,编码为UTF-8,DELETE,body为空,url为空", method: "DELETE", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", body: nil, want: map[string]interface{}{}},
		{name: "1.10 content-type为form,编码为UTF-8,DELETE,body为form,url参数为空", method: "DELETE", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", body: getTestBody(value, "UTF-8", "form"), want: map[string]interface{}{"key": value}},
		{name: "1.11 content-type为form,编码为UTF-8,DELETE,body为空,url带有参数", method: "DELETE", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.12 content-type为form,编码为UTF-8,DELETE,body为form,url带有参数", method: "DELETE", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "form"), want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "1.13 content-type为form,编码为UTF-8,PUT,body为空,url为空", method: "PUT", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", body: nil, want: map[string]interface{}{}},
		{name: "1.14 content-type为form,编码为UTF-8,PUT,body为form,url参数为空", method: "PUT", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", body: getTestBody(value, "UTF-8", "form"), want: map[string]interface{}{"key": value}},
		{name: "1.15 content-type为form,编码为UTF-8,PUT,body为空,url带有参数", method: "PUT", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.16 content-type为form,编码为UTF-8,PUT,body为form,url带有参数", method: "PUT", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "form"), want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "1.17 content-type为form,编码为UTF-8,GET,body为空,url为空", method: "PATCH", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", body: nil, want: map[string]interface{}{}},
		{name: "1.18 content-type为form,编码为UTF-8,PATCH,body为form,url参数为空", method: "PATCH", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", body: getTestBody(value, "UTF-8", "form"), want: map[string]interface{}{"key": value}},
		{name: "1.19 content-type为form,编码为UTF-8,PATCH,body为空,url带有参数", method: "PATCH", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.20 content-type为form,编码为UTF-8,PATCH,body为form,url带有参数", method: "PATCH", encoding: "UTF-8", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "UTF-8"), body: getTestBody(value, "UTF-8", "form"), want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "2.1 content-type为form,编码为GBK,GET,body为空,url为空", method: "POST", encoding: "GBK", contentType: "application/x-www-form-urlencoded", body: nil, want: map[string]interface{}{}},
		{name: "2.2 content-type为form,编码为GBK,POST,body为form,url参数为空", method: "POST", encoding: "GBK", contentType: "application/x-www-form-urlencoded", body: getTestBody(value, "GBK", "form"), want: map[string]interface{}{"key": value}},
		{name: "2.3 content-type为form,编码为GBK,POST,body为空,url带有参数", method: "POST", encoding: "GBK", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.4 content-type为form,编码为GBK,POST,body为form,url带有参数", method: "POST", encoding: "GBK", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "form"), want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "2.5 content-type为form,编码为GBK,GET,body为空,url为空", method: "GET", encoding: "GBK", contentType: "application/x-www-form-urlencoded", body: nil, want: map[string]interface{}{}},
		{name: "2.6 content-type为form,编码为GBK,GET,body为form,url参数为空", method: "GET", encoding: "GBK", contentType: "application/x-www-form-urlencoded", body: getTestBody(value, "GBK", "form"), want: map[string]interface{}{"key": value}},
		{name: "2.7 content-type为form,编码为GBK,GET,body为空,url带有参数", method: "GET", encoding: "GBK", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.8 content-type为form,编码为GBK,GET,body为form,url带有参数", method: "GET", encoding: "GBK", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "form"), want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "2.9 content-type为form,编码为GBK,DELETE,body为空,url为空", method: "DELETE", encoding: "GBK", contentType: "application/x-www-form-urlencoded", body: nil, want: map[string]interface{}{}},
		{name: "2.10 content-type为form,编码为GBK,DELETE,body为form,url参数为空", method: "DELETE", encoding: "GBK", contentType: "application/x-www-form-urlencoded", body: getTestBody(value, "GBK", "form"), want: map[string]interface{}{"key": value}},
		{name: "2.11 content-type为form,编码为GBK,DELETE,body为空,url带有参数", method: "DELETE", encoding: "GBK", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.12 content-type为form,编码为GBK,DELETE,body为form,url带有参数", method: "DELETE", encoding: "GBK", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "form"), want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "2.13 content-type为form,编码为GBK,PUT,body为空,url为空", method: "PUT", encoding: "GBK", contentType: "application/x-www-form-urlencoded", body: nil, want: map[string]interface{}{}},
		{name: "2.14 content-type为form,编码为GBK,PUT,body为form,url参数为空", method: "PUT", encoding: "GBK", contentType: "application/x-www-form-urlencoded", body: getTestBody(value, "GBK", "form"), want: map[string]interface{}{"key": value}},
		{name: "2.15 content-type为form,编码为GBK,PUT,body为空,url带有参数", method: "PUT", encoding: "GBK", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.16 content-type为form,编码为GBK,PUT,body为form,url带有参数", method: "PUT", encoding: "GBK", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "form"), want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "2.17 content-type为form,编码为GBK,PATCH,body为空,url为空", method: "PATCH", encoding: "GBK", contentType: "application/x-www-form-urlencoded", body: nil, want: map[string]interface{}{}},
		{name: "2.18 content-type为form,编码为GBK,PATCH,body为form,url参数为空", method: "PATCH", encoding: "GBK", contentType: "application/x-www-form-urlencoded", body: getTestBody(value, "GBK", "form"), want: map[string]interface{}{"key": value}},
		{name: "2.19 content-type为form,编码为GBK,PATCH,body为空,url带有参数", method: "PATCH", encoding: "GBK", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.20 content-type为form,编码为GBK,PATCH,body为form,url带有参数", method: "PATCH", encoding: "GBK", contentType: "application/x-www-form-urlencoded", queryRaw: getTestQueryRaw(value, "GBK"), body: getTestBody(value, "GBK", "form"), want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+tt.queryRaw, bytes.NewReader(tt.body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", tt.contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r
		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetMap()
		if tt.errStr != "" {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, gotS, tt.name)
		gotS2, err := w.GetMap()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}
}

func Test_body_GetMap__MIMEMultipartPOSTForm(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		encoding    string
		contentType string
		queryRaw    string
		value       url.Values
		errStr      string
		want        map[string]interface{}
	}{

		{name: "1.1 content-type为form-data,编码为UTF-8,GET,body为空,url为空", method: "POST", encoding: "UTF-8", value: nil, want: map[string]interface{}{}},
		{name: "1.2 content-type为form-data,编码为UTF-8,POST,body为form,url参数为空", method: "POST", encoding: "UTF-8", value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": value}},
		{name: "1.3 content-type为form-data,编码为UTF-8,POST,body为空,url带有参数", method: "POST", encoding: "UTF-8", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.4 content-type为form-data,编码为UTF-8,POST,body为form,url带有参数", method: "POST", encoding: "UTF-8", queryRaw: getTestQueryRaw(value, "UTF-8"), value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "1.5 content-type为form-data,编码为UTF-8,GET,body为空,url为空", method: "GET", encoding: "UTF-8", value: nil, want: map[string]interface{}{}},
		{name: "1.6 content-type为form-data,编码为UTF-8,GET,body为form,url参数为空", method: "GET", encoding: "UTF-8", value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": value}},
		{name: "1.7 content-type为form-data,编码为UTF-8,GET,body为空,url带有参数", method: "GET", encoding: "UTF-8", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.8 content-type为form-data,编码为UTF-8,GET,body为form,url带有参数", method: "GET", encoding: "UTF-8", queryRaw: getTestQueryRaw(value, "UTF-8"), value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "1.9 content-type为form-data,编码为UTF-8,DELETE,body为空,url为空", method: "DELETE", encoding: "UTF-8", value: nil, want: map[string]interface{}{}},
		{name: "1.10 content-type为form-data,编码为UTF-8,DELETE,body为form,url参数为空", method: "DELETE", encoding: "UTF-8", value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": value}},
		{name: "1.11 content-type为form-data,编码为UTF-8,DELETE,body为空,url带有参数", method: "DELETE", encoding: "UTF-8", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.12 content-type为form-data,编码为UTF-8,DELETE,body为form,url带有参数", method: "DELETE", encoding: "UTF-8", queryRaw: getTestQueryRaw(value, "UTF-8"), value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "1.13 content-type为form-data,编码为UTF-8,PUT,body为空,url为空", method: "PUT", encoding: "UTF-8", value: nil, want: map[string]interface{}{}},
		{name: "1.14 content-type为form-data,编码为UTF-8,PUT,body为form,url参数为空", method: "PUT", encoding: "UTF-8", value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": value}},
		{name: "1.15 content-type为form-data,编码为UTF-8,PUT,body为空,url带有参数", method: "PUT", encoding: "UTF-8", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.16 content-type为form-data,编码为UTF-8,PUT,body为form,url带有参数", method: "PUT", encoding: "UTF-8", queryRaw: getTestQueryRaw(value, "UTF-8"), value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "1.17 content-type为form-data,编码为UTF-8,GET,body为空,url为空", method: "PATCH", encoding: "UTF-8", value: nil, want: map[string]interface{}{}},
		{name: "1.18 content-type为form-data,编码为UTF-8,PATCH,body为form,url参数为空", method: "PATCH", encoding: "UTF-8", value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": value}},
		{name: "1.19 content-type为form-data,编码为UTF-8,PATCH,body为空,url带有参数", method: "PATCH", encoding: "UTF-8", queryRaw: getTestQueryRaw(value, "UTF-8"), want: map[string]interface{}{"key": value}},
		{name: "1.20 content-type为form-data,编码为UTF-8,PATCH,body为form,url带有参数", method: "PATCH", encoding: "UTF-8", queryRaw: getTestQueryRaw(value, "UTF-8"), value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "2.1 content-type为form-data,编码为GBK,GET,body为空,url为空", method: "POST", encoding: "GBK", value: nil, want: map[string]interface{}{}},
		{name: "2.2 content-type为form-data,编码为GBK,POST,body为form,url参数为空", method: "POST", encoding: "GBK", value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": value}},
		{name: "2.3 content-type为form-data,编码为GBK,POST,body为空,url带有参数", method: "POST", encoding: "GBK", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.4 content-type为form-data,编码为GBK,POST,body为form,url带有参数", method: "POST", encoding: "GBK", queryRaw: getTestQueryRaw(value, "GBK"), value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "2.5 content-type为form-data,编码为GBK,GET,body为空,url为空", method: "GET", encoding: "GBK", value: nil, want: map[string]interface{}{}},
		{name: "2.6 content-type为form-data,编码为GBK,GET,body为form,url参数为空", method: "GET", encoding: "GBK", value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": value}},
		{name: "2.7 content-type为form-data,编码为GBK,GET,body为空,url带有参数", method: "GET", encoding: "GBK", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.8 content-type为form-data,编码为GBK,GET,body为form,url带有参数", method: "GET", encoding: "GBK", queryRaw: getTestQueryRaw(value, "GBK"), value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "2.9 content-type为form-data,编码为GBK,DELETE,body为空,url为空", method: "DELETE", encoding: "GBK", value: nil, want: map[string]interface{}{}},
		{name: "2.10 content-type为form-data,编码为GBK,DELETE,body为form,url参数为空", method: "DELETE", encoding: "GBK", value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": value}},
		{name: "2.11 content-type为form-data,编码为GBK,DELETE,body为空,url带有参数", method: "DELETE", encoding: "GBK", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.12 content-type为form-data,编码为GBK,DELETE,body为form,url带有参数", method: "DELETE", encoding: "GBK", queryRaw: getTestQueryRaw(value, "GBK"), value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "2.13 content-type为form-data,编码为GBK,PUT,body为空,url为空", method: "PUT", encoding: "GBK", value: nil, want: map[string]interface{}{}},
		{name: "2.14 content-type为form-data,编码为GBK,PUT,body为form,url参数为空", method: "PUT", encoding: "GBK", value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": value}},
		{name: "2.15 content-type为form-data,编码为GBK,PUT,body为空,url带有参数", method: "PUT", encoding: "GBK", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.16 content-type为form-data,编码为GBK,PUT,body为form,url带有参数", method: "PUT", encoding: "GBK", queryRaw: getTestQueryRaw(value, "GBK"), value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},

		{name: "2.17 content-type为form-data,编码为GBK,PATCH,body为空,url为空", method: "PATCH", encoding: "GBK", value: nil, want: map[string]interface{}{}},
		{name: "2.18 content-type为form-data,编码为GBK,PATCH,body为form,url参数为空", method: "PATCH", encoding: "GBK", value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": value}},
		{name: "2.19 content-type为form-data,编码为GBK,PATCH,body为空,url带有参数", method: "PATCH", encoding: "GBK", queryRaw: getTestQueryRaw(value, "GBK"), want: map[string]interface{}{"key": value}},
		{name: "2.20 content-type为form-data,编码为GBK,PATCH,body为form,url带有参数", method: "PATCH", encoding: "GBK", queryRaw: getTestQueryRaw(value, "GBK"), value: url.Values{"key": []string{value}}, want: map[string]interface{}{"key": fmt.Sprintf("%s,%s", value, value)}},
	}

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		//写入文件和设置参数
		content, contentType := getTestMIMEMultipartPOSTForm(tt.value)
		data, _ := encoding.Encode(content, tt.encoding)
		body := string(data)

		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url?"+tt.queryRaw, bytes.NewReader([]byte(body)))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, tt.encoding))

		//替换gin上下文的请求
		c.Request = r
		w := ctx.NewBody(middleware.NewGinCtx(c), tt.encoding)
		gotS, err := w.GetMap()
		if tt.errStr != "" {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, gotS, tt.name)
		gotS2, err := w.GetMap()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}
}
