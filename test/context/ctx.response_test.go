package context

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/logger"
)

func TestContentErr(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		content interface{}
		wantRs  int
		wantRc  string
		wantCt  string
		header  http.Header
	}{
		{name: "1.1.内容为err,状态码未设置", status: 0, content: errs.NewError(999, "错误"), wantRs: 999, wantRc: "错误", wantCt: "text/plain; charset=utf-8"},
		{name: "1.2.内容为err,状态码300", status: 300, content: errors.New("err"), wantRs: 400, wantRc: "err", wantCt: "text/plain; charset=utf-8"},
		{name: "1.3.内容为err,状态码在200", status: 200, content: errors.New("err"), wantRs: 400, wantRc: "err", wantCt: "text/plain; charset=utf-8"},
		{name: "1.4.内容为err,状态码500", status: 500, content: errors.New("err"), wantRs: 500, wantRc: "err", wantCt: "text/plain; charset=utf-8"},
		{name: "1.5.内容为err,状态码900", status: 900, content: errors.New("err"), wantRs: 900, wantRc: "err", wantCt: "text/plain; charset=utf-8"},
	}

	confObj := mocks.NewConfBy("context_response_test", "response") //构建对象
	confObj.API(":8080")                                            //初始化参数
	serverConf := confObj.GetAPIConf()                              //获取配置
	meta := conf.NewMeta()
	global.IsDebug = true
	for _, tt := range tests {
		contx := &mocks.TestContxt{HttpHeader: tt.header}

		log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(contx, "", meta).GetRequestID())

		//构建response对象
		c := ctx.NewResponse(contx, serverConf, log, meta)
		err := c.Write(tt.status, tt.content)
		assert.Equal(t, nil, err, tt.name)

		//测试reponse状态码和内容
		rs, rc, cp := c.GetFinalResponse()
		assert.Equal(t, tt.wantRs, rs, tt.name)
		assert.Equal(t, tt.wantRc, rc, tt.name)
		assert.Equal(t, tt.wantCt, cp, tt.name)

	}
}

