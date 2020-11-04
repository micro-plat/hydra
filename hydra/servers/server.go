package servers

import (
	"fmt"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
)

//IServerCreatorHandler 服务器创建器
type IServerCreatorHandler func(app.IAPPConf) (IResponsiveServer, error)

//Create 创建服务
func (r IServerCreatorHandler) Create(c app.IAPPConf) (IResponsiveServer, error) {
	return r(c)
}

//IServerCreator 服务器构建嚣
type IServerCreator interface {
	Create(app.IAPPConf) (IResponsiveServer, error)
}

//IResponsiveServer 响应式服务器
type IResponsiveServer interface {
	Start() error
	Notify(app.IAPPConf) (bool, error)
	Shutdown()
}

var creators = make(map[string]IServerCreator)

//Register 注册服务器生成器
func Register(tp string, creator IServerCreatorHandler) {
	if _, ok := creators[tp]; ok {
		panic(fmt.Sprintf("服务器[%s]不能多次注册", tp))
	}
	global.ServerTypes = append(global.ServerTypes, tp)
	creators[tp] = creator
}

//GetServerTypes 获取支付的服务器类型
func GetServerTypes() []string {
	tps := make([]string, 0, len(creators))
	for _, s := range global.Def.ServerTypes {
		if _, ok := creators[s]; ok {
			tps = append(tps, s)
		}
	}
	return tps
}
