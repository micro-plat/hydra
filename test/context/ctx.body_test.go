package context

import (
	"bytes"
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
)

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

func Test_body_GetRawBody(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
		want        string
	}{
		{name: "1.1 content-type为application/xml,方法为POST", contentType: "application/xml", method: "POST", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.2 content-type为application/xml,方法为PUT", contentType: "application/xml", method: "PUT", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.3 content-type为application/xml,方法为DELETE", contentType: "application/xml", method: "DELETE", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.4 content-type为application/xml,方法为GET", contentType: "application/xml", method: "GET", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.5 content-type为application/xml,方法为PATCH", contentType: "application/xml", method: "PATCH", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},

		{name: "2.1 content-type为text/xml,方法为POST", contentType: "text/xml", method: "POST", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.2 content-type为text/xml,方法为PUT", contentType: "text/xml", method: "PUT", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.3 content-type为text/xml,方法为DELETE", contentType: "text/xml", method: "DELETE", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.4 content-type为text/xml,方法为GET", contentType: "text/xml", method: "GET", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.5 content-type为text/xml,方法为PATCH", contentType: "text/xml", method: "PATCH", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},

		{name: "3.1 content-type为application/json,方法为POST", contentType: "application/json", method: "POST", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.2 content-type为application/json,方法为PUT", contentType: "application/json", method: "PUT", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.3 content-type为application/json,方法为DELETE", contentType: "application/json", method: "DELETE", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.4 content-type为application/json,方法为GET", contentType: "application/json", method: "GET", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.5 content-type为application/json,方法为PATCH", contentType: "application/json", method: "PATCH", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},

		{name: "4.1 content-type为text/json,方法为POST", contentType: "text/json", method: "POST", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.2 content-type为text/json,方法为PUT", contentType: "text/json", method: "PUT", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.3 content-type为text/json,方法为DELETE", contentType: "json/json", method: "DELETE", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.4 content-type为text/json,方法为GET", contentType: "text/json", method: "GET", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.5 content-type为text/json,方法为PATCH", contentType: "text/json", method: "PATCH", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},

		{name: "5.1 content-type为application/x-yaml,方法为POST", contentType: "application/x-yaml", method: "POST", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.2 content-type为application/x-yaml,方法为PUT", contentType: "application/x-yaml", method: "PUT", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.3 content-type为application/x-yaml,方法为DELETE", contentType: "application/x-yaml", method: "DELETE", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.4 content-type为application/x-yaml,方法为GET", contentType: "application/x-yaml", method: "GET", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.5 content-type为application/x-yaml,方法为PATCH", contentType: "application/x-yaml", method: "PATCH", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},

		{name: "6.1 content-type为text/plain,方法为POST", contentType: "text/plain", method: "POST", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.2 content-type为text/plain,方法为PUT", contentType: "text/plain", method: "PUT", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.3 content-type为text/plain,方法为DELETE", contentType: "text/plain", method: "DELETE", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.4 content-type为text/plain,方法为GET", contentType: "text/plain", method: "GET", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.5 content-type为text/plain,方法为PATCH", contentType: "text/plain", method: "PATCH", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},

		{name: "7.1 content-type为application/x-www-form-urlencoded,方法为POST", contentType: "application/x-www-form-urlencoded", method: "POST", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.2 content-type为application/x-www-form-urlencoded,方法为PUT", contentType: "application/x-www-form-urlencoded", method: "PUT", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.3 content-type为application/x-www-form-urlencoded,方法为DELETE", contentType: "application/x-www-form-urlencoded", method: "DELETE", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.4 content-type为application/x-www-form-urlencoded,方法为GET", contentType: "application/x-www-form-urlencoded", method: "GET", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.5 content-type为application/x-www-form-urlencoded,方法为PATCH", contentType: "application/x-www-form-urlencoded", method: "GET", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},

		{name: "8.1 content-type为multipart/form-data,方法为POST", contentType: "multipart/form-data", method: "POST", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.2 content-type为multipart/form-data,方法为PUT", contentType: "multipart/form-data", method: "PUT", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.3 content-type为multipart/form-data,方法为DELETE", contentType: "multipart/form-data", method: "DELETE", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.4 content-type为multipart/form-data,方法为GET", contentType: "multipart/form-data", method: "GET", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.5 content-type为multipart/form-data,方法为PATCH", contentType: "multipart/form-data", method: "PATCH", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},

		{name: "9.1 content-type为text/html,方法为POST", contentType: "text/plain", method: "POST", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.2 content-type为text/html,方法为PUT", contentType: "text/plain", method: "PUT", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.3 content-type为text/html,方法为DELETE", contentType: "text/plain", method: "DELETE", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.4 content-type为text/html,方法为GET", contentType: "text/plain", method: "GET", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.4 content-type为text/html,方法为PATCH", contentType: "text/plain", method: "PATCH", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
	}

	router := gin.New()
	router.POST("/url", func(c *gin.Context) {
		return
	})

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
		router.HandleContext(c)

		w := ctx.NewBody(middleware.NewGinCtx(c), "utf-8")
		gotS, err := w.GetRawBody()
		str := string(gotS)
		if tt.contentType == "multipart/form-data" {
			s := strings.Split(str, "\r\n")
			s = s[1 : len(s)-2]
			str = strings.Join(s, "\r\n")
		}
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, str, tt.name)
		gotS2, err := w.GetRawBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}