func TestContentNil(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		content interface{}
		wantRs  int
		wantRc  string
		wantCt  string
		header  http.Header
	}{

		{name: "1.1.内容为空,未设置状态码,content-type未设置", status: 0, content: nil, wantRs: 200, wantRc: "", wantCt: "text/plain; charset=utf-8"},
		{name: "1.4.内容为空,未设置状态码,content-type为plain", status: 0, content: nil, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "", wantCt: "text/plain; charset=utf-8"},
		{name: "1.5.内容为空,未设置状态码,content-type为json", status: 0, content: nil, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: "", wantCt: "application/json; charset=utf-8"},
		{name: "1.6.内容为空,未设置状态码,content-type为xml", status: 0, content: nil, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "", wantCt: "application/xml; charset=utf-8"},
		{name: "1.6.内容为空,未设置状态码,content-type为yaml", status: 0, content: nil, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "", wantCt: "text/yaml; charset=utf-8"},

		{name: "2.1.内容为空,状态码为成功,content-type未设置", status: 200, content: nil, wantRs: 200, wantRc: "", wantCt: "text/plain; charset=utf-8"},
		{name: "2.2.内容为空,状态码为成功,content-type为plain", status: 200, content: nil, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "", wantCt: "text/plain; charset=utf-8"},
		{name: "2.3.内容为空,状态码为成功,content-type为json", status: 200, content: nil, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: "", wantCt: "application/json; charset=utf-8"},
		{name: "2.4.内容为空,状态码为成功,content-type为xml", status: 200, content: nil, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "", wantCt: "application/xml; charset=utf-8"},
		{name: "2.5.内容为空,状态码为成功,content-type为yaml", status: 200, content: nil, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "", wantCt: "text/yaml; charset=utf-8"},

		{name: "3.1.内容为空,状态码为600,content-type未设置", status: 600, content: nil, wantRs: 600, wantRc: "", wantCt: "text/plain; charset=utf-8"},
		{name: "3.2.内容为空,状态码为600,content-type为plain", status: 600, content: nil, wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "", wantCt: "text/plain; charset=utf-8"},
		{name: "3.3.内容为空,状态码为600,content-type为json", status: 600, content: nil, wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: "", wantCt: "application/json; charset=utf-8"},
		{name: "3.4.内容为空,状态码为600,content-type为xml", status: 600, content: nil, wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "", wantCt: "application/xml; charset=utf-8"},
		{name: "3.5.内容为空,状态码为600,content-type为yaml", status: 600, content: nil, wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "", wantCt: "text/yaml; charset=utf-8"},

		{name: "4.1.内容为空,状态码为400,content-type未设置", status: 400, content: nil, wantRs: 400, wantRc: "", wantCt: "text/plain; charset=utf-8"},
		{name: "4.2.内容为空,状态码为400,content-type为plain", status: 400, content: nil, wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "", wantCt: "text/plain; charset=utf-8"},
		{name: "4.3.内容为空,状态码为400,content-type为json", status: 400, content: nil, wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: "", wantCt: "application/json; charset=utf-8"},
		{name: "4.4.内容为空,状态码为400,content-type为xml", status: 400, content: nil, wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "", wantCt: "application/xml; charset=utf-8"},
		{name: "4.5.内容为空,状态码为400,content-type为yaml", status: 400, content: nil, wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "", wantCt: "text/yaml; charset=utf-8"},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	meta := conf.NewMeta()
	global.IsDebug = true
	for _, tt := range tests {
		contx := &mocks.TestContxt{HttpHeader: tt.header}

		log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(contx, "", meta).GetRequestID())

		//构建response对象
		c := ctx.NewResponse(contx, serverConf, log, meta)
		err := c.Write(tt.status, tt.content)
		assert.Equal(t, nil, err, tt.name)

		//测试reponse状态码和内容
		rs, rc, cp := c.GetFinalResponse()
		assert.Equal(t, tt.wantRs, rs, tt.name)
		assert.Equal(t, tt.wantRc, rc, tt.name)
		assert.Equal(t, tt.wantCt, cp, tt.name)

	}
}

