package mocks

import (
	"context"
	"net/http"
	"strings"

	extcontext "github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/lib4go/logger"
)

type MiddleContext struct {
	MockNext       func()
	MockMeta       conf.IMeta
	MockUser       *MockUser
	MockTFuncs     extcontext.TFuncs
	MockRequest    extcontext.IRequest
	MockResponse   extcontext.IResponse
	HttpRequest    *http.Request
	HttpResponse   http.ResponseWriter
	MockServerConf server.IServerConf
}

func (ctx *MiddleContext) Next() {
	if ctx.MockNext != nil {
		ctx.MockNext()
	}
}

func (ctx *MiddleContext) Meta() conf.IMeta {
	return ctx.MockMeta
}

//Request 请求信息
func (ctx *MiddleContext) Request() extcontext.IRequest {
	return ctx.MockRequest
}

//Response 响应信息
func (ctx *MiddleContext) Response() extcontext.IResponse {
	return ctx.MockResponse
}

//Context 控制超时的Context
func (ctx *MiddleContext) Context() context.Context {
	return context.Background()
}

//ServerConf 服务器配置
func (ctx *MiddleContext) ServerConf() server.IServerConf {
	return ctx.MockServerConf
}

//TmplFuncs 模板函数列表
func (ctx *MiddleContext) TmplFuncs() extcontext.TFuncs {
	return ctx.MockTFuncs
}

//User 用户信息
func (ctx *MiddleContext) User() extcontext.IUser {
	return ctx.MockUser
}

//Log 日志组件
func (ctx *MiddleContext) Log() logger.ILogger {
	return logger.Nil()
}

//Close 关闭并释放资源
func (ctx *MiddleContext) Close() {}

func (ctx *MiddleContext) Trace(...interface{}) {}

//GetHttpReqResp GetHttpReqResp
func (ctx *MiddleContext) GetHttpReqResp() (*http.Request, http.ResponseWriter) {
	return ctx.HttpRequest, ctx.HttpResponse
}

type MockUser struct {
	MockClientIP  string
	MockRequestID string
	MockAuth      extcontext.IAuth
}

//GetClientIP 获取客户端请求IP
func (u *MockUser) GetClientIP() string {
	return u.MockClientIP
}

//GetRequestID 获取请求编号
func (u *MockUser) GetRequestID() string {
	return u.MockRequestID
}

//Auth 认证信息
func (u *MockUser) Auth() extcontext.IAuth {
	return u.MockAuth
}

type MockResponse struct {
	SpecialList     []string
	MockHeader      map[string]string
	MockRaw         interface{}
	MockStatus      int
	MockContent     string
	MockError       error
	MockContentType string
}

//AddSpecial 添加特殊标记，用于在打印响应内容时知道当前请求进行了哪些特殊处理
func (res *MockResponse) AddSpecial(t string) {
	res.SpecialList = append(res.SpecialList, t)
}

//GetSpecials 获取特殊标识字段串，多个标记用"|"分隔
func (res *MockResponse) GetSpecials() string {
	return strings.Join(res.SpecialList, "|")
}

//Header 设置响应头
func (res *MockResponse) Header(key string, val string) {
	res.MockHeader[key] = val
}

//GetRaw 获取未经处理的响应内容
func (res *MockResponse) GetRaw() interface{} {
	return res.MockRaw
}

//StatusCode 设置状态码
func (res *MockResponse) StatusCode(code int) {
	res.MockStatus = code
}

//ContentType 设置Content-Type响应头
func (res *MockResponse) ContentType(v string) {
	res.MockContentType = v
}

//NoNeedWrite 无需写入响应数据到缓存
func (res *MockResponse) NoNeedWrite(status int) {
	res.MockStatus = status
}

//Render 修改最终渲染内容
func (res *MockResponse) Render(status int, content string, ctp string) {
	res.MockStatus = status
	res.MockContent = content
	res.MockContentType = ctp
	return
}

//Write 向响应流中写入状态码与内容(不会立即写入)
func (res *MockResponse) Write(s int, v interface{}) error {
	return nil
}

//WriteAny 向响应流中写入内容,状态码根据内容进行判断(不会立即写入)
func (res *MockResponse) WriteAny(v interface{}) error {
	return nil
}

//File 向响应流中写入文件(立即写入)
func (res *MockResponse) File(path string) {

}

//Abort 终止当前请求继续执行
func (res *MockResponse) Abort(code int, err error) {
	res.MockStatus = code
	res.MockError = err
}

//Stop 停止当前服务执行
func (res *MockResponse) Stop(code int) {
	res.MockStatus = code
}

//GetRawResponse 获取原始响应状态码与内容
func (res *MockResponse) GetRawResponse() (int, interface{}) {
	return res.MockStatus, res.MockRaw
}

//GetFinalResponse 获取最终渲染的响应状态码与内容
func (res *MockResponse) GetFinalResponse() (int, string) {
	return res.MockStatus, res.MockContent
}

//Flush 将当前内容写入响应流
func (res *MockResponse) Flush() {
}

var _ http.ResponseWriter = &MockResponseWriter{}

type MockResponseWriter struct {
	ResponseHeader http.Header
	ContentBytes   []byte
	StatusCode     int
}

func (w *MockResponseWriter) Header() http.Header {
	return w.ResponseHeader
}
func (w *MockResponseWriter) Write(bytes []byte) (int, error) {
	w.ContentBytes = bytes
	return w.StatusCode, nil
}
func (w *MockResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}
