package rpc

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

type testServer struct {
	addr string
	s    net.Listener
}

func (t *testServer) testListen() {
	t.s, _ = net.Listen("tcp", "0.0.0.0:64120")
}

func (t *testServer) close() {
	t.s.Close()
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

func TestNewResponsive(t *testing.T) {
	confObj := mocks.NewConfBy("rpcserver_resserivece_test", "testrpacsdf") //构建对象
	confObj.RPC(":35211")
	tests := []struct {
		name    string
		cnf     app.IAPPConf
		proto   string
		addr    string
		wantErr bool
	}{
		{name: "1. 构建rpc服务", addr: ":35211", proto: "tcp", cnf: confObj.GetRPCConf()},
	}
	for _, tt := range tests {
		gotH, err := rpc.NewResponsive(tt.cnf)
		assert.Equal(t, nil, err, tt.name)
		addr := fmt.Sprintf("%s://%s%s", tt.proto, global.LocalIP(), tt.addr)
		assert.Equal(t, addr, gotH.Server.GetAddress(), tt.name)
	}
}

func TestResponsive_Start(t *testing.T) {
	confObj := mocks.NewConfBy("rpcserver_resserivece_test1", "testrpacsdf1") //构建对象
	confObj.RPC(":64120")
	reg := confObj.Registry
	tests := []struct {
		name          string
		cnf           app.IAPPConf
		serverName    string
		serverType    string
		starting      func(app.IAPPConf) error
		closing       func(app.IAPPConf) error
		isConfStart   bool //禁用服务
		isServerStart bool //http服务启动失败
		isServerPub   bool //服务发布失败
		wantErr       string
		wantSubErr    string
		wantLog       string
	}{
		{name: "1. 启动rpc服务-starting报错", cnf: confObj.GetRPCConf(), serverType: "rpc", starting: func(app.IAPPConf) error { return fmt.Errorf("err") }, closing: func(app.IAPPConf) error { return nil }, wantErr: "err"},
		{name: "2. 启动rpc服务-禁用服务", cnf: confObj.GetRPCConf(), serverType: "rpc", serverName: "rpcserver", starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return nil }, isConfStart: true, wantLog: "rpc被禁用，未启动"},
		{name: "3. 启动rpc服务-失败", cnf: confObj.GetRPCConf(), serverType: "rpc", serverName: "rpcserver", starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return nil }, isServerStart: true, wantErr: "rpc启动失败 listen tcp 192.168.5.94:64120: bind: address already in use"},
		{name: "4. 启动rpc服务-注册中心服务发布失败", cnf: confObj.GetRPCConf(), serverType: "rpc", serverName: "rpcserver", starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return fmt.Errorf("closing_err") }, isServerPub: true, wantSubErr: "rpc服务发布失败 服务发布失败:", wantLog: "关闭[closing_err]服务,出现错误"},
		{name: "5. 启动rpc服务-成功", cnf: confObj.GetRPCConf(), serverType: "rpc", serverName: "rpcserver", starting: func(app.IAPPConf) error { return nil }, closing: func(app.IAPPConf) error { return nil }, wantLog: "启动成功(rpc,"},
	}
	for _, tt := range tests {
		//starting
		testInitServicesDef()
		services.Def.OnStarting(tt.starting, "rpc")
		//closing
		services.Def.OnClosing(tt.closing, "rpc")

		//禁用服务
		if tt.isConfStart {
			path := fmt.Sprintf("/rpcserver_resserivece_test1/%s/%s/testrpacsdf1/conf", tt.serverName, tt.serverType)
			err := reg.Update(path, `{"address":":64120","status":"stop"}`)
			assert.Equal(t, nil, err, tt.name+"禁用服务")
			tt.cnf, _ = app.NewAPPConf(path, reg)
		}

		//占用端口使服务启动失败
		var rpcServer *testServer
		if tt.isServerStart {
			rpcServer = &testServer{addr: "127.0.0.1:64120"}
			go rpcServer.testListen()
			time.Sleep(time.Second)
		}

		//创建节点使服务发布报错
		if tt.isServerPub {
			newConfObj := mocks.NewConfBy("hydra", "test", "fs://./") //构建对象
			newConfObj.RPC(":64120")
			tt.cnf = newConfObj.GetRPCConf()
			path := fmt.Sprintf("./hydra/%s/%s/test/servers", tt.serverName, tt.serverType)
			os.RemoveAll(path) //删除文件夹
			os.Create(path)    //使文件夹节点变成文件节点,让该节点下不能创建文件
		}

		//构建服务器
		rsp, _ := rpc.NewResponsive(tt.cnf)

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
		//释放端口
		if tt.isServerStart {
			rpcServer.close()
			time.Sleep(time.Second * 1)
		}

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
	confObj := mocks.NewConfBy("rpcserver_resserivece_test2", "testrpacsdf2") //构建对象
	confObj.RPC(":64121")                                                     //初始化参数
	cnf := confObj.GetRPCConf()
	rsp, err := rpc.NewResponsive(cnf)

	assert.Equal(t, nil, err, "构建服务错误")
	//节点未变动
	tChange, err := rsp.Notify(cnf)
	assert.Equal(t, nil, err, "通知变动错误")
	assert.Equal(t, false, tChange, "通知变动判断")

	path := "/rpcserver_resserivece_test2/rpcserver/rpc/testrpacsdf2/conf"
	registry, err := registry.NewRegistry("lm://./", logger.New("hydra"))
	//节点进行值变更 进行启动
	err = registry.Update(path, `{"status":"start","addr":":64121"}`)
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
	err = registry.Update(path, `{"status":"stop","addr":":64121"}`)
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
