package conf

import (
	"encoding/json"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNewMainConf(t *testing.T) {
	systemName := "sys1"
	confM := mocks.NewConfBy("hydra1", "cluter1")
	confM.API(":8080")
	confM.Conf().Pub("hydra1", systemName, "cluter1", "lm://.", true)

	pubObj := server.NewServerPub("hydra1", systemName, global.API, "cluter1")
	mainPath := pubObj.GetServerPath()
	mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
	assert.Equal(t, true, err == nil, "获取manconf异常")

	assert.Equal(t, false, mainConf.IsTrace(), "IsTrace与配置不匹配")
	assert.Equal(t, confM.Registry, mainConf.GetRegistry(), "注册中心不匹配")
	assert.Equal(t, true, mainConf.IsStarted(), "IsStarted与配置不匹配")

	data, vsion, err := confM.Registry.GetValue(mainPath)
	assert.Equal(t, true, err == nil, "注册中心获取主配置信息异常")
	assert.Equal(t, vsion, mainConf.GetVersion(), "mainPath版本号不匹配")
	assert.Equal(t, data, mainConf.GetMainConf().GetRaw(), "主配置数据信息不匹配")

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
	pubObj = server.NewServerPub("hydra2", systemName, global.API, "cluter2")
	mainPath = pubObj.GetServerPath()
	mainConf, err = server.NewServerConf(confN.PlatName, systemName, global.API, confN.ClusterName, confN.Registry)
	assert.Equal(t, true, err == nil, "获取manconf异常1")
	assert.Equal(t, true, mainConf.IsTrace(), "IsTrace与配置不匹配1")
	assert.Equal(t, false, mainConf.IsStarted(), "IsStarted与配置不匹配1")

	data, vsion, err = confN.Registry.GetValue(mainPath)
	assert.Equal(t, true, err == nil, "注册中心获取主配置信息异常1")
	assert.Equal(t, vsion, mainConf.GetVersion(), "mainPath版本号不匹配1")
	assert.Equal(t, data, mainConf.GetMainConf().GetRaw(), "主配置数据信息不匹配1")

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

func TestMainConf_IsTrace(t *testing.T) {
	platName, systemName, clusterName := "hydra2", "sys2", "cluter2"
	tests := []struct {
		name    string
		opts    []api.Option
		want    bool
		wantErr bool
	}{
		{name: "1. Conf-MainConfIsTrace-没有设置trace", opts: []api.Option{}, want: false, wantErr: true},
		{name: "2. Conf-MainConfIsTrace-没有设置trace", opts: []api.Option{api.WithTrace()}, want: true, wantErr: true},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confM.API(":8080", tt.opts...)
		confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
		mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.want, mainConf.IsTrace(), tt.name+"Treace")
	}
}

func TestMainConf_IsStarted(t *testing.T) {
	platName, systemName, clusterName := "hydra3", "sys3", "cluter3"
	tests := []struct {
		name    string
		opts    []api.Option
		want    bool
		wantErr bool
	}{
		{name: "1. Conf-MainConfIsStarted-默认设置status", opts: []api.Option{}, want: true, wantErr: true},
		{name: "2. Conf-MainConfIsStarted-设置statsu==stop", opts: []api.Option{api.WithDisable()}, want: false, wantErr: true},
		{name: "3. Conf-MainConfIsStarted-设置statsu==start", opts: []api.Option{api.WithEnable()}, want: true, wantErr: true},
	}

	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confM.API(":8080", tt.opts...)
		confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
		mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.want, mainConf.IsStarted(), tt.name+"Started")
	}
}

func TestMainConf_GetRegistry(t *testing.T) {
	platName, systemName, clusterName := "hydra4", "sys4", "cluter4"
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{name: "1. Conf-MainConfGetRegistry-获取注册中心对象", address: "lm://.", wantErr: true},
	}
	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confM.API(":8080")
		confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
		mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		rgst, err := registry.GetRegistry(tt.address, global.Def.Log())
		assert.Equal(t, true, err == nil, tt.name+",err")
		assert.Equal(t, rgst, mainConf.GetRegistry(), tt.name+"Started")
	}
}

