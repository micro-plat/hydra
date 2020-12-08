package mqc

import (
	"fmt"
	"io/ioutil"
	xhttp "net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/conf/vars/cache/cacheredis"
	"github.com/urfave/cli"

	varredis "github.com/micro-plat/hydra/conf/vars/redis"

	_ "github.com/micro-plat/hydra/components/caches/cache/redis"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/mqc"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func TestNewResponsive(t *testing.T) {
	confObj := mocks.NewConf() //构建对象
	confObj.Vars().Redis("5.79", varredis.New([]string{"192.168.5.79:6379"}))
	confObj.Vars().Queue().Redis("xxx", queueredis.New(queueredis.WithConfigName("5.79")))
	confObj.MQC("redis://xxx") //初始化参数
	tests := []struct {
		name    string
		proto   string
		cnf     app.IAPPConf
		wantErr bool
	}{
		{name: "1. 构建mqc服务", proto: "mqc", cnf: confObj.GetMQCConf()},
	}
	for _, tt := range tests {
		gotH, err := mqc.NewResponsive(tt.cnf)
		assert.Equal(t, nil, err, tt.name)
		addr := fmt.Sprintf("%s://%s", tt.proto, global.LocalIP())
		assert.Equal(t, addr, gotH.Server.GetAddress(), tt.name)
	}
}

func testInitServicesDef() {
	services.Def = services.New()
	services.Def.RegisterServer("api", func(g *services.Unit, ext ...interface{}) error {
		return services.API.Add(g.Path, g.Service, g.Actions, ext...)
	})
	services.Def.RegisterServer("ws", func(g *services.Unit, ext ...interface{}) error {
		return services.WS.Add(g.Path, g.Service, g.Actions, ext...)
	})
	services.Def.RegisterServer("web", func(g *services.Unit, ext ...interface{}) error {
		return services.WEB.Add(g.Path, g.Service, g.Actions, ext...)
	})
	services.Def.RegisterServer("rpc", func(g *services.Unit, ext ...interface{}) error {
		return services.RPC.Add(g.Path, g.Service, g.Actions, ext...)
	})

	services.Def.RegisterServer("cron", func(g *services.Unit, ext ...interface{}) error {
		for _, t := range ext {
			services.CRON.Add(t.(string), g.Service)
		}
		return nil
	})
	services.Def.RegisterServer("mqc", func(g *services.Unit, ext ...interface{}) error {
		for _, t := range ext {
			services.MQC.Add(t.(string), g.Service)
		}
		return nil
	})
}

var osStdOutLock sync.Mutex

func newTestStdOut(r *os.File, w *os.File) (rescueStdout *os.File) {
	osStdOutLock.Lock()
	rescueStdout = os.Stdout
	os.Stdout = w
	return
}

func rescueTestStdout(rescueStdout *os.File) {
	defer osStdOutLock.Unlock()
	os.Stdout = rescueStdout
	return
}

type testServer struct {
	addr string
	s    *xhttp.Server
}

func (t *testServer) testListen() {
	t.s = &xhttp.Server{Addr: t.addr, Handler: nil}
	t.s.ListenAndServe()
}

func (t *testServer) close() {
	t.s.Close()
}

