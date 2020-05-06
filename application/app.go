package application

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry/conf/server"
)

var app IApplication = &application{}

type application struct {
}

func (a *application) GetCMD() string {
	return ""
}

//GetHandler 获取服务对应的处理函数
func (a *application) GetHandler(tp string, service string) context.IHandler {
	return nil
}

//Server 获取服务器配置信息
func (a *application) Server(tp string) server.IServerConf {
	if s, ok := server.Get(tp); ok {
		return s
	}
	panic(fmt.Sprintf("[%s]服务器未启动", tp))
}

//CurrentContext 获取当前请求上下文
func (a *application) CurrentContext() context.IContext {
	return nil
}