func Test_body_GetBody_UTF8(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
		want        string
	}{
		{name: "1.1 content-type为application/xml,方法为POST", contentType: "application/xml", method: "POST", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.2 content-type为application/xml,方法为PUT", contentType: "application/xml", method: "PUT", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.3 content-type为application/xml,方法为DELETE", contentType: "application/xml", method: "DELETE", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.4 content-type为application/xml,方法为GET", contentType: "application/xml", method: "GET", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.5 content-type为application/xml,方法为PATCH", contentType: "application/xml", method: "PATCH", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},

		{name: "2.1 content-type为text/xml,方法为POST", contentType: "text/xml", method: "POST", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.2 content-type为text/xml,方法为PUT", contentType: "text/xml", method: "PUT", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.3 content-type为text/xml,方法为DELETE", contentType: "text/xml", method: "DELETE", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.4 content-type为text/xml,方法为GET", contentType: "text/xml", method: "GET", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.5 content-type为text/xml,方法为PATCH", contentType: "text/xml", method: "PATCH", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},

		{name: "3.1 content-type为application/json,方法为POST", contentType: "application/json", method: "POST", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.2 content-type为application/json,方法为PUT", contentType: "application/json", method: "PUT", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.3 content-type为application/json,方法为DELETE", contentType: "application/json", method: "DELETE", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.4 content-type为application/json,方法为GET", contentType: "application/json", method: "GET", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.5 content-type为application/json,方法为PATCH", contentType: "application/json", method: "PATCH", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},

		{name: "4.1 content-type为text/json,方法为POST", contentType: "text/json", method: "POST", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.2 content-type为text/json,方法为PUT", contentType: "text/json", method: "PUT", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.3 content-type为text/json,方法为DELETE", contentType: "json/json", method: "DELETE", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.4 content-type为text/json,方法为GET", contentType: "text/json", method: "GET", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.5 content-type为text/json,方法为PATCH", contentType: "text/json", method: "PATCH", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},

		{name: "5.1 content-type为application/x-yaml,方法为POST", contentType: "application/x-yaml", method: "POST", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.2 content-type为application/x-yaml,方法为PUT", contentType: "application/x-yaml", method: "PUT", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.3 content-type为application/x-yaml,方法为DELETE", contentType: "application/x-yaml", method: "DELETE", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.4 content-type为application/x-yaml,方法为GET", contentType: "application/x-yaml", method: "GET", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.5 content-type为application/x-yaml,方法为PATCH", contentType: "application/x-yaml", method: "PATCH", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},

		{name: "6.1 content-type为text/plain,方法为POST", contentType: "text/plain", method: "POST", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.2 content-type为text/plain,方法为PUT", contentType: "text/plain", method: "PUT", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.3 content-type为text/plain,方法为DELETE", contentType: "text/plain", method: "DELETE", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.4 content-type为text/plain,方法为GET", contentType: "text/plain", method: "GET", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.5 content-type为text/plain,方法为PATCH", contentType: "text/plain", method: "PATCH", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},

		{name: "7.1 content-type为application/x-www-form-urlencoded,方法为POST", contentType: "application/x-www-form-urlencoded", method: "POST", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.2 content-type为application/x-www-form-urlencoded,方法为PUT", contentType: "application/x-www-form-urlencoded", method: "PUT", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.3 content-type为application/x-www-form-urlencoded,方法为DELETE", contentType: "application/x-www-form-urlencoded", method: "DELETE", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.4 content-type为application/x-www-form-urlencoded,方法为GET", contentType: "application/x-www-form-urlencoded", method: "GET", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.5 content-type为application/x-www-form-urlencoded,方法为PATCH", contentType: "application/x-www-form-urlencoded", method: "GET", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},

		{name: "8.1 content-type为multipart/form-data,方法为POST", contentType: "multipart/form-data", method: "POST", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.2 content-type为multipart/form-data,方法为PUT", contentType: "multipart/form-data", method: "PUT", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.3 content-type为multipart/form-data,方法为DELETE", contentType: "multipart/form-data", method: "DELETE", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.4 content-type为multipart/form-data,方法为GET", contentType: "multipart/form-data", method: "GET", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.5 content-type为multipart/form-data,方法为PATCH", contentType: "multipart/form-data", method: "PATCH", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},

		{name: "9.1 content-type为text/html,方法为POST", contentType: "text/plain", method: "POST", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.2 content-type为text/html,方法为PUT", contentType: "text/plain", method: "PUT", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.3 content-type为text/html,方法为DELETE", contentType: "text/plain", method: "DELETE", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.4 content-type为text/html,方法为GET", contentType: "text/plain", method: "GET", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.4 content-type为text/html,方法为PATCH", contentType: "text/plain", method: "PATCH", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
	}

	router := gin.New()
	router.POST("/url", func(c *gin.Context) {
		return
	})

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "/url", strings.NewReader(tt.body))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", tt.contentType))

		//替换gin上下文的请求
		c.Request = r
		router.HandleContext(c)

		w := ctx.NewBody(middleware.NewGinCtx(c), "utf-8")
		gotS, err := w.GetBody()
		str := string(gotS)
		if tt.contentType == "multipart/form-data" {
			s := strings.Split(str, "\r\n")
			s = s[1 : len(s)-2]
			str = strings.Join(s, "\r\n")
		}
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, tt.want, str, tt.name)
		gotS2, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}

