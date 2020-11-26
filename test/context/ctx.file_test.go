package context

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
)

func Test_file_GetFileSize(t *testing.T) {
	tests := []struct {
		name        string
		encoding    string
		method      string
		contentType string
		fileKey     string
		wantbody    string
		wantSize    int64
		wantName    string
		wantErr     string
	}{
		{name: "1.1 content-type为multipart/form-data,POST", contentType: "multipart/form-data", method: "POST", fileKey: "upload", wantbody: "ADASDASDASFHNOJM~!@#$%^&*", wantName: "upload.test.txt", wantSize: 25},
		{name: "1.2 content-type为multipart/form-data,GET", contentType: "multipart/form-data", method: "GET", fileKey: "upload", wantbody: "ADASDASDASFHNOJM~!@#$%^&*", wantName: "upload.test.txt", wantSize: 25},
		{name: "1.3 content-type为multipart/form-data,DELETE", contentType: "multipart/form-data", method: "POST", fileKey: "upload", wantbody: "ADASDASDASFHNOJM~!@#$%^&*", wantName: "upload.test.txt", wantSize: 25},
		{name: "1.4 content-type为multipart/form-data,PUT", contentType: "multipart/form-data", method: "POST", fileKey: "upload", wantbody: "ADASDASDASFHNOJM~!@#$%^&*", wantName: "upload.test.txt", wantSize: 25},
		{name: "1.5 content-type为multipart/form-data,PATCH", contentType: "multipart/form-data", method: "POST", fileKey: "upload", wantbody: "ADASDASDASFHNOJM~!@#$%^&*", wantName: "upload.test.txt", wantSize: 25},

		{name: "2.1 content-type为application/x-www-form-urlencoded,POST", contentType: "application/x-www-form-urlencoded", method: "POST", fileKey: "upload", wantErr: "request Content-Type isn't multipart/form-data"},
		{name: "3.1 content-type为application/json,POST", contentType: "application/json", method: "POST", fileKey: "upload", wantErr: "request Content-Type isn't multipart/form-data"},
		{name: "4.1 content-type为application/xml,POST", contentType: "application/xml", method: "POST", fileKey: "upload", wantErr: "request Content-Type isn't multipart/form-data"},
	}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	for _, tt := range tests {
		boundary := multipart.NewWriter(bytes.NewBufferString("")).Boundary()
		data := "--" + boundary + "\r\n" + getUploadBody() + "\r\n" + "--" + boundary + "--"
		//构建请求 方法要与注册方法一致
		r, err := http.NewRequest(tt.method, "http://localhost:8080/url", bytes.NewReader([]byte(data)))
		assert.Equal(t, nil, err, "构建请求")
		//设置content-type
		r.Header.Set("Content-Type", fmt.Sprintf("%s;boundary=%s", tt.contentType, boundary))
		c.Request = r
		w := ctx.NewFile(middleware.NewGinCtx(c), conf.NewMeta())
		f, err := w.GetFileBody(tt.fileKey)
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, nil, err, "文件上传请求错误")
		s, _ := ioutil.ReadAll(f)
		assert.Equal(t, tt.wantbody, string(s), tt.name)

		name, _ := w.GetFileName(tt.fileKey)
		assert.Equal(t, tt.wantName, name, tt.name)
		size, _ := w.GetFileSize(tt.fileKey)
		assert.Equal(t, tt.wantSize, size, tt.name)
	}
}
