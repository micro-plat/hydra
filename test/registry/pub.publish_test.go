package registry

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/api"
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
		{name: "api服务发布1", serverName: "192.168.5.115:9999"},
		{name: "api服务发布2", serverName: "192.168.5.115:8899"},
		{name: "api服务发布3", serverName: "192.168.5.115:7799"},
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
	assert.Equal(t, len(tests), len(got), "api服务发布")

	//验证节点是否发布成功
	for path, sdata := range got {
		ldata, v, err := lm.GetValue(path)
		assert.Equal(t, nil, err, "api服务节点发布验证")
		assert.NotEqual(t, v, int32(0), "api服务节点发布验证")
		assert.Equal(t, string(ldata), sdata, "api服务节点发布验证")
	}
}

func TestPublisher_PubServerNode(t *testing.T) {
	tests := []struct {
		name       string
		serverName string
	}{
		{name: "server服务发布1", serverName: "192.168.5.115:9999"},
		{name: "server服务发布2", serverName: "192.168.5.115:8899"},
		{name: "server服务发布3", serverName: "192.168.5.115:7799"},
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
		got, err = p.PubServerNode(tt.serverName, data)
		assert.Equal(t, false, err != nil, tt.name)
	}
	assert.Equal(t, len(tests), len(got), "server服务发布")

	//验证节点是否发布成功
	for path, sdata := range got {
		ldata, v, err := lm.GetValue(path)
		assert.Equal(t, nil, err, "server服务节点发布验证")
		assert.NotEqual(t, v, int32(0), "server服务节点发布验证")
		assert.Equal(t, string(ldata), sdata, "server服务节点发布验证")
	}

}

func TestPublisher_PubDNSNode(t *testing.T) {
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
	lm := c.GetRegistry()

	p := pub.New(c)
	got := map[string]string{}
	var err error
	for _, tt := range tests {
		got, err = p.PubDNSNode(tt.serverName)
		assert.Equal(t, false, err != nil, tt.name)
	}
	assert.Equal(t, 1, len(got), "dns服务发布")

	//验证节点是否发布成功
	for path, sdata := range got {
		ldata, v, err := lm.GetValue(path)
		assert.Equal(t, nil, err, "dns服务节点发布验证")
		assert.NotEqual(t, v, int32(0), "dns服务节点发布验证")
		assert.Equal(t, string(ldata), sdata, "dns服务节点发布验证")
	}

	//验证节点未设置Domain
	confObj2 := mocks.NewConf() //构建对象
	confObj2.API(":8080")
	s2 := confObj2.GetAPIConf() //初始化参数
	c2 := s2.GetMainConf()      //获取配置
	p2 := pub.New(c2)
	got2, err2 := p2.PubDNSNode("192.168.5.115")
	assert.Equal(t, false, err2 != nil, "domain未设置不发布dns")
	assert.Equal(t, map[string]string{}, got2, "domain未设置不发布dns")
}

func TestPublisher_Publish(t *testing.T) {

	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080", api.WithDNS("192.168.0.101"))
	confObj.RPC(":9377")

	serverName := "192.168.5.115:9091"
	serviceAddr := "192.168.5.115:9091"

	//发布api节点和dns节点
	apiconf := confObj.GetAPIConf() //初始化参数
	c := apiconf.GetMainConf()      //获取配置
	lm := c.GetRegistry()
	p := pub.New(c)
	err := p.Publish(serverName, serviceAddr, c.GetServerID(), apiconf.GetRouterConf().GetPath()...)
	assert.Equal(t, false, err != nil, "发布api节点和dns节点")
	_, _, err = lm.GetValue("/hydra/apiserver/api/test/servers/")
	assert.Equal(t, nil, err, "servers服务节点发布验证")
	_, _, err = lm.GetValue("/hydra/services/api/providers/")
	assert.Equal(t, nil, err, "api服务节点发布验证")
	_, _, err = lm.GetValue("/dns/")
	assert.Equal(t, nil, err, "dns服务节点发布验证")

	//发布rpc节点
	confObj.Service.API.Add("/api1", "/api1", []string{"GET"})
	confObj.Service.API.Add("/api2", "/api1", []string{"GET"})
	rpcconf := confObj.GetRPCConf() //初始化参数
	c2 := rpcconf.GetMainConf()     //获取配置
	lm2 := c2.GetRegistry()
	p2 := pub.New(c2)
	p2.Publish(serverName, serviceAddr, c2.GetServerID(), rpcconf.GetRouterConf().GetPath()...)
	assert.Equal(t, false, err != nil, "发布rpc节点")
	_, _, err = lm2.GetValue("/hydra/rpcserver/rpc/test/servers/")
	assert.Equal(t, nil, err, "rpc服务节点发布验证")
}