func TestContentString(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		content interface{}
		wantRs  int
		wantRc  string
		wantCt  string
		header  http.Header
	}{

		{name: "1.1.内容字符串,未设置状态码,content-type未设置", status: 0, content: "success", wantRs: 200, wantRc: "success", wantCt: "text/plain; charset=utf-8"},
		{name: "1.2.内容字符串,未设置状态码,content-type为plain", status: 0, content: "success", wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "success", wantCt: "text/plain; charset=utf-8"},
		{name: "1.3.内容字符串,未设置状态码,content-type为json", status: 0, content: "success", wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: `success`, wantCt: "application/json; charset=utf-8"},
		{name: "1.4.内容字符串,未设置状态码,content-type为xml", status: 0, content: "success", wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "success", wantCt: "application/xml; charset=utf-8"},
		{name: "1.5.内容字符串,未设置状态码,content-type为yaml", status: 0, content: "success", wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "success", wantCt: "text/yaml; charset=utf-8"},

		{name: "2.1.内容字符串,状态码为成功,content-type未设置", status: 200, content: "success", wantRs: 200, wantRc: "success", wantCt: "text/plain; charset=utf-8"},
		{name: "2.2.内容字符串,状态码为成功,content-type为plain", status: 200, content: "success", wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "success", wantCt: "text/plain; charset=utf-8"},
		{name: "2.3.内容字符串,状态码为成功,content-type为json", status: 200, content: "success", wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: "success", wantCt: "application/json; charset=utf-8"},
		{name: "2.4.内容字符串,状态码为成功,content-type为xml", status: 200, content: "success", wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "success", wantCt: "application/xml; charset=utf-8"},
		{name: "2.5.内容字符串,状态码为成功,content-type为yaml", status: 200, content: "success", wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "success", wantCt: "text/yaml; charset=utf-8"},

		{name: "3.1.内容字符串,状态码为600,content-type未设置", status: 600, content: "success", wantRs: 600, wantRc: "success", wantCt: "text/plain; charset=utf-8"},
		{name: "3.2.内容字符串,状态码为600,content-type为plain", status: 600, content: "success", wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "success", wantCt: "text/plain; charset=utf-8"},
		{name: "3.3.内容字符串,状态码为600,content-type为json", status: 600, content: "success", wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: "success", wantCt: "application/json; charset=utf-8"},
		{name: "3.4.内容字符串,状态码为600,content-type为xml", status: 600, content: "success", wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "success", wantCt: "application/xml; charset=utf-8"},
		{name: "3.5.内容字符串,状态码为600,content-type为yaml", status: 600, content: "success", wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "success", wantCt: "text/yaml; charset=utf-8"},

		{name: "4.1.内容字符串,状态码为400,content-type未设置", status: 400, content: "success", wantRs: 400, wantRc: "success", wantCt: "text/plain; charset=utf-8"},
		{name: "4.2.内容字符串,状态码为400,content-type为plain", status: 400, content: "success", wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "success", wantCt: "text/plain; charset=utf-8"},
		{name: "4.3.内容字符串,状态码为400,content-type为json", status: 400, content: "success", wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: "success", wantCt: "application/json; charset=utf-8"},
		{name: "4.4.内容字符串,状态码为400,content-type为xml", status: 400, content: "success", wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "success", wantCt: "application/xml; charset=utf-8"},
		{name: "4.5.内容字符串,状态码为400,content-type为yaml", status: 400, content: "success", wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "success", wantCt: "text/yaml; charset=utf-8"},
	}

	confObj := mocks.NewConfBy("context_response_test1", "response1") //构建对象
	confObj.API(":8080")                                              //初始化参数
	serverConf := confObj.GetAPIConf()                                //获取配置
	meta := conf.NewMeta()
	global.IsDebug = true
	for _, tt := range tests {
		contx := &mocks.TestContxt{HttpHeader: tt.header}

		log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(contx, "", meta).GetRequestID())

		//构建response对象
		c := ctx.NewResponse(contx, serverConf, log, meta)
		err := c.Write(tt.status, tt.content)
		assert.Equal(t, nil, err, tt.name)

		//测试reponse状态码和内容
		rs, rc, cp := c.GetFinalResponse()
		assert.Equal(t, tt.wantRs, rs, tt.name)
		assert.Equal(t, tt.wantRc, rc, tt.name)
		assert.Equal(t, tt.wantCt, cp, tt.name)

	}
}

