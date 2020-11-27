package cron

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

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

func TestNewResponsive(t *testing.T) {
	confObj := mocks.NewConfBy("cronserver_resserivece_test", "testcronsdf") //构建对象
	confB := confObj.CRON()
	confB.Task(task.NewTask("@every 10s", "/cron/server1"), task.NewTask("@every 10s", "/cron/server2"))
	server, err := cron.NewServer(task.NewTask("@every 10s", "/cron/server1"), task.NewTask("@every 10s", "/cron/server2"))
	assert.Equalf(t, nil, err, "获取server对象异常")
	tests := []struct {
		name    string
		cnf     app.IAPPConf
		wantH   *cron.Server
		wantErr bool
	}{
		{name: "1. 初始化实体对象", cnf: confObj.GetCronConf(), wantH: server, wantErr: true},
	}
	for _, tt := range tests {
		_, err := cron.NewResponsive(tt.cnf)
		assert.Equalf(t, tt.wantErr, err == nil, tt.name+",err")
	}
}

func TestResponsive_Start(t *testing.T) {
	confObj := mocks.NewConfBy("cronserver_resserivece_test1", "testcronsdf1") //构建对象
	confObj.CRON()
	reg := confObj.Registry
	tests := []struct {
		name        string
		cnf         app.IAPPConf
		starting    func(app.IAPPConf) error
		closing     func(app.IAPPConf) error
		isConfStart bool //禁用服务
		isServerPub bool //服务发布失败
		wantErr     string
		wantSubErr  string
		wantLog     string
	}{
		{name: "1. 启动cron服务-starting报错", cnf: confObj.GetCronConf(), starting: func(app.IAPPConf) error { return fmt.Errorf("err") }, closing: func(app.IAPPConf) error { return nil }, wantErr: "err"},
		{name: "2. 启动cron服务-禁用服务", cnf: confObj.GetCronConf(), starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return nil }, isConfStart: true, wantLog: "cron被禁用，未启动"},
		{name: "3. 启动cron服务-注册中心服务发布失败", cnf: confObj.GetCronConf(), starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return fmt.Errorf("closing_err") }, isServerPub: true, wantSubErr: "cron服务发布失败 服务发布失败:", wantLog: "关闭[closing_err]服务,出现错误"},
		{name: "4. 启动cron服务-启动服务成功", cnf: confObj.GetCronConf(), starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return nil }, wantLog: "启动成功(cron,"},
	}
	for _, tt := range tests {
		//starting
		testInitServicesDef()
		services.Def.OnStarting(tt.starting, "cron")
		//closing
		services.Def.OnClosing(tt.closing, "cron")

		//禁用服务
		if tt.isConfStart {
			path := "/cronserver_resserivece_test1/cronserver/cron/testcronsdf1/conf"
			err := reg.Update(path, `{"status":"stop"}`)
			assert.Equal(t, nil, err, tt.name+"禁用服务")
			tt.cnf, _ = app.NewAPPConf(path, reg)
		}

		//创建节点使服务发布报错
		if tt.isServerPub {
			newConfObj := mocks.NewConfBy("cronserver_resserivece_test1", "testcronsdf1", "fs://./") //构建对象
			newConfObj.CRON()
			tt.cnf = newConfObj.GetCronConf()
			path := "./cronserver_resserivece_test1/cronserver/cron/testcronsdf1/servers"
			os.RemoveAll(path) //删除文件夹
			os.Create(path)    //使文件夹节点变成文件节点,让该节点下不能创建文件
		}

		//构建服务器
		rsp, _ := cron.NewResponsive(tt.cnf)

		//构建的新的os.Stdout
		r, w, _ := os.Pipe()
		rescueStdout := newTestStdOut(r, w)

		//启动服务器
		err := rsp.Start()

		//等待日志打印完成
		time.Sleep(time.Second)

		//获取输出并还原os.Stdout
		w.Close()
		out, _ := ioutil.ReadAll(r)
		rescueTestStdout(rescueStdout)
		// fmt.Println("xxxx:", string(out))

		//删除节点文件
		if tt.isServerPub {
			os.RemoveAll("./cronserver_resserivece_test1") //删除文件夹
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
	confObj := mocks.NewConfBy("cronserver_resserivece_test2", "testcronsdf2") //构建对象
	confObj.CRON()                                                             //初始化参数
	cnf := confObj.GetCronConf()
	rsp, err := cron.NewResponsive(cnf)

	assert.Equal(t, nil, err, "构建服务错误")
	//节点未变动
	tChange, err := rsp.Notify(cnf)
	assert.Equal(t, nil, err, "通知变动错误")
	assert.Equal(t, false, tChange, "通知变动判断")

	path := "/cronserver_resserivece_test2/cronserver/cron/testcronsdf2/conf"
	registry, err := registry.NewRegistry("lm://./", logger.New("hydra"))
	//节点进行值变更 进行启动
	err = registry.Update(path, `{"status":"start"}`)
	assert.Equalf(t, false, err != nil, "更新节点2")
	time.Sleep(time.Second * 1)
	conf, err := app.NewAPPConf(path, registry)
	assert.Equalf(t, false, err != nil, "获取最新配置2")
	tChange, err = rsp.Notify(conf)
	time.Sleep(time.Second)
	assert.Equal(t, nil, err, "通知变动错误2")
	assert.Equal(t, true, tChange, "通知变动判断2")

	//节点进行值变更 不用重启
	assert.Equalf(t, false, err != nil, "获取注册中心")
	err = registry.Update(path, `{"status":"stop"}`)
	assert.Equalf(t, false, err != nil, "更新节点")
	time.Sleep(time.Second * 1)
	conf, err = app.NewAPPConf(path, registry)
	assert.Equalf(t, false, err != nil, "获取最新配置")
	tChange, err = rsp.Notify(conf)
	time.Sleep(time.Second)
	assert.Equal(t, nil, err, "通知变动错误")
	assert.Equal(t, true, tChange, "通知变动判断")
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
