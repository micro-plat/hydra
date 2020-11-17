package context

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/lib4go/logger"
)

const (
	JSONF  = "application/json; charset=%s"
	XMLF   = "application/xml; charset=%s"
	YAMLF  = "text/yaml; charset=%s"
	HTMLF  = "text/html; charset=%s"
	PLAINF = "text/plain; charset=%s"
)

//Handler 业务处理Handler
type Handler func(IContext) interface{}

//Handle 处理业务流程
func (h Handler) Handle(c IContext) interface{} {
	return h(c)
}

//IHandler 业务处理接口
type IHandler interface {
	//Handle 业务处理
	Handle(IContext) interface{}
}

//IGetter 参数获取
type IGetter interface {
	GetKeys() []string
	Get(name string) (string, bool)
	GetString(name string, def ...string) string
	GetMax(name string, o ...int) int
	GetMin(name string, o ...int) int
	GetInt(name string, def ...int) int
	GetInt64(name string, def ...int64) int64
	GetFloat32(name string, def ...float32) float32
	GetFloat64(name string, def ...float64) float64
	GetBool(name string, def ...bool) bool
	GetDatetime(name string, format ...string) (time.Time, error)
	IsEmpty(name string) bool
}

//IPath 请求参数
type IPath interface {
	//GetMethod 获取服务请求方法GET POST PUT DELETE 等
	GetMethod() string

	//GetRouter 获取当前请求对应的路由信息
	GetRouter() (*router.Router, error)

	//GetRequestPath 获取请求路径
	GetRequestPath() string

	//GetURL 获取请求的URL信息
	GetURL() string

	//GetCookie 获取请求Cookie
	GetCookie(string) (string, bool)

	//GetHeader 获取头信息
	GetHeader(string) string

	//GetHeaders 获取请求头
	GetHeaders() http.Header

	//GetCookies 获取cookie信息
	GetCookies() map[string]string

	//Limit 设置限流信息
	Limit(isLimit bool, fallback bool)
	//IsLimited 是否已限流
	IsLimited() bool

	//AllowFallback 是否允许降级
	AllowFallback() bool
}

//IVariable 参与变量
type IVariable interface {
	IGetter
}

//@fix 提供上传文件的处理 @hj
type IFile interface {
	SaveFile(fileKey, dst string) error
	GetFileSize(fileKey string) (int64, error)
	GetFileName(fileKey string) (string, error)
	GetFileBody(fileKey string) (io.ReadCloser, error)
}

//IRequest 请求信息
type IRequest interface {
	//Path 地址、头、cookie相关信息
	Path() IPath

	//Param 路由参数
	Param(string) string

	//Bind 将请求的参数绑定到对象
	Bind(obj interface{}) error

	//Check 检查指定的字段是否有值
	Check(field ...string) error

	//GetMap 将当前请求转换为map并返回
	GetMap() (map[string]interface{}, error)

	//GetRawBody 获取请求的body参数
	GetRawBody(encoding ...string) (string, error)

	//GetBody 获取请求的参数
	GetBody(encoding ...string) (string, error)

	//GetBodyMap 将body转换为map
	GetRawBodyMap(encoding ...string) (map[string]interface{}, error)

	// //GetTrace 获取请求的trace信息
	// GetTrace() string

	//GetPlayload 更改名称 @fix
	GetPlayload() string

	IGetter
	IFile
}

//IResponse 响应信息
type IResponse interface {

	//AddSpecial 添加特殊标记，用于在打印响应内容时知道当前请求进行了哪些特殊处理
	AddSpecial(t string)

	//GetSpecials 获取特殊标识字段串，多个标记用"|"分隔
	GetSpecials() string

	//Header 设置响应头
	Header(string, string)

	//GetRaw 获取未经处理的响应内容
	GetRaw() interface{}

	//StatusCode 设置状态码
	StatusCode(int)

	//ContentType 设置Content-Type响应头
	ContentType(v string)

	//NoNeedWrite 无需写入响应数据到缓存
	NoNeedWrite(status int)

	//WriteFinal 修改最终渲染内容(不会立即写入)
	WriteFinal(status int, content string, ctp string)

	//Write 向响应流中写入状态码与内容(不会立即写入)
	Write(s int, v interface{}) error

	//WriteAny 向响应流中写入内容,状态码根据内容进行判断(不会立即写入)
	WriteAny(v interface{}) error

	//File 向响应流中写入文件(立即写入)
	File(path string)

	//Abort 终止当前请求继续执行(立即写入)
	Abort(int, error)

	//Stop 停止当前服务执行(立即写入)
	Stop(int)

	//GetRawResponse 获取原始响应状态码与内容
	GetRawResponse() (int, interface{})

	//GetFinalResponse 获取最终渲染的响应状态码与内容
	GetFinalResponse() (int, string)

	//Flush 将当前内容写入响应流(立即写入)
	Flush()

	//GetHeaders 获取返回数据
	GetHeaders() map[string][]string
}

//IAuth 认证信息
type IAuth interface {
	//Request 获取或设置用户请求的认证信息
	Request(...interface{}) interface{}

	//Response 获取或设置系统响应的认证信息
	Response(...interface{}) interface{}

	//Bind 将请求的认证对象绑定为特定的结构体
	Bind(out interface{}) error
}

//IUser 用户相关信息
type IUser interface {

	//GetGID 获取当前处理的goroutine id
	GetGID() string

	//GetClientIP 获取客户端请求IP
	GetClientIP() string

	//GetRequestID 获取请求编号
	GetRequestID() string

	//Auth 认证信息
	Auth() IAuth
}

//IContext 用于中间件处理的上下文管理
type IContext interface {

	//Meta 元数据
	Meta() conf.IMeta

	//Request 请求信息
	Request() IRequest

	//Response 响应信息
	Response() IResponse

	//Context 控制超时的Context
	Context() context.Context

	//APPConf 服务器配置
	APPConf() app.IAPPConf

	//User 用户信息
	User() IUser

	//Log 日志组件
	Log() logger.ILogger

	//Close 关闭并释放资源
	Close()
}

//TFuncs 用于模板翻译的函数列表
type TFuncs map[string]interface{}

//Add 添加一个自定义的函数
func (f TFuncs) Add(name string, handle interface{}) {
	f[name] = handle
}
