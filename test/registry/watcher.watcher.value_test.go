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
		{name: "获取支持协议的节点值监控对象", registryAddr: "lm://."},
		{name: "获取不支持协议的节点值监控对象", registryAddr: "cloud://.", wantNil: true, wantErr: true},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())
	for _, tt := range tests {
		gotR, err := watcher.NewValueWatcher(tt.registryAddr, []string{}, log)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.IsNil(t, tt.wantNil, gotR, tt.name)
	}
}

func TestNewValueWatcherByRegistry(t *testing.T) {

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	r := serverConf.GetServerConf().GetRegistry()
	gotR, err := watcher.NewValueWatcherByRegistry(r, []string{}, log)
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
		//	{name: "不传入servers", platName: "hydra", systemName: "apiserver", clusterName: "test", wantPanic: true},
		{name: "传入servers", platName: "hydra", systemName: "apiserver", serverTypes: []string{"api"}, clusterName: "test"},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())
	r := serverConf.GetServerConf().GetRegistry()
	for _, tt := range tests {
		defer func() {
			r := recover()
			assert.Equal(t, tt.wantPanic, r != nil, tt.name)
		}()
		gotR, err := watcher.NewValueWatcherByServers(r, tt.platName, tt.systemName, tt.serverTypes, tt.clusterName, log)
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