func TestContentMap(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		content interface{}
		wantRs  int
		wantRc  string
		wantCt  string
		header  http.Header
	}{

		{name: "1.1.内容Map,未设置状态码,content-type未设置", status: 0, content: map[string]interface{}{"id": 100}, wantRs: 200, wantRc: `{"id":100}`, wantCt: "application/json; charset=utf-8"},
		{name: "1.2.内容Map,未设置状态码,content-type为plain", status: 0, content: map[string]interface{}{"id": 100}, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "map[id:100]", wantCt: "text/plain; charset=utf-8"},
		{name: "1.3.内容Map,未设置状态码,content-type为json", status: 0, content: map[string]interface{}{"id": 100}, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: `{"id":100}`, wantCt: "application/json; charset=utf-8"},
		{name: "1.4.内容Map,未设置状态码,content-type为xml", status: 0, content: map[string]interface{}{"id": 100}, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "<id>100</id>", wantCt: "application/xml; charset=utf-8"},
		{name: "1.5.内容Map,未设置状态码,content-type为yaml", status: 0, content: map[string]interface{}{"id": 100}, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "id: 100\n", wantCt: "text/yaml; charset=utf-8"},

		{name: "2.1.内容Map,状态码为成功,content-type未设置", status: 200, content: map[string]interface{}{"id": 100}, wantRs: 200, wantRc: `{"id":100}`, wantCt: "application/json; charset=utf-8"},
		{name: "2.2.内容Map,状态码为成功,content-type为plain", status: 200, content: map[string]interface{}{"id": 100}, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "map[id:100]", wantCt: "text/plain; charset=utf-8"},
		{name: "2.3.内容Map,状态码为成功,content-type为json", status: 200, content: map[string]interface{}{"id": 100}, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: `{"id":100}`, wantCt: "application/json; charset=utf-8"},
		{name: "2.4.内容Map,状态码为成功,content-type为xml", status: 200, content: map[string]interface{}{"id": 100}, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "<id>100</id>", wantCt: "application/xml; charset=utf-8"},
		{name: "2.5.内容Map,状态码为成功,content-type为yaml", status: 200, content: map[string]interface{}{"id": 100}, wantRs: 200, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "id: 100\n", wantCt: "text/yaml; charset=utf-8"},

		{name: "3.1.内容Map,状态码为600,content-type未设置", status: 600, content: map[string]interface{}{"id": 100}, wantRs: 600, wantRc: `{"id":100}`, wantCt: "application/json; charset=utf-8"},
		{name: "3.2.内容Map,状态码为600,content-type为plain", status: 600, content: map[string]interface{}{"id": 100}, wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "map[id:100]", wantCt: "text/plain; charset=utf-8"},
		{name: "3.3.内容Map,状态码为600,content-type为json", status: 600, content: map[string]interface{}{"id": 100}, wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: `{"id":100}`, wantCt: "application/json; charset=utf-8"},
		{name: "3.4.内容Map,状态码为600,content-type为xml", status: 600, content: map[string]interface{}{"id": 100}, wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "<id>100</id>", wantCt: "application/xml; charset=utf-8"},
		{name: "3.5.内容Map,状态码为600,content-type为yaml", status: 600, content: map[string]interface{}{"id": 100}, wantRs: 600, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "id: 100\n", wantCt: "text/yaml; charset=utf-8"},

		{name: "4.1.内容Map,状态码为400,content-type未设置", status: 400, content: map[string]interface{}{"id": 100}, wantRs: 400, wantRc: `{"id":100}`, wantCt: "application/json; charset=utf-8"},
		{name: "4.2.内容Map,状态码为400,content-type为plain", status: 400, content: map[string]interface{}{"id": 100}, wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8PLAIN}}, wantRc: "map[id:100]", wantCt: "text/plain; charset=utf-8"},
		{name: "4.3.内容Map,状态码为400,content-type为json", status: 400, content: map[string]interface{}{"id": 100}, wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8JSON}}, wantRc: `{"id":100}`, wantCt: "application/json; charset=utf-8"},
		{name: "4.4.内容Map,状态码为400,content-type为xml", status: 400, content: map[string]interface{}{"id": 100}, wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8XML}}, wantRc: "<id>100</id>", wantCt: "application/xml; charset=utf-8"},
		{name: "4.5.内容Map,状态码为400,content-type为yaml", status: 400, content: map[string]interface{}{"id": 100}, wantRs: 400, header: http.Header{"Content-Type": []string{context.UTF8YAML}}, wantRc: "id: 100\n", wantCt: "text/yaml; charset=utf-8"},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	meta := conf.NewMeta()
	global.IsDebug = true
	for _, tt := range tests {
		contx := &mocks.TestContxt{HttpHeader: tt.header}

		log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(contx, "", meta).GetRequestID())

		//构建response对象
		c := ctx.NewResponse(contx, serverConf, log, meta)
		err := c.Write(tt.status, tt.content)
		assert.Equal(t, nil, err, tt.name)

		//测试reponse状态码和内容
		rs, rc, cp := c.GetFinalResponse()
		assert.Equal(t, tt.wantRs, rs, tt.name)
		assert.Equal(t, tt.wantRc, rc, tt.name)
		assert.Equal(t, tt.wantCt, cp, tt.name)

	}
}

