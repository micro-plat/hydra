package servers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
)

func Test_ginCtx_Get(t *testing.T) {
	//构建上下文
	c := &gin.Context{}
	router := gin.New()
	path := "/getbody"
	router.POST(path, func(c *gin.Context) {
		return
	})

	//构建请求 方法要与注册方法一致
	body := `{"a":"1"}`
	r, err := http.NewRequest("POST", "http://localhost:9091/getbody", strings.NewReader(body))
	assert.Equal(t, nil, err, "构建请求")

	//设置content-type
	ctp := "application/json"
	r.Header.Set("Content-Type", ctp)

	//添加cookie
	cookie := http.Cookie{Name: "session", Value: "value"}
	r.AddCookie(&cookie)

	//替换gin上下文的请求
	c.Request = r
	router.HandleContext(c)

	g := middleware.NewGinCtx(c)

	//GetRouterPath
	gotFullPath := g.GetRouterPath()
	assert.Equal(t, path, gotFullPath, "GetRouterPath")

	//GetMethod
	gotMethod := g.GetMethod()
	assert.Equal(t, "POST", gotMethod, "GetMethod")

	//GetURL
	goturl := g.GetURL()
	assert.Equal(t, &url.URL{Scheme: "http", Host: "localhost:9091", Path: path}, goturl, "GetURL")

	//GetHeaders
	gotHeader := g.GetHeaders()
	assert.Equal(t, http.Header{"Content-Type": []string{ctp}, "Cookie": []string{"session=value"}}, gotHeader, "GetHeaders")

	//GetCookies
	gotCookies := g.GetCookies()
	assert.Equal(t, []*http.Cookie{&cookie}, gotCookies, "GetCookies")

	//GetBody
	gotBody, _ := ioutil.ReadAll(g.GetBody())
	assert.Equal(t, body, string(gotBody), "GetBody")
}

func Test_ginCtx_Get_WithForm(t *testing.T) {
	//构建上下文
	c := &gin.Context{}
	router := gin.New()
	path := "/getbody"
	router.POST(path, func(c *gin.Context) {
		return
	})

	//构建请求 方法要与注册方法一致
	body := `a=1&b=2&c=3`
	r, err := http.NewRequest("POST", "http://localhost:9091/getbody", strings.NewReader(body))
	assert.Equal(t, nil, err, "构建请求")

	//设置content-type
	ctp := "application/x-www-form-urlencoded"
	r.Header.Set("Content-Type", ctp)

	//添加cookie
	cookie := http.Cookie{Name: "session", Value: "value"}
	r.AddCookie(&cookie)

	//替换gin上下文的请求
	c.Request = r
	router.HandleContext(c)

	g := middleware.NewGinCtx(c)

	//GetRouterPath
	gotFullPath := g.GetRouterPath()
	assert.Equal(t, path, gotFullPath, "GetRouterPath")

	//GetMethod
	gotMethod := g.GetMethod()
	assert.Equal(t, "POST", gotMethod, "GetMethod")

	//GetURL
	goturl := g.GetURL()
	assert.Equal(t, &url.URL{Scheme: "http", Host: "localhost:9091", Path: path}, goturl, "GetURL")

	//GetHeaders
	gotHeader := g.GetHeaders()
	assert.Equal(t, http.Header{"Content-Type": []string{ctp}, "Cookie": []string{"session=value"}}, gotHeader, "GetHeaders")

	//GetCookies
	gotCookies := g.GetCookies()
	assert.Equal(t, []*http.Cookie{&cookie}, gotCookies, "GetCookies")

	//GetBody body里面的数据被处理到http.Request对象 getbody为空
	gotBody, _ := ioutil.ReadAll(g.GetBody())
	assert.Equal(t, "", string(gotBody), "GetBody")

	//GetForm
	gotForm := g.GetPostForm()
	assert.Equal(t, url.Values{"a": []string{"1"}, "b": []string{"2"}, "c": []string{"3"}}, gotForm, "GetForm")
}

//gin.Context无法设置engine
// func Test_ginCtx_Get_WithMultiForm(t *testing.T) {
// 	//构建上下文
// 	c := &gin.Context{}
// 	router := gin.New()
// 	path := "/upload"
// 	router.POST(path, func(c *gin.Context) {
// 		return
// 	})

// 	//构建请求 方法要与注册方法一致
// 	fileName := "pkgs.middleware.middle.test.txt"
// 	file, _ := os.Open(fileName)
// 	defer file.Close()
// 	body := &bytes.Buffer{}
// 	// 文件写入 body
// 	writer := multipart.NewWriter(body)
// 	part, _ := writer.CreateFormFile("upload", filepath.Base(fileName))
// 	io.Copy(part, file)
// 	writer.Close()

// 	r, err := http.NewRequest("POST", "http://localhost:9091/upload", body)
// 	assert.Equal(t, nil, err, "构建请求")

// 	//设置content-type
// 	ctp := "multipart/form-data"
// 	r.Header.Set("Content-Type", ctp)

// 	//替换gin上下文的请求
// 	c.Request = r
// 	router.HandleContext(c)

// 	g := middleware.NewGinCtx(c)

// 	//GetForm
// 	gotFileName, rbody, gotSize, gotErr := g.GetFile("upload")
// 	assert.Equal(t, fileName, gotFileName, "GetForm")
// 	assert.Equal(t, nil, gotErr, "GetForm")
// 	assert.Equal(t, 10, gotSize, "GetForm")

// 	gotBody, _ := ioutil.ReadAll(rbody)
// 	assert.Equal(t, "", string(gotBody), "GetForm")
// }
