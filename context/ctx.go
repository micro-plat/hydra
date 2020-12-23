package context

import (
	"context"
	"io"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
)

const (
	UserName = "UserName"

	XRequestID = "X-Request-Id"

	JSONF  = "application/json; charset=%s"
	XMLF   = "application/xml; charset=%s"
	YAMLF  = "text/yaml; charset=%s"
	HTMLF  = "text/html; charset=%s"
	PLAINF = "text/plain; charset=%s"

	UTF8JSON  = "application/json; charset=utf-8"
	UTF8XML   = "application/xml; charset=utf-8"
	UTF8YAML  = "text/yaml; charset=utf-8"
	UTF8HTML  = "text/html; charset=utf-8"
	UTF8PLAIN = "text/plain; charset=utf-8"
)

var EmptyReponseResult = &EmptyResult{}

type EmptyResult struct{}

//Handler 业务处理Handler
type Handler func(IContext) interface{}

//Handle 处理业务流程
func (h Handler) Handle(c IContext) interface{} {
	return h(c)
}

//VoidHandler 无返回值的Handler处理
type VoidHandler func(IContext)

//Handle 处理业务流程
func (h VoidHandler) Handle(c IContext) interface{} {
	h(c)
	return EmptyReponseResult
}

//IHandler 业务处理接口
type IHandler interface {
	//Handle 业务处理
	Handle(IContext) interface{}
}

//IGetter 参数获取
type IGetter interface {

	//GetKeys 获取所有参数名称
	GetKeys() []string

	//Get 获取请求参数
	Get(name string) (string, bool)

	//GetString 获取字符串
	GetString(name string, def ...string) string

	//GetInt 获取类型为int的值
	GetInt(name string, def ...int) int

	//GetInt32 获取类型为int32的值
	GetInt32(name string, def ...int32) int32

	//GetInt64 获取类型为int64的值
	GetInt64(name string, def ...int64) int64

	//GetFloat32 获取类型为float32的值
	GetFloat32(name string, def ...float32) float32

	//GetFloat64 获取类型为float64的值
	GetFloat64(name string, def ...float64) float64

	//GetDecimal 获取类型为Decimal的值
	GetDecimal(name string, def ...types.Decimal) types.Decimal

	//GetDatetime 获取日期类型的值
	GetDatetime(name string, format ...string) (time.Time, error)
	IsEmpty(name string) bool
}

//IPath 请求参数
type IPath interface {
	//GetMethod 获取服务请求方法GET POST PUT DELETE 等
	GetMethod() string

	//GetService 获取服务名称
	GetService() string

	//Param 路由参数
	Params() types.XMap

	//GetRouter 获取当前请求对应的路由信息
	GetRouter() (*router.Router, error)

	//GetRequestPath 获取请求路径
	GetRequestPath() string

	//GetURL 获取请求的URL信息
	GetURL() string

	//Limit 设置限流信息
	Limit(isLimit bool, fallback bool)

	//IsLimited 是否已限流
	IsLimited() bool

	//AllowFallback 是否允许降级
	AllowFallback() bool

	GetEncoding() string
}

//IVariable 参与变量
type IVariable interface {
	IGetter
}

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

	//Bind 将请求的参数绑定到对象
	Bind(obj interface{}) error

	//Check 检查指定的字段是否有值
	Check(field ...string) error

	//GetMap 将当前请求转换为map并返回
	GetMap() types.XMap

	//GetFullRaw 获取请求的body参数
	GetFullRaw() (body []byte, query string, err error)

	//GetError 请求解析过程中发生的异常
	GetError() error

	//GetBody 获取请求的参数
	GetBody() (body []byte, err error)

	//GetPlayload
	GetPlayload() string

	//Headers 获取请求头
	Headers() types.XMap

	//Cookies 获取cookie信息
	Cookies() types.XMap

	types.IXMap

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

	//ContentType 设置Content-Type响应头,自动增加charset或编码值,如指定值为:application/json;或application/json; charset=%s
	//最终输出结果为 application/json; charset=utf-8 或 application/json; charset=gbk
	//具体的charset值与服务配置和请求的Content-Type中指定的charset有关
	ContentType(v string, xmlRoot ...string)

	//NoNeedWrite 无需写入响应数据到缓存
	NoNeedWrite(status int)

	//JSON json输出响应内容
	JSON(code int, data interface{}) interface{}

	//XML xml输出响应内容
	XML(code int, data interface{}, header string, rootNode ...string) interface{}

	//以text/html输出响应内容
	HTML(code int, data string) interface{}

	//YAML yaml输出响应内容
	YAML(code int, data interface{}) interface{}

	//以text/plain格式输出响应内容
	Plain(code int, data string) interface{}

	//Data 使用已设置的Content-Type输出内容，未设置时自动根据内容识别输出格式，内容无法识别时(map,struct)使用application/json
	//格式输出内容
	Data(code int, contentType string, data interface{}) interface{}

	//WriteAny 使用已设置的Content-Type输出内容，未设置时自动根据内容识别输出格式，内容无法识别时(map,struct)使用application/json
	//格式输出内容
	WriteAny(v interface{}) error

	//Write 使用已设置的Content-Type输出内容，未设置时自动根据内容识别输出格式，内容无法识别时(map,struct)使用application/json
	//格式输出内容
	Write(s int, v ...interface{}) error

	//File 向响应流中写入文件(立即写入)
	File(path string)

	//Abort 停止当前服务执行(立即写入)
	Abort(int, ...interface{})

	//GetRawResponse 获取原始响应状态码与内容
	GetRawResponse() (int, interface{}, string)

	//GetFinalResponse 获取最终渲染的响应状态码与内容
	GetFinalResponse() (statusCode int, content string, contentType string)

	//Flush 将当前内容写入响应流(立即写入)
	Flush()

	//GetHeaders 获取返回数据
	GetHeaders() types.XMap
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

	//GetUserName 获取用户名
	GetUserName() string

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

	//链路跟踪器
	Tracer() ITracer

	//Invoke 调用本地服务
	Invoke(service string) interface{}

	//Close 关闭并释放资源
	Close()
}

//ITracer 链路跟踪器
type ITracer interface {
	ITraceSpan

	//Root 根结点Tracer
	Root() ITraceSpan
}

//ITraceSpan 跟踪处理器
type ITraceSpan interface {
	IEnd

	//Available 是否可用
	Available() bool

	//Start 开发跟踪
	Start() IEnd

	//NewSpan 新的时间片
	NewSpan(opertor string) ITraceSpan
}

//IEnd 关闭
type IEnd interface {
	//End 结束跟踪
	End()
}
