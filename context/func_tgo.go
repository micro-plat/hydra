package context

import (
	"strings"

	"github.com/micro-plat/hydra/context/internal"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/tgo"
)

//GetTGOModules 获取tgo的模块配置信息
func GetTGOModules() []*tgo.Module {
	request := tgo.NewModule("request").
		Add("getString", tgo.FuncASRS(func(input string) string { ctx := Current(); return ctx.Request().GetString(input) })).
		Add("getPath", tgo.FuncARS(func() string { ctx := Current(); return ctx.Request().Path().GetRequestPath() })).
		Add("getHeader", tgo.FuncASRS(func(input string) string { ctx := Current(); return ctx.Request().Headers().GetString(input) })).
		Add("getCookie", tgo.FuncASRS(func(input string) string { ctx := Current(); return ctx.Request().Cookies().GetString(input) })).
		Add("getClientIP", tgo.FuncARS(func() string { ctx := Current(); return ctx.User().GetClientIP() })).
		Add("getTraceID", tgo.FuncARS(func() string { ctx := Current(); return ctx.User().GetTraceID() }))

	response := tgo.NewModule("response").
		Add("getStatus", tgo.FuncARI(func() int { ctx := Current(); s, _, _ := ctx.Response().GetFinalResponse(); return s })).
		Add("getContent", tgo.FuncARS(func() string { ctx := Current(); _, s, _ := ctx.Response().GetFinalResponse(); return s })).
		Add("getRaw", internal.IASANY(func() interface{} { ctx := Current(); return ctx.Response().GetRaw() }))

	app := tgo.NewModule("app").
		Add("getServerID", tgo.FuncARS(func() string { ctx := Current(); return ctx.APPConf().GetServerConf().GetServerID() })).
		Add("getPlatName", tgo.FuncARS(func() string { ctx := Current(); return ctx.APPConf().GetServerConf().GetPlatName() })).
		Add("getSysName", tgo.FuncARS(func() string { ctx := Current(); return ctx.APPConf().GetServerConf().GetSysName() })).
		Add("getServerType", tgo.FuncARS(func() string { ctx := Current(); return ctx.APPConf().GetServerConf().GetServerType() })).
		Add("getClusterNameBy", tgo.FuncASRS(getClusterNameBy)).
		Add("getCurrentClusterName", tgo.FuncARS(func() string { ctx := Current(); return ctx.APPConf().GetServerConf().GetClusterName() })).
		Add("getAllClusterNames", tgo.FuncARSs(func() []string { ctx := Current(); return ctx.APPConf().GetServerConf().GetClusterNames() })).
		Add("getServerName", tgo.FuncARS(func() string { ctx := Current(); return ctx.APPConf().GetServerConf().GetServerName() })).
		Add("getServerPath", tgo.FuncARS(func() string { ctx := Current(); return ctx.APPConf().GetServerConf().GetServerPath() }))

	types := tgo.NewModule("types").
		Add("getStringByIndex", internal.GetStringByIndex).
		Add("getIntByIndex", internal.GetIntByIndex).
		Add("getFloatByIndex", internal.GetFloatByIndex).
		Add("exclude", internal.Exclude).
		Add("translate", internal.Translate)

	users := tgo.NewModule("user").
		Add("getUserInfo", internal.IASANY(getUserInfo))

	return []*tgo.Module{request, response, app, types, users}
}

func init() {
	global.AddTGOModules(GetTGOModules()...)
}
func getClusterNameBy(n string) string {
	ctx := Current()
	names := ctx.APPConf().GetServerConf().GetClusterNames()
	for _, nm := range names {
		if strings.Contains(nm, n) {
			return nm
		}
	}
	return ""
}

func getUserInfo() interface{} {
	ctx := Current()
	mp := make(map[string]interface{})
	ctx.User().Auth().Bind(&mp)
	return mp
}
