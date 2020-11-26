package registry

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func TestNewValueWatcher(t *testing.T) {
	tests := []struct {
		name         string
		registryAddr string
		wantNil      bool
		wantErr      bool
	}{
		{name: "1.1 ValueWatcher-获取对象失败-获取不支持协议的节点监控对象", registryAddr: "cloud://.", wantNil: true, wantErr: true},
		{name: "1.2 ValueWatcher-获取对象失败-错误的协议格式", registryAddr: "lm//dd", wantNil: true, wantErr: true},
		{name: "1.3 ValueWatcher-获取对象失败-没有协议名", registryAddr: "://dd", wantNil: true, wantErr: true},
		{name: "1.4 ValueWatcher-获取对象失败-没有服务节点地址", registryAddr: "fs://", wantNil: true, wantErr: true},

		{name: "2.1 ValueWatcher-获取对象成功-基于zk", registryAddr: "zk://192.168.0.101"},
		{name: "2.2 ValueWatcher-获取对象成功-基于lm", registryAddr: "lm://."},
		{name: "2.3 ValueWatcher-获取对象成功-基于fs", registryAddr: "fs://."},
	}

	confObj := mocks.NewConfBy("hydra_rgst_watcher_value", "rgtwatchevaluetest") //构建对象
	confObj.API(":8080")                                                         //初始化参数
	serverConf := confObj.GetAPIConf()                                           //获取配置
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())
	for _, tt := range tests {
		gotR, err := watcher.NewValueWatcher(tt.registryAddr, []string{}, log)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.IsNil(t, tt.wantNil, gotR, tt.name)
	}
}

func TestNewValueWatcherByRegistry(t *testing.T) {
	confObj := mocks.NewConfBy("hydra_rgst_watcher_value1", "rgtwatchevaluetest1") //构建对象
	log := logger.GetSession("hydra_rgst_watcher_value1", ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	gotR, err := watcher.NewValueWatcherByRegistry(confObj.Registry, []string{}, log)
	assert.Equal(t, nil, err, "构建ValueWatcher")
	assert.IsNil(t, false, gotR, "构建ValueWatcher")

}

func TestNewValueWatcherByServers(t *testing.T) {
	tests := []struct {
		name        string
		platName    string
		systemName  string
		serverTypes []string
		clusterName string
		wantPanic   bool
	}{
		{name: "1. ValueWatcherByServers-初始化对象-serverTypes不存在", platName: "hydra", systemName: "apiserver", clusterName: "test", wantPanic: true},
		{name: "2. ValueWatcherByServers-初始化对象-platName不存在", platName: "", systemName: "apiserver", serverTypes: []string{"api"}, clusterName: "test", wantPanic: true},
		{name: "3. ValueWatcherByServers-初始化对象-systemName不存在", platName: "hydra", systemName: "", clusterName: "test", serverTypes: []string{"api"}, wantPanic: true},
		{name: "4. ValueWatcherByServers-初始化对象-clusterName不存在", platName: "hydra", systemName: "apiserver", clusterName: "", serverTypes: []string{"api"}, wantPanic: true},
		{name: "5. ValueWatcherByServers-初始化正确对象", platName: "hydra", systemName: "apiserver", serverTypes: []string{"api"}, clusterName: "test"},
	}

	confObj := mocks.NewConfBy("hydra_rgst_watcher_value2", "rgtwatchevaluetest2") //构建对象
	log := logger.GetSession("hydra_rgst_watcher_value2", ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())
	for _, tt := range tests {
		defer func() {
			r := recover()
			assert.Equal(t, tt.wantPanic, r != nil, tt.name)
		}()
		gotR, err := watcher.NewValueWatcherByServers(confObj.Registry, tt.platName, tt.systemName, tt.serverTypes, tt.clusterName, log)
		assert.Equal(t, nil, err, tt.name)
		assert.IsNil(t, false, gotR, tt.name)
		gotC, err := gotR.Start()
		assert.Equal(t, nil, err, tt.name)

	LOOP:
		for {
			select {
			case c := <-gotC:
				assert.Equal(t, registry.Join(tt.platName, tt.systemName, tt.serverTypes[0], tt.clusterName, "conf"), c.Path, "构建结果验证")
			default:
				break LOOP
			}
		}
	}
}
