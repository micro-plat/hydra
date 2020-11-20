package http

import (
	"fmt"
	"io/ioutil"
	xhttp "net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNewResponsive(t *testing.T) {
	confObj := mocks.NewConf() //构建对象
	confObj.API(":55003")      //初始化参数
	confObj.Web(":55002")      //初始化参数
	confObj.WS(":55001")       //初始化参数
	tests := []struct {
		name    string
		proto   string
		addr    string
		cnf     app.IAPPConf
		wantErr bool
	}{
		{name: "构建ws服务", addr: ":55001", proto: "ws", cnf: confObj.GetWSConf()},
		{name: "构建web服务", addr: ":55002", proto: "http", cnf: confObj.GetWebConf()},
		{name: "构建api服务", addr: ":55003", proto: "http", cnf: confObj.GetAPIConf()},
	}
	for _, tt := range tests {
		gotH, err := http.NewResponsive(tt.cnf)
		assert.Equal(t, nil, err, tt.name)
		addr := fmt.Sprintf("%s://%s%s", tt.proto, global.LocalIP(), tt.addr)
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
	confObj.API(":55004")      //初始化参数
	confObj.WS(":55005")       //初始化参数
	confObj.Web(":55006")      //初始化参数
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
		{name: "starting报错",
			cnf:        confObj.GetAPIConf(),
			serverType: "api",
			starting:   func(app.IAPPConf) error { return fmt.Errorf("err") },
			closing:    func(app.IAPPConf) error { return nil },
			wantErr:    "err",
		},
		{name: "禁用服务",
			cnf:         confObj.GetAPIConf(),
			serverType:  "api",
			serverName:  "apiserver",
			starting:    func(app.IAPPConf) error { return nil },
			closing:     func(app.IAPPConf) error { return nil },
			isConfStart: true,
			wantLog:     "api被禁用，未启动",
		},
		{name: "http服务启动失败",
			cnf:           confObj.GetAPIConf(),
			serverType:    "api",
			serverName:    "apiserver",
			starting:      func(app.IAPPConf) error { return nil },
			closing:       func(app.IAPPConf) error { return nil },
			isServerStart: true,
			wantErr:       "api启动失败 listen tcp 0.0.0.0:55004: bind: address already in use",
		},
		{name: "注册中心服务发布失败",
			cnf:         confObj.GetAPIConf(),
			serverType:  "api",
			serverName:  "apiserver",
			starting:    func(app.IAPPConf) error { return nil },
			closing:     func(app.IAPPConf) error { return fmt.Errorf("closing_err") },
			isServerPub: true,
			wantSubErr:  "api服务发布失败 服务发布失败:",
			wantLog:     "关闭[closing_err]服务,出现错误",
		},
		{name: "启动api服务成功",
			cnf:        confObj.GetAPIConf(),
			serverType: "api",
			serverName: "apiserver",
			starting:   func(app.IAPPConf) error { return nil },
			closing:    func(app.IAPPConf) error { return nil },
			wantLog:    "启动成功(api,http:",
		},
		{name: "启动ws服务成功",
			cnf:        confObj.GetWSConf(),
			serverType: "ws",
			serverName: "wsserver",
			starting:   func(app.IAPPConf) error { return nil },
			closing:    func(app.IAPPConf) error { return nil },
			wantLog:    "启动成功(ws,ws:",
		},
		{name: "启动web服务成功",
			cnf:        confObj.GetWebConf(),
			serverType: "web",
			serverName: "webserver",
			starting:   func(app.IAPPConf) error { return nil },
			closing:    func(app.IAPPConf) error { return nil },
			wantLog:    "启动成功(web,http:",
		},
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

		//占用端口使服务启动失败
		var httpServer *testServer
		if tt.isServerStart {
			httpServer = &testServer{addr: "127.0.0.1:55004"}
			go httpServer.testListen()
		}

		//创建节点使服务发布报错
		if tt.isServerPub {
			newConfObj := mocks.NewConfBy("hydra", "test", "fs://./") //构建对象
			newConfObj.API(":55004")
			tt.cnf = newConfObj.GetAPIConf()
			path := fmt.Sprintf("./hydra/%s/%s/test/servers", tt.serverName, tt.serverType)
			os.RemoveAll(path) //删除文件夹
			os.Create(path)    //使文件夹节点变成文件节点,让该节点下不能创建文件
		}

		//构建服务器
		rsp, _ := http.NewResponsive(tt.cnf)

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

		//释放端口
		if tt.isServerStart {
			httpServer.close()
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

		//fmt.Println("xxxx:", string(out))
		if tt.wantLog != "" {
			assert.Equalf(t, true, strings.Contains(string(out), tt.wantLog), tt.name+"log")
		}
	}
}
