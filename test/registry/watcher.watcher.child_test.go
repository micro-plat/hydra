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

		{name: "1.1 ChildWatcher-获取对象失败-获取不支持协议的节点监控对象", registryAddr: "cloud://.", wantNil: true, wantErr: true},
		{name: "1.2 ChildWatcher-获取对象失败-错误的协议格式", registryAddr: "lm//dd", wantNil: true, wantErr: true},
		{name: "1.3 ChildWatcher-获取对象失败-没有协议名", registryAddr: "://dd", wantNil: true, wantErr: true},
		{name: "1.4 ChildWatcher-获取对象失败-没有服务节点地址", registryAddr: "fs://", wantNil: true, wantErr: true},

		{name: "2.1 ChildWatcher-获取对象成功-基于zk", registryAddr: "zk://192.168.0.101"},
		{name: "2.2 ChildWatcher-获取对象成功-基于lm", registryAddr: "lm://."},
		{name: "2.3 ChildWatcher-获取对象成功-基于fs", registryAddr: "fs://."},
	}

	confObj := mocks.NewConfBy("hydra_rgst_watcher_clid", "rgtwatcherclidtest") //构建对象
	confObj.API(":8080")                                                        //初始化参数
	serverConf := confObj.GetAPIConf()                                          //获取配置
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())
	for _, tt := range tests {
		gotR, err := watcher.NewChildWatcher(tt.registryAddr, []string{}, log)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.IsNil(t, tt.wantNil, gotR, tt.name)
	}
}

func TestNewChildWatcherByRegistry(t *testing.T) {

	confObj := mocks.NewConfBy("hydra_rgst_watcher_clid1", "rgtwatcherclidtest1") //构建对象
	log := logger.GetSession("hydra_rgst_watcher_clid1", ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	gotR, err := watcher.NewChildWatcherByRegistry(confObj.Registry, []string{}, log)
	assert.Equal(t, nil, err, "构建ChildWatcher")
	assert.IsNil(t, false, gotR, "构建ChildWatcher")
}
