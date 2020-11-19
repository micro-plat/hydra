package servers

import (
	"testing"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/test/assert"
)

func TestRegister(t *testing.T) {

	fn := func(c app.IAPPConf) (servers.IResponsiveServer, error) {
		return http.NewResponsive(c)
	}
	//注册服务
	servers.Register("api-1", fn)

	global.Def.ServerTypes = []string{}
	//获取服务
	tps := servers.GetServerTypes()
	assert.Equal(t, tps, []string{}, "获取服务")

	global.Def.ServerTypes = []string{"api", "api-1"}

	//获取服务
	tps = servers.GetServerTypes()
	assert.Equal(t, tps, []string{"api", "api-1"}, "获取服务2")

	//再次注册服务
	assert.PanicError(t, "服务器[api-1]不能多次注册", func() { servers.Register("api-1", fn) }, "再次注册服务")
}