func Test_body_GetBody_GBK(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
		want        string
	}{
		{name: "1.1 content-type为application/xml,方法为POST", contentType: "application/xml", method: "POST", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.2 content-type为application/xml,方法为PUT", contentType: "application/xml", method: "PUT", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.3 content-type为application/xml,方法为DELETE", contentType: "application/xml", method: "DELETE", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.4 content-type为application/xml,方法为GET", contentType: "application/xml", method: "GET", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "1.5 content-type为application/xml,方法为PATCH", contentType: "application/xml", method: "PATCH", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},

		{name: "2.1 content-type为text/xml,方法为POST", contentType: "text/xml", method: "POST", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.2 content-type为text/xml,方法为PUT", contentType: "text/xml", method: "PUT", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.3 content-type为text/xml,方法为DELETE", contentType: "text/xml", method: "DELETE", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.4 content-type为text/xml,方法为GET", contentType: "text/xml", method: "GET", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},
		{name: "2.5 content-type为text/xml,方法为PATCH", contentType: "text/xml", method: "PATCH", body: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`, want: `<?xml version="1.0"?><key>value中文~!@#$%^&*()_+{}|:"<>?</key>`},

		{name: "3.1 content-type为application/json,方法为POST", contentType: "application/json", method: "POST", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.2 content-type为application/json,方法为PUT", contentType: "application/json", method: "PUT", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.3 content-type为application/json,方法为DELETE", contentType: "application/json", method: "DELETE", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.4 content-type为application/json,方法为GET", contentType: "application/json", method: "GET", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "3.5 content-type为application/json,方法为PATCH", contentType: "application/json", method: "PATCH", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},

		{name: "4.1 content-type为text/json,方法为POST", contentType: "text/json", method: "POST", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.2 content-type为text/json,方法为PUT", contentType: "text/json", method: "PUT", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.3 content-type为text/json,方法为DELETE", contentType: "json/json", method: "DELETE", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.4 content-type为text/json,方法为GET", contentType: "text/json", method: "GET", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},
		{name: "4.5 content-type为text/json,方法为PATCH", contentType: "text/json", method: "PATCH", body: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`, want: `{"key":"value中文~!@#$%^&*()_+{}|:"<>?"}`},

		{name: "5.1 content-type为application/x-yaml,方法为POST", contentType: "application/x-yaml", method: "POST", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.2 content-type为application/x-yaml,方法为PUT", contentType: "application/x-yaml", method: "PUT", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.3 content-type为application/x-yaml,方法为DELETE", contentType: "application/x-yaml", method: "DELETE", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.4 content-type为application/x-yaml,方法为GET", contentType: "application/x-yaml", method: "GET", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "5.5 content-type为application/x-yaml,方法为PATCH", contentType: "application/x-yaml", method: "PATCH", body: `key: value中文~!@#$%^&*()_+{}|:"<>?`, want: `key: value中文~!@#$%^&*()_+{}|:"<>?`},

		{name: "6.1 content-type为text/plain,方法为POST", contentType: "text/plain", method: "POST", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.2 content-type为text/plain,方法为PUT", contentType: "text/plain", method: "PUT", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.3 content-type为text/plain,方法为DELETE", contentType: "text/plain", method: "DELETE", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.4 content-type为text/plain,方法为GET", contentType: "text/plain", method: "GET", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "6.5 content-type为text/plain,方法为PATCH", contentType: "text/plain", method: "PATCH", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},

		{name: "7.1 content-type为application/x-www-form-urlencoded,方法为POST", contentType: "application/x-www-form-urlencoded", method: "POST", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.2 content-type为application/x-www-form-urlencoded,方法为PUT", contentType: "application/x-www-form-urlencoded", method: "PUT", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.3 content-type为application/x-www-form-urlencoded,方法为DELETE", contentType: "application/x-www-form-urlencoded", method: "DELETE", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.4 content-type为application/x-www-form-urlencoded,方法为GET", contentType: "application/x-www-form-urlencoded", method: "GET", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "7.5 content-type为application/x-www-form-urlencoded,方法为PATCH", contentType: "application/x-www-form-urlencoded", method: "GET", body: `key=value中文~!@#$%^&*()_+{}|:"<>?`, want: `key=value中文~!@#$%^&*()_+{}|:"<>?`},

		{name: "8.1 content-type为multipart/form-data,方法为POST", contentType: "multipart/form-data", method: "POST", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.2 content-type为multipart/form-data,方法为PUT", contentType: "multipart/form-data", method: "PUT", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.3 content-type为multipart/form-data,方法为DELETE", contentType: "multipart/form-data", method: "DELETE", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.4 content-type为multipart/form-data,方法为GET", contentType: "multipart/form-data", method: "GET", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},
		{name: "8.5 content-type为multipart/form-data,方法为PATCH", contentType: "multipart/form-data", method: "PATCH", body: getTestMIMEMultipartPOSTForm(), want: getUploadBody()},

		{name: "9.1 content-type为text/html,方法为POST", contentType: "text/plain", method: "POST", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.2 content-type为text/html,方法为PUT", contentType: "text/plain", method: "PUT", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.3 content-type为text/html,方法为DELETE", contentType: "text/plain", method: "DELETE", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.4 content-type为text/html,方法为GET", contentType: "text/plain", method: "GET", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
		{name: "9.4 content-type为text/html,方法为PATCH", contentType: "text/plain", method: "PATCH", body: `value中文~!@#$%^&*()_+{}|:"<>?`, want: `value中文~!@#$%^&*()_+{}|:"<>?`},
	}

	router := gin.New()
	router.POST("/url", func(c *gin.Context) {
		return
	})

	for _, tt := range tests {
		//构建上下文
		c := &gin.Context{}

		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "/url", strings.NewReader(url.QueryEscape(tt.body)))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s; charset=gbk", tt.contentType))

		//替换gin上下文的请求
		c.Request = r
		router.HandleContext(c)

		w := ctx.NewBody(middleware.NewGinCtx(c), "gbk")
		gotS, err := w.GetBody()
		str := string(gotS)
		if tt.contentType == "multipart/form-data" {
			s := strings.Split(str, "\r\n")
			s = s[1 : len(s)-2]
			str = strings.Join(s, "\r\n")
		}
		fmt.Print(str)
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, Utf8ToGbk(tt.want), str, tt.name)
		gotS2, err := w.GetBody()
		assert.Equal(t, nil, err, tt.name)
		assert.Equal(t, gotS, gotS2, tt.name+"再次读取body")
	}

}
