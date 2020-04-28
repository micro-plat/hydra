package ctx

import "github.com/qxnw/lib4go/logger"

//IRequest 请求信息
type IRequest interface {
	GetMethod() string
	GetService() string
	GetRequestPath() string
	GetCookie(string) (string, bool)
	GetHeader(string) string
	GetClientIP() string
}

//IResponse 响应信息
type IResponse interface {
	Header(string, string)
	GetStatusCode() int
	GetExt() string
	Abort(int)
	AbortWithError(int, error)
	File(string)
	GetResponseParam() map[string]interface{}
}

//IContext 用于中间件处理的上下文管理
type IContext interface {
	Request() IRequest
	Response() IResponse
	Log() logger.ILogger
	Next()
	Close()
}
