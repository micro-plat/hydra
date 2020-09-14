package test

import (
	"io/ioutil"
	xhttp "net/http"
	"os"
	"testing"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/lib4go/logger"
)

func getServer() *hydra.MicroApp {
	//设置服务器参数
	raw := `	
	
	local ips = {}
	local upstream = ""
	
	
	function getUpStream()
		return upstream;
	end
		
	function go2UpStream() 			
		local req = require("request")
		return req.getClientIP()
	end`

	hydra.Conf.API(":8081").Gray(raw)

	app := hydra.NewApp(hydra.WithServerTypes(http.API),
		hydra.WithUsage("apiserver"),
		hydra.WithRegistry("lm://."),
		hydra.WithPlatName("test"),
		hydra.WithClusterName("t"),
		hydra.WithSystemName("apiserver"),
	)
	app.Micro("/hello", func(hydra.IContext) interface{} {
		return "SUCCESS"
	})
	return app

}
func TestGray(t *testing.T) {

	os.Args = []string{"test", "run"}
	app := getServer()
	go app.Start()
	time.Sleep(time.Second * 1)

	resp, err := xhttp.Get("http://localhost:8081/hello")
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 || string(bodyBytes) != "SUCCESS" {
		t.Error(resp.Status, string(bodyBytes))
	}
}

func BenchmarkServerGray(t *testing.B) {

	os.Args = []string{"test", "run"}
	logger.Pause()
	app := getServer()
	go app.Start()
	defer app.Close()

	time.Sleep(time.Second * 1)

	request := func() {
		resp, err := xhttp.Get("http://localhost:8081/hello")
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 || string(bodyBytes) != "SUCCESS" {
			t.Error("ERR:", resp.Status, string(bodyBytes))
		}
	}
	t.ResetTimer()
	t.N = 10000
	for i := 0; i < t.N; i++ {
		request()
	}
}
