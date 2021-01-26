package mock

import (
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/creator"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/http"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	"github.com/micro-plat/lib4go/types"
)

//NewContext 创建mock类型的Context包
func NewContext(content string, opts ...Option) context.IContext {

	//构建mock
	mk := newMock(content)
	for _, opt := range opts {
		opt(mk)
	}

	//初始化参数
	global.Def.PlatName = types.GetString(global.Def.PlatName, "mock_plat")
	global.Def.SysName = types.GetString(global.Def.SysName, "tserver")
	global.Def.ClusterName = types.GetString(global.Def.ClusterName, "test")
	global.Def.RegistryAddr = types.GetString(global.Def.RegistryAddr, "lm://.")
	global.Def.ServerTypes = []string{http.API}

	//发布配置
	err := creator.Conf.Pub(global.Current().GetPlatName(),
		global.Current().GetSysName(),
		global.Current().GetClusterName(),
		global.Def.RegistryAddr,
		true)
	if err != nil {
		panic(err)
	}

	//初始化缓存
	err = app.PullAndSave()
	if err != nil {
		panic(err)
	}

	//构建Context
	return ctx.NewCtx(mk, global.Def.ServerTypes[0])
}
