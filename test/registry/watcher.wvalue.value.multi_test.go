package registry

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/hydra/registry/watcher/wvalue"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func TestNewMultiValueWatcher(t *testing.T) {

	//构建配置对象
	confObj := mocks.NewConf()
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetServerConf()
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	w, _ := wvalue.NewMultiValueWatcher(c.GetRegistry(), []string{"a", "b", "c"}, log)
	assert.Equal(t, 3, len(w.Watchers), "构建的值监控对象")
}

func TestMultiValueWatcher_Close(t *testing.T) {
	//构建配置对象
	confObj := mocks.NewConf()
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetServerConf()
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	w, _ := wvalue.NewMultiValueWatcher(c.GetRegistry(), []string{"a", "b", "c"}, log)
	w.Close()
	for _, v := range w.Watchers {
		assert.Equal(t, true, v.Done, "节点值监控对象关闭")
		_, ok := <-v.CloseChan
		assert.Equal(t, false, ok, "节点值监控对象关闭")
	}
}

func TestMultiValueWatcher_Start(t *testing.T) {

	//构建配置对象
	confObj := mocks.NewConf()
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetServerConf()

	tests := []struct {
		name    string
		path    string
		wantOp  int
		r       registry.IRegistry
		wantErr bool
	}{
		{name: "监控过程中,注册中心节点的值未发生改变", path: c.GetServerPubPath(), r: c.GetRegistry(), wantOp: watcher.ADD, wantErr: false},
		{name: "监控过程中,注册中心子节点的值发生改变", path: "/platname/apiserver/api/test/hosts/server1",
			r: mocks.NewTestRegistry("platname", "apiserver", "test", ""), wantOp: watcher.ADD, wantErr: false},
	}

	router, _ := apiconf.GetRouterConf()
	pub.New(c).Publish("192.168.0.1:9091", "192.168.0.2:9091", c.GetServerID(), router.GetPath()...)
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	for _, tt := range tests {

		data, version, err := tt.r.GetValue(tt.path)

		w, _ := wvalue.NewMultiValueWatcher(tt.r, []string{tt.path}, log)
		got, err := w.Start()
		assert.Equal(t, nil, err, tt.name)

		//获取子节点,监控结果验证
	LOOP:
		for {
			select {
			case c := <-got:
				newData, newVersion, err := tt.r.GetValue(tt.path)
				assert.Equal(t, nil, err, tt.name)

				if c.OP == watcher.ADD {
					assert.Equal(t, version, c.Version, tt.name)
					assert.Equal(t, data, c.Content, tt.name)
				}
				if c.OP == watcher.CHANGE {
					assert.Equal(t, newVersion, c.Version, tt.name)
					assert.Equal(t, newData, c.Content, tt.name)
				}

				assert.Equal(t, tt.path, c.Path, tt.name)
				tt.wantOp = watcher.CHANGE
			default:
				break LOOP
			}
		}
	}

}
