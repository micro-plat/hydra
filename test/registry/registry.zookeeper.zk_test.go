package registry

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func Test_zookeeperFactory_Create(t *testing.T) {

	tests := []struct {
		name      string
		addr      string
		path      string
		isCreate  bool
		isNil     bool
		wantErr1  bool
		errStr    string
		wantErr2  bool
		wantValue string
	}{
		{name: "1. 初始化zk，地址为空", addr: "", wantErr1: true, errStr: "，包含多个://。格式:[proto]://[address]"},
		{name: "2. 初始化zk，不指定服务器ip", addr: "zk://", wantErr1: true, errStr: "zk://，地址不能为空。格式:[proto]://[address]"},
		{name: "3. 初始化zk，指定错误的proto", addr: "z1k://", wantErr1: true, errStr: "z1k://，地址不能为空。格式:[proto]://[address]"},
		{name: "4. 初始化zk，指定测试服务器,获取不存在的节点值", addr: "zk://192.168.0.101", path: "/hydratest_rg/apiserver1/api/test/conf1", wantErr2: true},
		{name: "5. 初始化zk，指定测试服务器,获取存在的节点值", addr: "zk://192.168.0.101", path: "/hydratest_rg/apiserver2/api/test/conf", isCreate: true, wantValue: `{"address":":51001"}`},
		{name: "6. 初始化zk，指定测试服务器,获取存在的节点值", addr: "zk://192.168.0.101", path: "/hydratest_rg/apiserver3/api", isCreate: true, wantValue: ``},
	}

	confObj := mocks.NewConfBy("hydra_rgst_zook_test", "rgtzooktest")
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	for _, tt := range tests {
		z, err := registry.NewRegistry(tt.addr, log)
		assert.Equal(t, tt.wantErr1, err != nil, tt.name)
		if tt.wantErr1 {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
			continue
		}

		if tt.isCreate {
			err := z.CreateTempNode(tt.path, tt.wantValue)
			assert.Equal(t, nil, err, tt.name, err)
		}
		//测试连接是否正常
		data, _, err := z.GetValue(tt.path)
		assert.Equal(t, tt.wantErr2, err != nil, tt.name)
		assert.Equal(t, tt.wantValue, string(data), tt.name)
	}
}
