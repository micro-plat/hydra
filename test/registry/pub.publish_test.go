package registry

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/jsons"
)

func getTestData(serviceAddr, clusterID string) string {
	input := map[string]interface{}{}
	input["addr"] = serviceAddr
	input["cluster_id"] = clusterID
	input["time"] = time.Now().Unix()

	buff, _ := jsons.Marshal(input)
	return string(buff)
}

func TestPublisher_PubRPCServiceNode(t *testing.T) {

	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080")
	confObj.Service.API.Add("/api1", "/api1", []string{"GET"})
	confObj.Service.API.Add("/api2", "/api1", []string{"GET"})
	confObj.Service.API.Add("/api3", "/api1", []string{"GET"})
	confObj.Service.API.Add("/api4", "/api1", []string{"GET"})
	confObj.Service.API.Add("/api5", "/api1", []string{"GET"})
	s := confObj.GetAPIConf() //初始化参数
	c := s.GetMainConf()      //获取配置
	lm := c.GetRegistry()

	p := pub.New(c)
	i := 1
	for _, service := range s.GetRouterConf().GetPath() {
		got, err := p.PubRPCServiceNode("192.168.5.115:9999", service, getTestData("192.168.5.115:9999", c.GetServerID()))
		assert.Equal(t, false, err != nil, "rpc服务发布")
		//验证pubs长度
		assert.Equal(t, i, len(got), "rpc服务发布")
		i++

		//验证节点是否发布成功
		for path, sdata := range got {
			ldata, v, err := lm.GetValue(path)
			assert.Equal(t, nil, err, "rpc服务节点发布验证")
			assert.NotEqual(t, v, int32(0), "rpc服务节点发布验证")
			assert.Equal(t, string(ldata), sdata, "rpc服务节点发布验证")
		}
	}
}

func TestPublisher_PubAPIServiceNode(t *testing.T) {
	tests := []struct {
		name       string
		serverName string
	}{
		{name: "api服务首次发布", serverName: "192.168.5.115:9999"},
		{name: "api服务再次发布", serverName: "192.168.5.115:8899"},
		{name: "api服务多次发布", serverName: "192.168.5.115:7799"},
	}
	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080")
	s := confObj.GetAPIConf() //初始化参数
	c := s.GetMainConf()      //获取配置
	lm := c.GetRegistry()

	p := pub.New(c)
	got := map[string]string{}
	var err error
	for _, tt := range tests {
		data := getTestData(tt.serverName, c.GetServerID())
		got, err = p.PubAPIServiceNode(tt.serverName, data)
		assert.Equal(t, false, err != nil, tt.name)
	}
	assert.Equal(t, len(tests), len(got), "api服务发布结果验证")

	//验证节点是否发布成功
	for path, sdata := range got {
		ldata, v, err := lm.GetValue(path)
		assert.Equal(t, nil, err, "api服务节点发布结果获取")
		assert.NotEqual(t, v, int32(0), "api服务节点结果版本号获取")
		assert.Equal(t, string(ldata), sdata, "api服务节点发布结果比对")
	}
}

func TestPublisher_PubServerNode(t *testing.T) {
	tests := []struct {
		name       string
		serverName string
	}{
		{name: "server服务首次发布", serverName: "192.168.5.115:9999"},
		{name: "server服务再次发布", serverName: "192.168.5.115:8899"},
		{name: "server服务多次发布", serverName: "192.168.5.115:7799"},
	}
	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080")
	s := confObj.GetAPIConf() //初始化参数
	c := s.GetMainConf()      //获取配置

	got := map[string]string{}
	var err error
	for _, tt := range tests {
		data := getTestData(tt.serverName, c.GetServerID())
		got, err = pub.New(c).PubServerNode(tt.serverName, data)
		assert.Equal(t, false, err != nil, tt.name)
	}
	assert.Equal(t, len(tests), len(got), "server服务发布验证")

	//验证节点是否发布成功
	lm := c.GetRegistry()
	for path, sdata := range got {
		ldata, v, err := lm.GetValue(path)
		assert.Equal(t, nil, err, "server服务节点发布结果获取")
		assert.NotEqual(t, v, int32(0), "server服务节点发布版本获取")
		assert.Equal(t, string(ldata), sdata, "server服务节点发布结果比对")
	}

}

func TestPublisher_PubDNSNode_WithDomain(t *testing.T) {
	tests := []struct {
		name       string
		serverName string
	}{
		{name: "dns发布", serverName: "192.168.5.115:9999"},
		{name: "dns再次发布", serverName: "192.168.5.115:8899"},
	}

	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080", api.WithDNS("192.168.0.101"))
	s := confObj.GetAPIConf() //初始化参数
	c := s.GetMainConf()      //获取配置

	got := map[string]string{}
	var err error
	for _, tt := range tests {
		got, err = pub.New(c).PubDNSNode(tt.serverName)
		assert.Equal(t, false, err != nil, tt.name)
	}
	assert.Equal(t, 1, len(got), "dns服务发布")

	//验证节点是否发布成功
	lm := c.GetRegistry()
	for path, sdata := range got {
		ldata, v, err := lm.GetValue(path)
		assert.Equal(t, nil, err, "dns服务节点发布验证")
		assert.NotEqual(t, v, int32(0), "dns服务节点发布验证")
		assert.Equal(t, string(ldata), sdata, "dns服务节点发布验证")
	}
}

