package servers

import (
	"fmt"

	"github.com/micro-plat/hydra/registry/conf/server"
)

type IServerCreatorHandler func(server.IServerConf) (IResponsiveServer, error)

//Create 创建服务
func (r IServerCreatorHandler) Create(c server.IServerConf) (IResponsiveServer, error) {
	return r(c)
}

//IServerCreator 服务器构建嚣
type IServerCreator interface {
	Create(server.IServerConf) (IResponsiveServer, error)
}

//IResponsiveServer 响应式服务器
type IResponsiveServer interface {
	Start() error
	Notify(server.IServerConf) error
	Shutdown() error
}

var creators = make(map[string]IServerCreator)

//Register 注册服务器生成器
func Register(tp string, creator IServerCreator) {
	if _, ok := creators[tp]; ok {
		panic(fmt.Sprintf("服务器[%s]不能多次注册", tp))
	}
	creators[tp] = creator
}
