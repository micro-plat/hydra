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
		wantErr1  bool
		errStr    string
		wantErr2  bool
		wantValue string
	}{
		{name: "不指定服务器", addr: "zk://", wantErr1: true, errStr: "zk://，地址不能为空。格式:[proto]://[address]"},
		{name: "指定测试服务器,获取不存在的节点值", addr: "zk://192.168.0.101", path: "/hydra/apiserver/api/test/conf1", wantErr2: true},
		{name: "指定测试服务器,获取存在的节点值", addr: "zk://192.168.0.101", path: "/hydra/apiserver/api/test/conf", wantValue: `{"address":":51001"}`},
		{name: "指定测试服务器,获取存在的节点值", addr: "zk://192.168.0.101", path: "/hydra/apiserver/api", wantValue: ``},
	}

	confObj := mocks.NewConf()
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, conf.NewMeta()).GetRequestID())

	for _, tt := range tests {
		z, err := registry.NewRegistry(tt.addr, log)
		assert.Equal(t, tt.wantErr1, err != nil, tt.name)
		if tt.wantErr1 {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
			continue
		}
		//测试连接是否正常
		data, _, err := z.GetValue(tt.path)
		assert.Equal(t, tt.wantErr2, err != nil, tt.name)
		assert.Equal(t, tt.wantValue, string(data), tt.name)
	}
}
