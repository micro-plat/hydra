package registry

import (
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/hydra/registry/watcher/wchild"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func TestChildWatcher_Close(t *testing.T) {
	confObj := mocks.NewConf()
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetServerConf()
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())
	w := wchild.NewChildWatcher(c.GetRegistry(), c.GetServerPubPath(), log)

	w.Close()
	w.Close()

	for _, v := range w.Watchers {
		assert.Equal(t, true, v.Done, "childWatcher关闭测试")
		_, ok := <-v.CloseChan
		assert.Equal(t, false, ok, "childWatcher关闭测试")
	}
}

func TestChildWatcher_Start(t *testing.T) {

	//构建配置对象
	confObj := mocks.NewConfBy("TestChildWatcher", "start")
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetServerConf()
	addr1 := "192.168.5.115:9091"
	addr2 := "192.168.5.116:9091"

	tests := []struct {
		name    string
		path    string
		deep    int
		r       registry.IRegistry
		wantOp  int
		wantErr bool
	}{
		// {name: "获取错误路径的节点变动", path: "/a/b/c", r: c.GetRegistry(), wantErr: false},  //@watch方法 陷入死循环
		// {name: "监控过程中,注册中心节点存在,获取子节点发生错误", path: "/platname/apiserver/api/test/hosts1/",
		// r: mocks.NewTestRegistry("platname", "apiserver", "test", ""), wantErr: false},   //watch方法陷入死循环
		{name: "深度为1,监控过程中,注册中心子节点未发生改变", path: registry.Join(c.GetServerPubPath(), addr1),
			deep: 1, r: c.GetRegistry(), wantOp: watcher.ADD, wantErr: false},
		{name: "深度为2,监控过程中,注册中心子节点未发生改变", path: c.GetServerPubPath(), deep: 2, r: c.GetRegistry(), wantOp: watcher.ADD, wantErr: false},
		{name: "深度为3,监控过程中,注册中心子节点未发生改变", path: "/hydra/apiserver/api/test/", deep: 3, r: c.GetRegistry(), wantOp: watcher.ADD, wantErr: false},
		{name: "深度为4,监控过程中,注册中心子节点未发生改变", path: "/hydra/apiserver/api/", deep: 4, r: c.GetRegistry(), wantOp: watcher.ADD, wantErr: false},
	}

	//发布节点到注册中心
	router, _ := apiconf.GetRouterConf()
	pub.New(c).Publish(addr1, addr1, c.GetServerID(), router.GetPath()...)
	pub.New(c).Publish(addr2, addr2, c.GetServerID(), router.GetPath()...)
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	for _, tt := range tests {

		//启动子节点监控
		w := wchild.NewChildWatcherByDeep(tt.path, tt.deep, tt.r, log)
		gotC, err := w.Start()

		//保证测试退出前 线程执行完
		time.Sleep(time.Second * 2)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)

		if tt.wantErr {
			continue
		}

		//获取子节点,监控结果验证
	LOOP:
		for {
			select {
			case c := <-gotC:
				//比对当前的节点返回的子节点的信息
				children, version, _ := tt.r.GetChildren(c.Parent)
				assert.Equal(t, version, c.Version, tt.name)
				assert.Equal(t, len(children), len(c.Children), tt.name)

				lk := c.Parent[len(tt.path):]
				d := strings.Split(lk, "/")
				if len(d) == 1 {
					if d[0] == "" {
						assert.Equal(t, tt.deep-len(d)+1, c.Deep, tt.name)
					}
				} else {
					assert.Equal(t, tt.deep-len(d), c.Deep, tt.name)
				}

				assert.Equal(t, tt.wantOp, c.OP, tt.name)
			default:
				break LOOP
			}
		}
	}
}

