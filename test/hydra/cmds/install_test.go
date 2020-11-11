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

const installRegistryAddr = "lm://."

func Test_install_Normal(t *testing.T) {

	fileName := fmt.Sprint(time.Now().Nanosecond())
	file, _ := os.Create(fileName)
	orgStd := *os.Stdout
	*os.Stdout = *file
	//defer os.Remove(fileName)
	//正常的安装
	args := []string{"xxtest", "install", "-r", installRegistryAddr, "-c", "c"}

	var app = hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithPlatName("xxtest"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("c"),
	)
	os.Args = args
	app.Start()
	time.Sleep(time.Second)

	//还原std
	*os.Stdout = orgStd

	file.Close()
	time.Sleep(time.Second)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		orgStd.WriteString(err.Error())
	}
	lines := strings.Split(string(bytes), "\r")
	for _, row := range lines {
		if runtime.GOOS == "linux" {
			//unbuntu/centos
			result := strings.Contains(row, "sudo") || strings.Contains(row, "OK")
			assert.Equal(t, true, result, "正常参数的安装")
		}
		if runtime.GOOS == "windows" {

			assert.Equal(t, true, result, "正常参数的安装")
		}
		return
	}
}