func TestResponsive_Start(t *testing.T) {
	confObj := mocks.NewConf() //构建对象
	confObj.Vars().Redis("5.79", varredis.New([]string{"192.168.5.79:6379"}))
	confObj.Vars().Queue().Redis("xxx", queueredis.New(queueredis.WithConfigName("5.79")))
	confObj.MQC("redis://xxx") //初始化参数
	reg := confObj.Registry
	tests := []struct {
		name           string
		cnf            app.IAPPConf
		serverName     string
		serverType     string
		starting       func(app.IAPPConf) error
		closing        func(app.IAPPConf) error
		isConfStart    bool //禁用服务
		isGetCluster   bool //集群获取错误
		isGetMainConf  bool //主配置获取错误
		isServerResume bool //server回复错误
		isServerPub    bool //服务发布失败
		wantErr        string
		wantSubErr     string
		wantLog        string
	}{
		{name: "1. 启动mqc服务-starting报错", cnf: confObj.GetMQCConf(), serverType: "mqc", starting: func(app.IAPPConf) error { return fmt.Errorf("err") }, closing: func(app.IAPPConf) error { return nil }, wantErr: "err"},
		{name: "2. 启动mqc服务-禁用服务", cnf: confObj.GetMQCConf(), serverType: "mqc", serverName: "mqcserver", starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return nil }, isConfStart: true, wantLog: "mqc被禁用，未启动"},
		{name: "3. 启动mqc服务-mqc服务获取集群监控失败", cnf: confObj.GetMQCConf(), serverType: "mqc", serverName: "mqcserver", starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return nil }, isGetCluster: true, wantLog: "当前集群节点不可用"},
		{name: "4. 启动mqc服务-mqc服务恢复失败", cnf: confObj.GetMQCConf(), serverType: "mqc", serverName: "mqcserver", starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return nil }, isServerResume: true, wantLog: "恢复mqc服务器失败: 队列名字不能为空"},
		{name: "5. 启动mqc服务-注册中心服务发布失败", cnf: confObj.GetMQCConf(), serverType: "mqc", serverName: "mqcserver", starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return fmt.Errorf("closing_err") }, isServerPub: true, wantSubErr: "mqc服务发布失败 服务发布失败:", wantLog: "关闭[closing_err]服务,出现错误"},
		{name: "6. 启动mqc服务-启动mqc服务成功", cnf: confObj.GetMQCConf(), serverType: "mqc", serverName: "mqcserver", starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return nil }, wantLog: "启动成功(mqc,mqc://"},
	}
	for _, tt := range tests {

		//starting
		testInitServicesDef()
		services.Def.OnStarting(tt.starting)

		//closing
		services.Def.OnClosing(tt.closing)

		//禁用服务
		if tt.isConfStart {
			path := fmt.Sprintf("/hydra/%s/%s/test/conf", tt.serverName, tt.serverType)
			err := reg.Update(path, `{"address":":55004","status":"stop"}`)
			assert.Equal(t, nil, err, tt.name+"禁用服务")
			tt.cnf, _ = app.NewAPPConf(path, reg)
		}

		//创建节点使服务发布报错
		if tt.isServerPub {
			newConfObj := mocks.NewConfBy("hydra", "test", "fs://./") //构建对象
			newConfObj.Vars().Redis("5.79", varredis.New([]string{"192.168.5.79:6379"}))
			newConfObj.Vars().Queue().Redis("xxx", queueredis.New(queueredis.WithConfigName("5.79")))
			newConfObj.MQC("redis://xxx") //初始化参数
			tt.cnf = newConfObj.GetMQCConf()
			path := fmt.Sprintf("./hydra/%s/%s/test/servers", tt.serverName, tt.serverType)
			os.RemoveAll(path) //删除文件夹
			os.Create(path)    //使文件夹节点变成文件节点,让该节点下不能创建文件
		}

		//构建服务器
		rsp, _ := mqc.NewResponsive(tt.cnf)

		//构建的新的os.Stdout
		r, w, _ := os.Pipe()
		rescueStdout := newTestStdOut(r, w)

		//添加空队列使server服务恢复失败
		if tt.isServerResume {
			rsp.Server.Processor.Add(queue.NewQueue("", "services1"))
		}

		//启动服务器
		err := rsp.Start()

		//删除集群配置
		if tt.isGetCluster {
			err := reg.Delete(fmt.Sprintf("/hydra/%s/%s/test/servers", tt.serverName, tt.serverType))
			fmt.Println("y:", err)
		}

		//等待日志打印完成
		time.Sleep(time.Second)

		//获取输出并还原os.Stdout
		w.Close()
		out, _ := ioutil.ReadAll(r)
		rescueTestStdout(rescueStdout)

		//删除节点文件
		if tt.isServerPub {
			os.RemoveAll("./hydra") //删除文件夹
		}

		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name+"err")
			continue
		}

		if tt.wantSubErr != "" {
			assert.Equal(t, true, strings.Contains(err.Error(), tt.wantSubErr), tt.name+"sub_err")
		} else {
			assert.Equal(t, nil, err, tt.name)
		}

		if tt.wantLog != "" {
			assert.Equalf(t, true, strings.Contains(string(out), tt.wantLog), tt.name+"log")
		}
	}
}