func TestMainConf_GetVersion(t *testing.T) {
	platName, systemName, clusterName := "hydra5", "sys5", "cluter5"
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "1. Conf-MainConfGetVersion-获取主节点的版本号", wantErr: true},
	}
	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confM.API(":8080")
		confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
		mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		pubObj := server.NewServerPub(platName, systemName, global.API, clusterName)
		mainPath := pubObj.GetServerPath()
		_, vsion, err := confM.Registry.GetValue(mainPath)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err1")
		assert.Equal(t, vsion, mainConf.GetVersion(), tt.name+",vsion")
	}
}

func TestMainConf_GetMainConf(t *testing.T) {
	platName, systemName, clusterName := "hydra5", "sys5", "cluter5"
	tests := []struct {
		name    string
		opts    []api.Option
		wantErr bool
	}{
		{name: "1. Conf-MainConfGetMainConf-获取空主节点的数据", opts: []api.Option{}, wantErr: true},
		{name: "2. Conf-MainConfGetMainConf-获取主节点的数据", opts: []api.Option{api.WithTrace(), api.WithDNS("192.168.30.11")}, wantErr: true},
	}
	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confM.API(":8080", tt.opts...)
		confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
		mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		pubObj := server.NewServerPub(platName, systemName, global.API, clusterName)
		mainPath := pubObj.GetServerPath()
		data, _, err := confM.Registry.GetValue(mainPath)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err1")
		assert.Equal(t, data, mainConf.GetMainConf().GetRaw(), tt.name+",data")
	}
}

func TestMainConf_GetMainObject(t *testing.T) {
	platName, systemName, clusterName := "hydra6", "sys6", "cluter6"
	tests := []struct {
		name    string
		opts    []api.Option
		wantErr bool
	}{
		{name: "1. Conf-MainConfGetMainObject-空主节点获取对象", opts: []api.Option{}, wantErr: true},
		{name: "2. Conf-MainConfGetMainObject-设置好的主节点获取对象", opts: []api.Option{api.WithDisable(), api.WithDNS("192.168.0.101"), api.WithTrace(), api.WithTimeout(20, 20), api.WithHeaderReadTimeout(15)}, wantErr: true},
	}
	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confM.API(":8080", tt.opts...)
		confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
		mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		pubObj := server.NewServerPub(platName, systemName, global.API, clusterName)
		mainPath := pubObj.GetServerPath()
		data, vsion, err := confM.Registry.GetValue(mainPath)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err1")
		apiObj1 := api.Server{}
		json.Unmarshal(data, &apiObj1)
		apiObj := api.Server{}
		version, err := mainConf.GetMainObject(&apiObj)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, vsion, version, tt.name+",vsion")
		assert.Equal(t, apiObj1, apiObj, tt.name+",data")
	}
}

func TestMainConf_GetSubConf(t *testing.T) {
	platName, systemName, clusterName := "hydra7", "sys7", "cluter7"
	tests := []struct {
		name    string
		isSet   bool
		secert  string
		opts    []apikey.Option
		wantErr bool
	}{
		{name: "1. Conf-MainConfGetSubConf-不设置子节点", isSet: false, secert: "", opts: []apikey.Option{}, wantErr: true},
		{name: "2. Conf-MainConfGetSubConf-设置好的子节点获取对象", isSet: true, secert: "123456", opts: []apikey.Option{apikey.WithDisable(), apikey.WithExcludes("xxx", "yyy"), apikey.WithSHA1Mode()}, wantErr: true},
	}
	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080", api.WithDisable())
		if tt.isSet {
			confN.APIKEY(tt.secert, tt.opts...)
		}
		confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
		mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		subObj, err := mainConf.GetSubConf("/auth/apikey")
		if tt.isSet {
			assert.Equal(t, true, err == nil, tt.name+",err")
			apikeykk := apikey.New(tt.secert, tt.opts...)
			assert.Equal(t, apikeykk.Secret, subObj.GetString("secret"), tt.name+"secret不匹配")
			assert.Equal(t, apikeykk.Mode, subObj.GetString("mode"), tt.name+"mod不匹配")
			assert.Equal(t, apikeykk.Disable, subObj.GetBool("disable"), tt.name+"置disable不匹配")
			assert.Equal(t, len(apikeykk.Excludes), len(subObj.GetArray("excludes")), tt.name+"excludes不匹配")
		} else {
			assert.Equal(t, true, err == conf.ErrNoSetting, tt.name+",err1")
		}
	}
}

