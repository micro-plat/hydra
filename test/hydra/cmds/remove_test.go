package cmds

import (
	"fmt"
	"io/ioutil"
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
	//正常的删除
	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	//defer os.Remove(fileName)

	//1. 先安装服务
	args := []string{"xxtest", "install", "-r", removeRegistryAddr, "-c", "c"}
	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	os.Args = args
	app.Start()
	time.Sleep(time.Second * 2)

	//2. 删除服务
	args = []string{"xxtest", "remove"}
	os.Args = args
	app.Start()
	time.Sleep(time.Second * 2)

	//还原std
	*os.Stdout = orgStd

	file.Close()
	time.Sleep(time.Second)
	bytes, err := ioutil.ReadFile(fileName)

	fmt.Println("bytes:", string(bytes), err)

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