func TestPublisher_PubDNSNode_NoDomain(t *testing.T) {
	//验证节点未设置Domain
	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080")
	s := confObj.GetAPIConf() //初始化参数
	c := s.GetMainConf()      //获取配置
	got, err := pub.New(c).PubDNSNode("192.168.5.115")
	assert.Equal(t, false, err != nil, "domain未设置不发布dns")
	assert.Equal(t, map[string]string{}, got, "domain未设置不发布dns")
}

func TestPublisher_Publish_API(t *testing.T) {

	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080", api.WithDNS("192.168.0.101"))
	confObj.RPC(":9377")
	apiconf := confObj.GetAPIConf() //初始化参数
	c := apiconf.GetMainConf()      //获取配置

	//发布api节点和dns节点
	err := pub.New(c).Publish("192.168.5.115:9091", "192.168.5.115:9091", c.GetServerID(), apiconf.GetRouterConf().GetPath()...)
	assert.Equal(t, false, err != nil, "发布api节点和dns节点")

	//验证节点发布结果
	lm := c.GetRegistry()
	_, _, err = lm.GetValue(c.GetServerPubPath(c.GetClusterName()))
	assert.Equal(t, nil, err, "servers服务节点发布验证")
	_, _, err = lm.GetValue(c.GetServicePubPath())
	assert.Equal(t, nil, err, "api服务节点发布验证")
	_, _, err = lm.GetValue(c.GetDNSPubPath("192.168.0.101"))
	assert.Equal(t, nil, err, "dns服务节点发布验证")

}

func TestPublisher_Publish_RPC(t *testing.T) {

	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080", api.WithDNS("192.168.0.101"))
	confObj.Service.API.Add("/api1", "/api1", []string{"GET"})
	confObj.Service.API.Add("/api2", "/api1", []string{"GET"})
	confObj.RPC(":9377")
	rpcconf := confObj.GetRPCConf() //初始化参数
	s := confObj.GetAPIConf()       //初始化参数
	c := rpcconf.GetMainConf()      //获取配置

	//发布rpc节点
	err := pub.New(c).Publish("192.168.5.115:9091", "192.168.5.115:9091", c.GetServerID(), s.GetRouterConf().GetPath()...)
	assert.Equal(t, false, err != nil, "发布rpc节点")

	lm := c.GetRegistry()
	_, _, err = lm.GetValue(c.GetRPCServicePubPath("api1"))
	assert.Equal(t, nil, err, "rpc服务节点发布验证")
	_, _, err = lm.GetValue(c.GetRPCServicePubPath("api2"))
	assert.Equal(t, nil, err, "rpc服务节点发布验证")
}

func TestPublisher_Update(t *testing.T) {

	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080", api.WithDNS("192.168.0.101"))
	confObj.RPC(":9377")
	apiconf := confObj.GetAPIConf() //初始化参数
	c := apiconf.GetMainConf()      //获取配置
	lm := c.GetRegistry()
	path := c.GetServerPubPath(c.GetClusterName())

	tests := []struct {
		name            string
		path            string
		k               string
		v               string
		wantUpdateErr   bool
		wantChildrenLen int
		wantValueErr    bool
	}{
		{name: "更新已经存在的API节点的不存在值", path: path, k: "key1", v: "value1", wantChildrenLen: 1},
		{name: "更新已经存在的API节点的存在值", path: path, k: "key1", v: "value1-1", wantChildrenLen: 1},
		{name: "更新不存在的API节点的值", path: path + "/ss", k: "key1", v: "value1", wantChildrenLen: 0},
	}

	p := pub.New(c)
	p.Publish("192.168.5.118:9091", "192.168.5.118:9091", c.GetServerID())

	for _, tt := range tests {
		//更新节点
		err := p.Update("192.168.5.118:9091", "192.168.5.118:9091", c.GetServerID(), tt.k, tt.v)
		assert.Equal(t, tt.wantUpdateErr, err != nil, tt.name)
		//获取更新结果
		paths, _, _ := lm.GetChildren(tt.path)
		assert.Equal(t, tt.wantChildrenLen, len(paths), tt.name)
		//验证
		for _, v := range paths {
			ldata, _, err := lm.GetValue(registry.Join(path, v))
			assert.Equal(t, tt.wantValueErr, err != nil, tt.name)
			value := map[string]string{}
			json.Unmarshal(ldata, &value)
			assert.Equal(t, tt.v, value[tt.k], tt.name)
		}
	}
}

//测试自动恢复节点
func TestNew(t *testing.T) {
	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf() //初始化参数
	c := apiconf.GetMainConf()      //获取配置

	//发布节点
	pub.New(c).Publish("192.168.5.118:9091", "192.168.5.118:9091", c.GetServerID())

	//删除节点
	lm := c.GetRegistry()
	pPath := c.GetServerPubPath(c.GetClusterName())
	paths, _, _ := lm.GetChildren(pPath)
	for _, v := range paths {
		path := registry.Join(pPath, v)
		err := lm.Delete(path)
		assert.Equal(t, nil, err, "NEW()-测试自动恢复节点")
		fmt.Printf("节点%s已删除", path)
	}

	//自动恢复节点
	paths, _, _ = lm.GetChildren(pPath)
	assert.Equal(t, 0, len(paths), "NEW()-测试自动恢复节点")

	time.Sleep(time.Second * 35)
	paths, _, _ = lm.GetChildren(pPath)
	assert.Equal(t, 1, len(paths), "NEW()-测试自动恢复节点")

}
