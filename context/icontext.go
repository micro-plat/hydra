package context

import (
	"time"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/registry/conf/server"
	"github.com/micro-plat/lib4go/logger"
)

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
	GetMethod() string
	GetService() string
	GetPath() string
	GetCookie(string) (string, bool)
	GetHeader(string) string
	GetHeaders() map[string][]string
	GetCookies() map[string]string
}

//IVariable 参与变量
type IVariable interface {
	IGetter
}

//IRequest 请求信息
type IRequest interface {
	Path() IPath
	Bind(obj interface{}) error
	Check(field ...string) error
	GetData() (map[string]interface{}, error)
	GetBody(encoding ...string) (string, error)
	GetBodyMap(encoding ...string) (map[string]interface{}, error)
	GetTrace() string
	IGetter
}

//IResponse 响应信息
type IResponse interface {
	AddSpecial(t string)
	GetSpecials() string
	SetHeader(string, string)
	GetStatusCode() int
	SetStatusCode(int)
	Write(s int, v string) error
	WriteAny(v interface{}) error
	Written() bool
	File(path string)
	Abort(int)
	AbortWithError(int, error)
	GetTrace() string
}

//IUser 用户相关信息
type IUser interface {
	GetClientIP() string
	GetRequestID() string
	SaveJwt(v interface{})
	GetJwt() interface{}
	BindJwt(out interface{})
}

//IApplication 应用信息

//IContext 用于中间件处理的上下文管理
type IContext interface {
	Request() IRequest
	Response() IResponse
	Server() server.IServerConf
	Component() components.IComponent
	User() IUser
	Log() logger.ILogger
	Close()
}