func TestChildWatcher_Start_2(t *testing.T) {

	//构建配置对象
	confObj := mocks.NewConfBy("TestChildWatcher_Start_2", "TestChildWatcher_Start_2Clu")
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetServerConf()
	addr1 := "192.168.5.115:9091"
	addr2 := "192.168.5.116:9091"

	tests := []struct {
		name    string
		path    string
		deep    int
		r       *mocks.TestRegistry
		wantErr bool
	}{
		{name: "深度为1,监控过程中,注册中心添加了当前节点的子节点", path: "/platname/apiserver/api/test/hosts/server6",
			r: mocks.NewTestRegistry("platname", "apiserver", "test", ""), deep: 1, wantErr: false},
		{name: "深度为2,监控过程中,注册中心添加了当前节点的子节点", path: "/platname/apiserver/api/test/hosts",
			r: mocks.NewTestRegistry("platname", "apiserver", "test", ""), deep: 2, wantErr: false},
		{name: "深度为3,监控过程中,注册中心添加了当前节点的子节点", path: "/platname/apiserver/api/test",
			r: mocks.NewTestRegistry("platname", "apiserver", "test", ""), deep: 3, wantErr: false},
	}

	//发布节点到注册中心
	router, _ := apiconf.GetRouterConf()
	pub.New(c).Publish(addr1, addr1, c.GetServerID(), router.GetPath()...)
	pub.New(c).Publish(addr2, addr2, c.GetServerID(), router.GetPath()...)
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	for _, tt := range tests {
		tt.r.Deep = tt.deep

		//启动子节点监控
		w := wchild.NewChildWatcherByDeep(tt.path, tt.deep, tt.r, log)
		gotC, err := w.Start()

		//保证测试退出前 线程执行完
		time.Sleep(time.Second * 2)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)

		if tt.wantErr {
			continue
		}

		//获取子节点,监控结果验证
	LOOP:
		for {
			select {
			case c := <-gotC:
				//对返回的节点进行验证
				children, _, _ := tt.r.GetChildren(c.Parent)

				if c.OP == watcher.ADD && tt.path == c.Parent { //未添加节点前的返回值
					assert.Equal(t, len(children)-1, len(c.Children), tt.name)
				} else { //添加节点后的返回值
					assert.Equal(t, len(children), len(c.Children), tt.name)
				}
				_, cVersion, _ := tt.r.GetValue(registry.Join(c.Parent, c.Children[0]))
				assert.Equal(t, cVersion, c.Version, tt.name)

				lk := c.Parent[len(tt.path):]
				d := strings.Split(lk, "/")
				assert.Equal(t, tt.deep-len(d)+1, c.Deep, tt.name)

				names := strings.Split(strings.Trim(c.Parent, "/"), "/")
				assert.Equal(t, names[len(names)-1], c.Name, tt.name)

			default:
				break LOOP
			}
		}
	}
}

//节点删除时进行处理
func TestChildWatcher_deleted(t *testing.T) {
	//构建配置对象
	confObj := mocks.NewConfBy("TestChildWatcher_Deleted", "TestChildWatcher_Deleted_Clu")
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf()
	c := apiconf.GetServerConf()
	addr1 := "192.168.5.115:9091"
	addr2 := "192.168.5.116:9091"

	tests := []struct {
		name    string
		path    string
		deep    int
		r       *mocks.TestRegistry
		wantErr bool
	}{
		{name: "深度为2,监控过程中,注册中心删除了当前节点的子节点", path: "/platname/apiserver/api/test/hosts_delete",
			r: mocks.NewTestRegistry("platname", "apiserver", "test", ""), deep: 2, wantErr: false},
	}

	//发布节点到注册中心
	router, _ := apiconf.GetRouterConf()
	pub.New(c).Publish(addr1, addr1, c.GetServerID(), router.GetPath()...)
	pub.New(c).Publish(addr2, addr2, c.GetServerID(), router.GetPath()...)
	log := logger.GetSession(apiconf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", conf.NewMeta()).GetRequestID())

	for _, tt := range tests {
		tt.r.Deep = tt.deep

		//获取当前子节点
		beforeChildren, _, _ := tt.r.GetChildren(tt.path)

		//启动子节点监控
		w := wchild.NewChildWatcherByDeep(tt.path, tt.deep, tt.r, log)
		gotC, err := w.Start()

		//保证测试退出前 线程执行完
		time.Sleep(time.Second * 2)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)

		if tt.wantErr {
			continue
		}

		//获取子节点,监控结果验证
	LOOP:
		for {
			select {
			case c := <-gotC:
				//对返回的节点进行验证
				//children, _, _ := tt.r.GetChildren(c.Parent)

				if c.OP == watcher.ADD && tt.path == c.Parent { //未删除节点前的返回值
					assert.Equal(t, len(beforeChildren), len(c.Children), tt.name)
				} else { //删除节点后的返回值
					assert.Equal(t, len(beforeChildren)-1, len(c.Children), tt.name)
				}
				var version int32
				version = 0
				assert.Equal(t, version, c.Version, tt.name)

				lk := c.Parent[len(tt.path):]
				d := strings.Split(lk, "/")
				assert.Equal(t, tt.deep-len(d)+1, c.Deep, tt.name)

				names := strings.Split(strings.Trim(c.Parent, "/"), "/")
				assert.Equal(t, names[len(names)-1], c.Name, tt.name)

			default:
				break LOOP
			}
		}
	}
}
