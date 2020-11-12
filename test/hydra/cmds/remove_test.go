package cmds

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/test/assert"
)

const removeRegistryAddr = "lm://."

func Test_remove_Normal(t *testing.T) {
	defunc, fileCallback := injectStdOutFile()
	defer defunc()

	args := []string{"xxtest", "install", "-r", removeRegistryAddr, "-c", "c"}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)

	//1. 清除服务(保证没有服务安装)
	os.Args = []string{"xxtest", "remove"}
	go app.Start()
	time.Sleep(time.Second * 2)

	//2. 先安装服务
	os.Args = args
	app.Start()
	time.Sleep(time.Second * 2)

	//3. 删除服务
	args = []string{"xxtest", "remove"}
	os.Args = args
	app.Start()

	time.Sleep(time.Second * 2)
	bytes, err := fileCallback()

	if err != nil {
		t.Error(err)
		return
	}
	lines := strings.Split(string(bytes), "\r")
	for _, row := range lines {
		//找到响应数据行
		if !(strings.Contains(row, "Install") || strings.Contains(row, "Removing")) {
			continue
		}
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			//unbuntu/centos
			result := strings.Contains(row, "sudo") || strings.Contains(row, "OK")
			assert.Equal(t, true, result, "正常安装再删除")
		}
		if runtime.GOOS == "windows" {
			result := strings.Contains(row, "OK")
			assert.Equal(t, true, result, "正常安装再删除")
		}
		return
	}
	time.Sleep(time.Second)
}

func Test_remove_NotExists(t *testing.T) {
	defunc, fileCallback := injectStdOutFile()
	defer defunc()

	//1. 删除服务
	args := []string{"xxtest", "remove"}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	os.Args = args
	app.Start()

	time.Sleep(time.Second * 2)
	bytes, err := fileCallback()
	if err != nil {
		t.Error("读取文件报错：", err)
	}

	lines := strings.Split(string(bytes), "\r")
	for _, row := range lines {
		//找到响应数据行
		if !(strings.Contains(row, "not")) {
			continue
		}
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			//unbuntu/centos
			result := strings.Contains(row, "sudo") || strings.Contains(row, "Service is not installed")
			assert.Equal(t, true, result, "删除不存在的服务")
		}
		if runtime.GOOS == "windows" {
			result := strings.Contains(row, "service does not exist as an installed service")
			assert.Equal(t, true, result, "删除不存在的服务")
		}
		return
	}
	time.Sleep(time.Second)
}