func TestResponsive_Notify(t *testing.T) {
	confObj := mocks.NewConf() //构建对象
	confObj.Vars().Redis("5.79", varredis.New([]string{"192.168.5.79:6379"}))
	confObj.Vars().Queue().Redis("xxx", queueredis.New(queueredis.WithConfigName("5.79")))
	confObj.MQC("redis://xxx") //初始化参数
	cnf := confObj.GetMQCConf()
	rsp, err := mqc.NewResponsive(cnf)

	assert.Equal(t, nil, err, "构建服务错误")
	//节点未变动
	tChange, err := rsp.Notify(cnf)
	assert.Equal(t, nil, err, "通知变动错误")
	assert.Equal(t, false, tChange, "通知变动判断")

	path := "/hydra/mqcserver/mqc/test/conf"
	registry, err := registry.GetRegistry("lm://./", logger.New("hydra"))
	//节点进行值变更 进行启动
	err = registry.Update(path, `{"stat11us":"start","addr":"redis://xxx"}`)
	assert.Equalf(t, false, err != nil, "更新节点2")
	time.Sleep(time.Second * 1)
	conf, err := app.NewAPPConf(path, registry)
	assert.Equalf(t, false, err != nil, "获取最新配置2")
	tChange, err = rsp.Notify(conf)
	time.Sleep(time.Second * 2)
	assert.Equal(t, nil, err, "通知变动错误2")
	assert.Equal(t, true, tChange, "通知变动判断2")

	//节点进行值变更 不用重启
	assert.Equalf(t, false, err != nil, "获取注册中心")
	err = registry.Update(path, `{"status":"stop","addr":"redis://xxx"}`)
	assert.Equalf(t, false, err != nil, "更新节点")
	time.Sleep(time.Second * 1)
	conf, err = app.NewAPPConf(path, registry)
	assert.Equalf(t, false, err != nil, "获取最新配置")
	tChange, err = rsp.Notify(conf)
	time.Sleep(time.Second)
	assert.Equal(t, nil, err, "通知变动错误")
	assert.Equal(t, true, tChange, "通知变动判断")

}

func TestResponsive_Start_2(t *testing.T) {

	testInitServicesDef()

	confObj := mocks.NewConfBy("hydra_mqc", "test") //构建对象
	confObj.Vars().Redis("5.79", varredis.New([]string{"192.168.5.79:6379"}))
	confObj.Vars().Queue().Redis("xxx", queueredis.New(queueredis.WithConfigName("5.79")))
	confObj.Vars().Cache().Redis("xxx", cacheredis.New(cacheredis.WithConfigName("5.79")))
	confObj.MQC("redis://xxx").Queue(queue.NewQueue("queue1", "/mqc/test/service1"), queue.NewQueue("queue2", "/mqc/test/service2"))
	global.FlagVal.PlatName = "hydra_mqc"
	global.FlagVal.ClusterName = "test"
	global.Def.Bind(&cli.Context{})

	//注册服务
	services.Def.MQC("services1", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("services1ss")
		return "success"
	}, "queue1")
	services.Def.MQC("services2", func(ctx context.IContext) (r interface{}) {
		ctx.Log().Info("services1ss")
		return "success"
	}, "queue2")

	cnf := confObj.GetMQCConf()

	//往redis添加数据
	app.Cache.Save(cnf)
	c := components.Def.Cache().GetRegularCache("xxx")
	for i := 0; i < 10000; i++ {
		c.Add("queue1", "value1", 120)
		c.Add("queue2", "value2", 120)
	}

	//构建服务器
	rsp, _ := mqc.NewResponsive(cnf)

	//启动服务器
	rsp.Start()

	//等待日志打印完成
	time.Sleep(time.Second * 1)

}
