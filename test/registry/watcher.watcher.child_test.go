package registry

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func TestNewChildWatcher(t *testing.T) {
	tests := []struct {
		name         string
		registryAddr string
		wantNil      bool
		wantErr      bool
	}{
		{name: "获取支持协议的节点监控对象", registryAddr: "lm://."},
		{name: "获取不支持协议的节点监控对象", registryAddr: "cloud://.", wantNil: true, wantErr: true},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, conf.NewMeta()).GetRequestID())
	for _, tt := range tests {
		gotR, err := watcher.NewChildWatcher(tt.registryAddr, []string{}, log)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.IsNil(t, tt.wantNil, gotR, tt.name)
	}
}

func TestNewChildWatcherByRegistry(t *testing.T) {

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, conf.NewMeta()).GetRequestID())

	r := serverConf.GetServerConf().GetRegistry()
	gotR, err := watcher.NewChildWatcherByRegistry(r, []string{}, log)
	assert.Equal(t, nil, err, "构建ChildWatcher")
	assert.IsNil(t, false, gotR, "构建ChildWatcher")

}