func TestPublisher_Update(t *testing.T) {
	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080", api.WithDNS("192.168.0.101"))
	confObj.RPC(":9377")

	serverName := "192.168.5.118:9091"
	serviceAddr := "192.168.5.118:9091"

	apiconf := confObj.GetAPIConf() //初始化参数
	c := apiconf.GetMainConf()      //获取配置
	lm := c.GetRegistry()
	p := pub.New(c)

	//更新不存在的API节点的值
	err := p.Update(serverName, serviceAddr, c.GetServerID(), "key1", "value1")
	assert.Equal(t, false, err != nil, "更新不存在的API节点的值")
	_, _, err = lm.GetValue("/hydra/apiserver/api/test/servers/")
	assert.Equal(t, "节点[%!w(string=/hydra/apiserver/api/test/servers)]不存在", err.Error(), "servers服务节点发布验证")

	//更新已经存在的API节点的不存在值
	p.Publish(serverName, serviceAddr, c.GetServerID())
	err = p.Update(serverName, serviceAddr, c.GetServerID(), "key1", "value1")
	assert.Equal(t, false, err != nil, "更新存在的API节点的值")
	paths, _, _ := lm.GetChildren("/hydra/apiserver/api/test/servers/")
	for _, v := range paths {
		ldata, _, err := lm.GetValue("/hydra/apiserver/api/test/servers/" + v)
		assert.Equal(t, nil, err, "更新存在的API节点的值")
		value := map[string]string{}
		json.Unmarshal(ldata, &value)
		assert.Equal(t, "value1", value["key1"], "更新存在的API节点的值")
	}

	//更新已经存在的api节点存在的值
	err = p.Update(serverName, serviceAddr, c.GetServerID(), "key1", "value1-1")
	assert.Equal(t, false, err != nil, "更新存在的API节点的值")
	paths, _, _ = lm.GetChildren("/hydra/apiserver/api/test/servers/")
	for _, v := range paths {
		ldata, _, err := lm.GetValue("/hydra/apiserver/api/test/servers/" + v)
		assert.Equal(t, nil, err, "更新存在的API节点的值")
		value := map[string]string{}
		json.Unmarshal(ldata, &value)
		assert.Equal(t, "value1-1", value["key1"], "更新存在的API节点的值")
	}

}

//测试自动恢复节点
func TestNew(t *testing.T) {
	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf() //初始化参数
	c := apiconf.GetMainConf()      //获取配置
	lm := c.GetRegistry()

	p := pub.New(c)
	p.Publish("192.168.5.118:9091", "192.168.5.118:9091", c.GetServerID())
	//删除节点
	paths, _, _ := lm.GetChildren("/hydra/apiserver/api/test/servers/")
	for _, v := range paths {
		path := "/hydra/apiserver/api/test/servers/" + v
		err := lm.Delete(path)
		assert.Equal(t, nil, err, "NEW()-测试自动恢复节点")
		fmt.Printf("节点%s已删除", path)
	}

	paths, _, _ = lm.GetChildren("/hydra/apiserver/api/test/servers/")
	assert.Equal(t, 0, len(paths), "NEW()-测试自动恢复节点")

	time.Sleep(time.Second * 10)
	paths, _, _ = lm.GetChildren("/hydra/apiserver/api/test/servers/")
	assert.Equal(t, 0, len(paths), "NEW()-测试自动恢复节点")

	time.Sleep(time.Second * 10)
	paths, _, _ = lm.GetChildren("/hydra/apiserver/api/test/servers/")
	assert.Equal(t, 0, len(paths), "NEW()-测试自动恢复节点")

	time.Sleep(time.Second * 15)
	paths, _, _ = lm.GetChildren("/hydra/apiserver/api/test/servers/")
	assert.Equal(t, 1, len(paths), "NEW()-测试自动恢复节点")

}