func Test_response_Header(t *testing.T) {
	confObj := mocks.NewConfBy("context_response_test2", "response2") //构建对象
	confObj.API(":8080")                                              //初始化参数
	serverConf := confObj.GetAPIConf()                                //获取配置
	meta := conf.NewMeta()
	rc := &mocks.TestContxt{HttpHeader: http.Header{}}
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(rc, "", meta).GetRequestID())
	c := ctx.NewResponse(rc, serverConf, log, meta)

	//设置header
	c.Header("header1", "value1")
	assert.Equal(t, http.Header{"header1": []string{"value1"}}, rc.GetHeaders(), "设置header")

	//再次设置header
	c.Header("header1", "value1-1")
	assert.Equal(t, http.Header{"header1": []string{"value1-1"}}, rc.GetHeaders(), "再次设置header")

	//设置不同的header
	c.Header("header2", "value2")
	assert.Equal(t, http.Header{"header1": []string{"value1-1"}, "header2": []string{"value2"}}, rc.GetHeaders(), "再次设置header")
}

func Test_response_ContentType(t *testing.T) {
	confObj := mocks.NewConfBy("context_response_test3", "response3") //构建对象
	confObj.API(":8080")                                              //初始化参数
	serverConf := confObj.GetAPIConf()                                //获取配置
	meta := conf.NewMeta()
	rc := &mocks.TestContxt{HttpHeader: http.Header{}}
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(rc, "", meta).GetRequestID())
	c := ctx.NewResponse(rc, serverConf, log, meta)

	//设置content-type
	c.ContentType("application/json")
	assert.Equal(t, "application/json", rc.ContentType(), "设置content-type")

	//再次设置header
	c.ContentType("text/plain")
	assert.Equal(t, "text/plain", rc.ContentType(), "再次设置content-type")
}

func Test_response_Abort(t *testing.T) {
	confObj := mocks.mocks.NewConfBy("context_response_test4", "response4") //构建对象
	confObj.API(":8080")                                                    //初始化参数
	serverConf := confObj.GetAPIConf()                                      //获取配置
	meta := conf.NewMeta()
	context := &mocks.TestContxt{HttpHeader: http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}}
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(context, "", meta).GetRequestID())
	c := ctx.NewResponse(context, serverConf, log, meta)

	//测试Abort
	c.Abort(200, fmt.Errorf("终止"))
	rs, rc, cp := c.GetFinalResponse()
	assert.Equal(t, 400, rs, "验证状态码")
	assert.Equal(t, []byte(rc), context.Content, "验证返回内容")
	assert.Equal(t, rs, context.StatusCode, "验证上下文中的状态码")
	assert.Equal(t, context.HttpHeader["Content-Type"][0], cp, "验证上下文中的content-type")
	assert.Equal(t, true, context.WrittenStatus, "验证上下文中的写入状态")
	assert.Equal(t, true, context.Doen, "验证上下文中的abort状态")
}

func Test_response_Stop(t *testing.T) {
	confObj := mocks.NewConfBy("context_response_test5", "response5") //构建对象
	confObj.API(":8080")                                              //初始化参数
	serverConf := confObj.GetAPIConf()                                //获取配置
	meta := conf.NewMeta()
	context := &mocks.TestContxt{HttpHeader: http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}}
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(context, "", meta).GetRequestID())
	c := ctx.NewResponse(context, serverConf, log, meta)

	//测试Stop
	c.Abort(200)
	rs, rc, cp := c.GetFinalResponse()
	assert.Equal(t, 200, rs, "验证状态码")
	assert.Equal(t, []byte(rc), context.Content, "验证返回内容")
	assert.Equal(t, rs, context.StatusCode, "验证上下文中的状态码")
	assert.Equal(t, context.HttpHeader["Content-Type"][0], cp, "验证上下文中的content-type")
	assert.Equal(t, true, context.WrittenStatus, "验证上下文中的写入状态")
	assert.Equal(t, true, context.Doen, "验证上下文中的abort状态")
}

func Test_response_StatusCode(t *testing.T) {
	tests := []struct {
		name       string
		s          int
		wantStatus int
	}{
		{name: "设置状态码为200", s: 200, wantStatus: 200},
		{name: "设置状态码为300", s: 300, wantStatus: 300},
		{name: "设置状态码为400", s: 400, wantStatus: 400},
		{name: "设置状态码为500", s: 500, wantStatus: 500},
	}
	confObj := mocks.NewConfBy("context_response_test6", "response6") //构建对象
	confObj.API(":8080")                                              //初始化参数
	serverConf := confObj.GetAPIConf()                                //获取配置
	meta := conf.NewMeta()
	context := &mocks.TestContxt{HttpHeader: http.Header{}}
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(context, "", meta).GetRequestID())
	c := ctx.NewResponse(context, serverConf, log, meta)
	for _, tt := range tests {
		c.Write(tt.s)
		rs, _, _ := c.GetFinalResponse()
		assert.Equal(t, tt.wantStatus, rs, tt.name)
		// assert.Equal(t, context.Status(), rs, tt.name)
	}
}

