package gray

import (
	"testing"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/gray"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/pkgs/lua"
	"github.com/micro-plat/hydra/registry"
)

var modules = lua.Modules{
	"request": map[string]lua.LGFunction{
		"getClientIP": func(ls *lua.LState) int {
			ls.Push(lua.LString("abc"))
			return 1
		},
	},
}

func getGray() (*gray.Gray, error) {
	//设置服务器参数

	hydra.NewApp(hydra.WithServerTypes(http.API),
		hydra.WithUsage("apiserver"),
		hydra.WithRegistry("lm://."),
		hydra.WithPlatName("test"),
		hydra.WithClusterName("t"),
		hydra.WithSystemName("apiserver"),
	)
	if err := pkgs.Pub2Registry(true); err != nil {
		return nil, err
	}
	r, err := registry.NewRegistry(global.Def.RegistryAddr, global.Def.Log())
	if err != nil {
		return nil, err
	}
	mconf, err := server.NewMainConf(global.Def.PlatName, global.Def.SysName, global.ServerTypes[0], global.Def.ClusterName, r)
	if err != nil {
		return nil, err
	}

	return gray.GetConf(mconf), nil
}
func TestNoSetting(t *testing.T) {
	nconf, err := getGray()
	if err != nil {
		t.Error(err)
		return
	}
	if !nconf.Disable {
		t.Error("配置状态不正确")
		return
	}
}

func Test_RunScript(t *testing.T) {
	raw := `
	local ips = {}
	local upstream = ""
	
	
	function getUpStream()
		return upstream
	end
		
	function go2UpStream() 
		local req = require("request")
		local ip = req.getClientIP()
		if ips[ip] ~= nil then
			return true
		end
		return false
	end`

	hydra.Conf.API(":8081", api.WithTrace()).Gray(raw)
	nconf, err := getGray()
	if err != nil {
		t.Error(err)
		return
	}
	ok, err := nconf.NeedGo2UpStream(modules)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ok)
}

func BenchmarkCodeBlockMode(b *testing.B) {
	raw := `	
	function getUpStream()
		return upstream
	end
		
	function go2UpStream() 		
		return false
	end`

	hydra.Conf.API(":8081", api.WithTrace()).Gray(raw)
	g, err := getGray()
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := g.NeedGo2UpStream(modules)
		if err != nil {
			b.Error(err)
			return
		}
	}
}
