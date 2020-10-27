package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/global"

	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNewMainConf(t *testing.T) {
	systemName := "sys1"
	confM := mocks.NewConfBy("hydra1", "cluter1")
	confM.API(":8080")
	confM.Conf().Pub("hydra1", systemName, "cluter1", "lm://.", true)

	pubObj := server.NewPub("hydra1", systemName, global.API, "cluter1")
	mainPath := pubObj.GetMainPath()
	mainConf, err := server.NewMainConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
	assert.Equal(t, true, err == nil, "获取manconf异常")

	assert.Equal(t, false, mainConf.IsTrace(), "IsTrace与配置不匹配")
	assert.Equal(t, confM.Registry, mainConf.GetRegistry(), "注册中心不匹配")
	assert.Equal(t, true, mainConf.IsStarted(), "IsStarted与配置不匹配")

	data, vsion, err := confM.Registry.GetValue(mainPath)
	assert.Equal(t, true, err == nil, "注册中心获取主配置信息异常")
	assert.Equal(t, vsion, mainConf.GetVersion(), "mainPath版本号不匹配")
	assert.Equal(t, data, mainConf.GetRootConf().GetRaw(), "主配置数据信息不匹配")

	apiObj := api.Server{}
	vsion, err = mainConf.GetMainObject(&apiObj)
	assert.Equal(t, true, err == nil, "获取主配置对象反序列化异常")
	assert.Equal(t, vsion, mainConf.GetVersion(), "配置版本号不相同")
	assert.Equal(t, ":8080", apiObj.Address, "配置的端口号不相同")
	assert.Equal(t, "start", apiObj.Status, "配置的Status不相同")

	_, err = mainConf.GetSubConf("header")
	assert.Equal(t, conf.ErrNoSetting, err, "获取header子配置异常")

	_, err = mainConf.GetCluster("cluter1")
	assert.Equal(t, true, err == nil, "获取集群对象异常")

	_, err = mainConf.GetSubObject("header", &header.Headers{})
	assert.Equal(t, conf.ErrNoSetting, err, "获取header子配置对象异常")

	assert.Equal(t, true, mainConf.Has("router"), "router子配置是否存在判断失败")
	assert.Equal(t, false, mainConf.Has("/auth/apikey"), "/auth/apikey子配置是否存在判断失败")

	//从新设置完善的节点
	systemName = "sys2"
	confN := mocks.NewConfBy("hydra2", "cluter2")
	confN.API(":8080")
	subCOnf := confN.API(":8081", api.WithDisable(), api.WithDNS("192.168.0.101"), api.WithTrace(), api.WithTimeout(20, 20), api.WithHeaderReadTimeout(15))
	subCOnf.APIKEY("111111", apikey.WithSHA256Mode())
	confN.Conf().Pub("hydra2", systemName, "cluter2", "lm://.", true)
	pubObj = server.NewPub("hydra2", systemName, global.API, "cluter2")
	mainPath = pubObj.GetMainPath()
	mainConf, err = server.NewMainConf(confN.PlatName, systemName, global.API, confN.ClusterName, confN.Registry)
	assert.Equal(t, true, err == nil, "获取manconf异常1")
	assert.Equal(t, true, mainConf.IsTrace(), "IsTrace与配置不匹配1")
	assert.Equal(t, false, mainConf.IsStarted(), "IsStarted与配置不匹配1")

	data, vsion, err = confN.Registry.GetValue(mainPath)
	assert.Equal(t, true, err == nil, "注册中心获取主配置信息异常1")
	assert.Equal(t, vsion, mainConf.GetVersion(), "mainPath版本号不匹配1")
	assert.Equal(t, data, mainConf.GetRootConf().GetRaw(), "主配置数据信息不匹配1")

	apiObj = api.Server{}
	vsion, err = mainConf.GetMainObject(&apiObj)
	assert.Equal(t, true, err == nil, "获取主配置对象反序列化异常1")
	assert.Equal(t, vsion, mainConf.GetVersion(), "配置版本号不相同1")
	assert.Equal(t, ":8081", apiObj.Address, "配置的端口号不相同")
	assert.Equal(t, "stop", apiObj.Status, "配置的Status不相同")
	assert.Equal(t, 20, apiObj.RTimeout, "配置的RTimeout不相同")
	assert.Equal(t, 20, apiObj.WTimeout, "配置的WTimeout不相同")
	assert.Equal(t, true, apiObj.Trace, "配置的Trace不相同")
	assert.Equal(t, 15, apiObj.RHTimeout, "配置的RHTimeout不相同")
	assert.Equal(t, "192.168.0.101", apiObj.Domain, "配置的Domain不相同")

	subObj, err := mainConf.GetSubConf("/auth/apikey")
	assert.Equal(t, true, err == nil, "获取apikey子配置异常")
	assert.Equal(t, "111111", subObj.GetString("secret"), "apikey子配置secret不匹配1")
	assert.Equal(t, "SHA256", subObj.GetString("mode"), "apikey子配置mod不匹配1")

	vsion, err = mainConf.GetSubObject("/auth/apikey", &apikey.APIKeyAuth{})
	assert.Equal(t, vsion, subObj.GetVersion(), "获取apikey子配置版本号不匹配")

	assert.Equal(t, true, mainConf.Has("router"), "router子配置是否存在判断失败1")
	assert.Equal(t, true, mainConf.Has("/auth/apikey"), "/auth/apikey子配置是否存在判断失败1")
}
