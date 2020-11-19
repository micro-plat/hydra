package registry

import (
	"testing"
	"time"

	"github.com/micro-plat/hydra/registry/watcher"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/hydra/registry/watcher/wvalue"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func TestSingleValueWatcher_Close(t *testing.T) {
	//构建配置对象
	confObj := mocks.NewConf()
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetServerConf()
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	w := wvalue.NewSingleValueWatcher(c.GetRegistry(), c.GetServerPubPath(), log)
	w.Close()
	assert.Equal(t, true, w.Done, "valueWatcher关闭测试")
	_, ok := <-w.CloseChan
	assert.Equal(t, false, ok, "valueWatcher关闭测试")
}

func TestSingleValueWatcher_Start(t *testing.T) {
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
		//	{name: "监控的错误节点的值变动", path: "/a/b/c", r: c.GetRegistry(), wantOp: watcher.ADD, wantErr: false}, //@watch方法 陷入死循环
		// {name: "监控过程中,注册中心节点存在,获取子节点值错误", path: "/platname/apiserver/api/test/hosts1/",
		// r: mocks.NewTestRegistry("platname", "apiserver", "test", ""), wantOp: watcher.ADD, wantErr: false}, //@watch方法 陷入死循环
		{name: "监控过程中,注册中心节点的值未发生改变", path: c.GetServerPubPath(), r: c.GetRegistry(), wantOp: watcher.ADD, wantErr: false},
		{name: "监控过程中,注册中心子节点的值发生改变", path: "/platname/apiserver/api/test/hosts/server1",
			r: mocks.NewTestRegistry("platname", "apiserver", "test", ""), wantOp: watcher.ADD, wantErr: false},
	}

	//发布节点到注册中心
	router, _ := apiconf.GetRouterConf()
	pub.New(c).Publish("192.168.0.1:9091", "192.168.0.2:9091", c.GetServerID(), router.GetPath()...)
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	for _, tt := range tests {
		//变化之前的值
		data, version, err := tt.r.GetValue(tt.path)

		w := wvalue.NewSingleValueWatcher(tt.r, tt.path, log)
		gotC, err := w.Start()
		assert.Equal(t, tt.wantErr, err != nil, tt.name)

		//保存测试退出前,线程执行完
		time.Sleep(time.Second * 2)

		if tt.wantErr {
			continue
		}

		//获取子节点,监控结果验证
		newData, newVersion, err := tt.r.GetValue(tt.path)
		assert.Equal(t, nil, err, tt.name)
	LOOP:
		for {
			select {
			case c := <-gotC:
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