func TestMainConf_GetCluster(t *testing.T) {
	platName, systemName, clusterName := "hydra8", "sys8", "cluter8"
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "1. Conf-MainConfGetCluster-获取主节点的集群对象", wantErr: true},
	}
	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confM.API(":8080", api.WithDisable())
		confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
		mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		_, err = mainConf.GetCluster(clusterName)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
	}
}

func TestMainConf_GetSubObject(t *testing.T) {
	platName, systemName, clusterName := "hydra9", "sys9", "cluter9"
	tests := []struct {
		name    string
		isSet   bool
		secert  string
		opts    []apikey.Option
		wantErr bool
	}{
		{name: "1. Conf-MainConfGetSubObject-不设置子节点", isSet: false, secert: "", opts: []apikey.Option{}, wantErr: true},
		{name: "2. Conf-MainConfGetSubObject-设置好的子节点获取对象", isSet: true, secert: "123456", opts: []apikey.Option{apikey.WithDisable(), apikey.WithExcludes("xxx", "yyy"), apikey.WithSHA1Mode()}, wantErr: true},
	}
	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080", api.WithDisable())
		if tt.isSet {
			confN.APIKEY(tt.secert, tt.opts...)
		}
		confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
		mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		apikeyObj := apikey.APIKeyAuth{}
		vsion, err := mainConf.GetSubObject("/auth/apikey", &apikeyObj)
		if tt.isSet {
			assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
			pubObj := server.NewServerPub(platName, systemName, global.API, clusterName)
			data, vsion1, err := confM.Registry.GetValue(pubObj.GetServerPath() + "/auth/apikey")
			assert.Equal(t, tt.wantErr, err == nil, tt.name+",err1")
			apikeyObj1 := apikey.APIKeyAuth{}
			json.Unmarshal(data, &apikeyObj1)
			assert.Equal(t, vsion1, vsion, tt.name+",vsion")
			assert.Equal(t, apikeyObj1, apikeyObj, tt.name+"secret不匹配")
		} else {
			assert.Equal(t, tt.wantErr, err == conf.ErrNoSetting, tt.name+",err1")
			assert.Equal(t, int32(0), vsion, tt.name+",vsion")
		}
	}
}

func TestMainConf_Has(t *testing.T) {
	platName, systemName, clusterName := "hydra10", "sys10", "cluter10"
	tests := []struct {
		name    string
		opts    []api.Option
		subOpts []apikey.Option
		wantErr bool
	}{
		{name: "1. Conf-MainConfHas-设置主节点", opts: []api.Option{api.WithDisable()}, subOpts: []apikey.Option{}, wantErr: true},
		{name: "2. Conf-MainConfHas-同时设置了apikey子节点", opts: []api.Option{api.WithDisable(), api.WithTrace(), api.WithTimeout(20, 20)}, subOpts: []apikey.Option{}, wantErr: true},
	}
	for _, tt := range tests {
		confM := mocks.NewConfBy(platName, clusterName)
		confN := confM.API(":8080", tt.opts...)
		if len(tt.subOpts) > 0 {
			confN.APIKEY("123456", tt.subOpts...)
		}
		confM.Conf().Pub(platName, systemName, clusterName, "lm://.", true)
		mainConf, err := server.NewServerConf(confM.PlatName, systemName, global.API, confM.ClusterName, confM.Registry)
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, true, mainConf.Has("router"), "router子配置是否存在判断失败1")
		if len(tt.subOpts) > 0 {
			assert.Equal(t, true, mainConf.Has("/auth/apikey"), "/auth/apikey子配置是否存在判断失败1")
		}
	}
}
