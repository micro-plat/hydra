package servers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	queueredis "github.com/micro-plat/hydra/conf/vars/queue/redis"
	varredis "github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/registry"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

// func TestRspServers_delayPub(t *testing.T) {

// 	tests := []struct {
// 		name string
// 		p    string
// 		done bool
// 		want string
// 	}{
// 		{name: "延迟发布过程中 服务未关闭", p: "/path", want: "/path"},
// 		{name: "延迟发布过程中 服务已关闭", p: "/path", done: true},
// 	}
// 	r := servers.NewRspServers("lm://./", "hydra_test", "servers", []string{"api"}, "t")
// 	for _, tt := range tests {
// 		r.delayPub(tt.p)
// 		time.Sleep(time.Second * 1)
// 		r.done = tt.done
// 		time.Sleep(time.Second * 2)
// 		if r.done {
// 			assert.Equal(t, 0, len(r.delayChan), tt.name)
// 			continue
// 		}
// 		select {
// 		case v := <-r.delayChan:
// 			assert.Equal(t, tt.want, v, tt.name)
// 		default:
// 			t.Error("测试未通过", tt.name)
// 		}
// 	}
// }

func TestRspServers_Start(t *testing.T) {

	tests := []struct {
		name       string
		serverName string
		sysType    string
		isFirst    bool
		wantErr    bool
	}{
		{name: "启动cronServer", serverName: "cronserver", sysType: "cron", isFirst: true},
		{name: "启动apiServer", serverName: "apiserver", sysType: "api"},
		{name: "启动mqcServer", serverName: "mqcserver", sysType: "mqc"},
		//	{name: "启动rpcServer", serverName: "rpcserver", sysType: "rpc"},
	}

	platName := "servershydra_test"
	clusterName := "serv_test_go"
	registryAddr := "lm://./"

	for _, tt := range tests {

		//构建的新的os.Stdout
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		rsp := servers.NewRspServers(registryAddr, platName, tt.serverName, []string{tt.sysType}, clusterName)

		err := rsp.Start()
		assert.Equalf(t, tt.wantErr, err != nil, tt.name)
		time.Sleep(time.Second * 1)

		//初始化注册中心
		sc := mocks.NewConfBy(platName, clusterName)
		sc.API(":50001")
		sc.GetCronConf()
		sc.GetAPIConf()
		sc.Vars().Redis("5.79", varredis.New([]string{"192.168.5.79:6379"}))
		sc.Vars().Queue().Redis("xxx", queueredis.New(queueredis.WithConfigName("5.79")))
		sc.MQC("redis://xxx")
		sc.GetMQCConf()
		time.Sleep(time.Second * 1)

		//注册中心节点值发生变化
		path := fmt.Sprintf("/servershydra_test/%s/%s/serv_test_go/conf", tt.serverName, tt.sysType)
		registry, err := registry.NewRegistry(registryAddr, logger.New("hydra"))
		assert.Equalf(t, false, err != nil, tt.name)
		err = registry.Update(path, `{"status":"start","addr":"redis://xxx"}`)
		assert.Equalf(t, false, err != nil, tt.name)
		time.Sleep(time.Second * 1)

		rsp.Shutdown()
		time.Sleep(time.Second * 1)

		//获取输出
		w.Close()
		out, err := ioutil.ReadAll(r)
		assert.Equalf(t, false, err != nil, tt.name)

		//还原os.Stdout
		os.Stdout = rescueStdout

		//	fmt.Println("out:", string(out))

		wantLog := fmt.Sprintf("初始化: %s", path)
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"初始化")
		if tt.isFirst {
			wantLog = fmt.Sprintf("监听服务器配置...")
			assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"监听")
		}
		wantLog = fmt.Sprintf("启动[%s]服务...", tt.sysType)
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"启动")
		wantLog = fmt.Sprintf("启动成功")
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"成功")
		wantLog = fmt.Sprintf("配置发生变化%s", path)
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"配置变化")
		wantLog = fmt.Sprintf("配置更新完成")
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"配置更新完成")
		wantLog = fmt.Sprintf("关闭[%s]服务...", tt.sysType)
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"关闭")
	}
}

type testServer struct {
	addr string
	s    *http.Server
}

func (t *testServer) testListen() {
	t.s = &http.Server{Addr: t.addr, Handler: nil}
	t.s.ListenAndServe()
}

func (t *testServer) close() {
	t.s.Close()
}

func TestRspServers_Start_ServerStartErr(t *testing.T) {

	tests := []struct {
		name       string
		serverName string
		sysType    string
		isFirst    bool
		wantErr    bool
	}{
		{name: "启动apiServer失败,之后延迟启动成功", serverName: "apiserver", sysType: "api"},
	}

	platName := "servershydra_test1"
	clusterName := "serv_test_go1"
	registryAddr := "lm://./"

	for _, tt := range tests {

		//构建的新的os.Stdout
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//初始化注册中心
		sc := mocks.NewConfBy(platName, clusterName)
		sc.API(":50002")
		sc.GetAPIConf()
		time.Sleep(time.Second * 2)
		//占用端口
		s := &testServer{addr: "127.0.0.1:50002"}
		go s.testListen()

		rsp := servers.NewRspServers(registryAddr, platName, tt.serverName, []string{tt.sysType}, clusterName)
		rsp.DelayTime = time.Second
		err := rsp.Start()
		assert.Equalf(t, tt.wantErr, err != nil, tt.name)
		time.Sleep(time.Second * 2)

		//释放端口
		s.close()
		time.Sleep(time.Second * 1)

		rsp.Shutdown()
		time.Sleep(time.Second * 1)

		//获取输出
		w.Close()
		out, err := ioutil.ReadAll(r)
		assert.Equalf(t, false, err != nil, tt.name)

		//还原os.Stdout
		os.Stdout = rescueStdout

		// fmt.Println("out:", string(out))
		path := fmt.Sprintf("/servershydra_test1/%s/%s/serv_test_go1/conf", tt.serverName, tt.sysType)
		wantLog := fmt.Sprintf("初始化: %s", path)
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"初始化")
		wantLog = fmt.Sprintf("启动[%s]服务...", tt.sysType)
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"启动")
		wantLog = fmt.Sprintf("[%s]服务器启动失败:%s启动失败 listen tcp 0.0.0.0:50002: bind: address already in use", tt.sysType, tt.sysType)
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"占用")
		wantLog = fmt.Sprintf("监听服务器配置...")
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"监听")
		wantLog = fmt.Sprintf("启动成功")
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"成功")
		wantLog = fmt.Sprintf("关闭[%s]服务...", tt.sysType)
		assert.Equalf(t, true, strings.Contains(string(out), wantLog), tt.name+"关闭")
	}
}