func Test_response_File(t *testing.T) {
	confObj := mocks.NewConfBy("context_response_test7", "response7") //构建对象
	confObj.API(":8080")                                              //初始化参数
	serverConf := confObj.GetAPIConf()                                //获取配置
	meta := conf.NewMeta()
	context := &mocks.TestContxt{HttpHeader: http.Header{}}
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(context, "", meta).GetRequestID())
	c := ctx.NewResponse(context, serverConf, log, meta)

	//测试File
	c.File("file")
	assert.Equal(t, true, context.WrittenStatus, "验证上下文中的文件内容")
	assert.Equal(t, true, context.Doen, "验证上下文中的abort状态")
}

func Test_response_WriteFinal(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		content string
		ctp     string
		wantS   int
		wantC   string
		wantCP  string
	}{
		{name: "写入200状态码和json数据", status: 200, content: `{"a":"b"}`, ctp: "application/json", wantS: 200, wantC: `{"a":"b"}`},
		{name: "写入300状态码和空数据", status: 300, content: ``, ctp: "application/json", wantS: 300, wantC: ``},
		{name: "写入400状态码和错误数据", status: 400, content: `错误`, ctp: "application/json", wantS: 400, wantC: "错误"},
		{name: "写入空状态码和空数据", ctp: "application/json", wantS: 200, wantC: ""},
	}

	confObj := mocks.NewConfBy("context_response_test8", "response8") //构建对象
	confObj.API(":8080")                                              //初始化参数
	serverConf := confObj.GetAPIConf()                                //获取配置
	meta := conf.NewMeta()
	context := &mocks.TestContxt{HttpHeader: http.Header{}}
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(context, "", meta).GetRequestID())
	c := ctx.NewResponse(context, serverConf, log, meta)

	for _, tt := range tests {
		c.Write(tt.status, tt.content)
		rs, rc, cp := c.GetFinalResponse()
		assert.Equal(t, tt.wantS, rs, tt.name)
		assert.Equal(t, tt.wantC, rc, tt.name)
		assert.Equal(t, tt.wantCP, cp, tt.name)
	}
}

func Test_response_Redirect(t *testing.T) {

	confObj := mocks.NewConfBy("context_response_test9", "response9") //构建对象
	confObj.API(":8080")                                              //初始化参数
	serverConf := confObj.GetAPIConf()                                //获取配置
	meta := conf.NewMeta()
	context := &mocks.TestContxt{HttpHeader: http.Header{}}
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(context, "", meta).GetRequestID())
	c := ctx.NewResponse(context, serverConf, log, meta)

	c.Redirect(200, "url")
	assert.Equal(t, true, context.WrittenStatus, "验证上下文中的写入状态")
	assert.Equal(t, "url", context.Url, "验证上下文中的url")
	assert.Equal(t, 200, context.StatusCode, "验证上下文中的状态码")
}

func Test_response_Redirect_WithHttp(t *testing.T) {

	startServer()
	resp, err := http.Post("http://localhost:9091/response/redirect", "application/json", strings.NewReader(""))
	assert.Equal(t, false, err != nil, "重定向请求错误")
	defer resp.Body.Close()
	assert.Equal(t, "application/json; charset=UTF-8", resp.Header["Content-Type"][0], "重定向响应头错误")
	assert.Equal(t, "200 OK", resp.Status, "重定向响应状态错误")
	assert.Equal(t, 200, resp.StatusCode, "重定向响应码错误")
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, false, err != nil, "重定向响应体读取错误")
	assert.Equal(t, "success", string(body))

}
