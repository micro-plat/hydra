package context

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

var starting bool

func startServer() {
	if starting {
		return
	}
	starting = true
	app := hydra.NewApp(
		hydra.WithPlatName("test"),
		hydra.WithSystemName("test"),
		hydra.WithServerTypes(http.API),
		hydra.WithDebug(),
		hydra.WithClusterName("t"),
		hydra.WithRegistry("lm://."),
	)

	hydra.Conf.API(":9091")
	app.API("/getbodymap", GetBodyMap)
	app.API("/form", GetBodyMapFormData)
	app.API("/upload", UploadFile)
	app.API("/download", DownloadFile)
	os.Args = []string{"startserver", "run"}
	go app.Start()
	time.Sleep(time.Second * 2)
}

func GetBodyMap(ctx hydra.IContext) interface{} {
	raw, err := ctx.Request().GetRawBodyMap()
	if err != nil {
		return fmt.Errorf("getBody出错")
	}
	header := ctx.Request().Path().GetHeader("Content-Type")
	ctx.Response().Header("Content-Type", header)
	return raw
}

func GetBodyMapFormData(ctx hydra.IContext) interface{} {
	raw, err := ctx.Request().GetBody()
	if err != nil {
		return fmt.Errorf("getBody出错")
	}
	ctx.Response().Header("Content-Type", "text/plain")
	return raw
}

func UploadFile(ctx hydra.IContext) interface{} {
	fileName, err := ctx.Request().GetFileName("upload")
	if err != nil {
		return fmt.Errorf("UploadFile出错1")
	}
	size, err := ctx.Request().GetFileSize("upload")
	if err != nil {
		return fmt.Errorf("UploadFile出错2")
	}
	// err = ctx.Request().SaveFile("upload", "save.txt")
	// if err != nil {
	// 	return fmt.Errorf("UploadFile出错s")
	// }
	body, err := ctx.Request().GetFileBody("upload")
	if err != nil {
		return fmt.Errorf("UploadFile出错3")
	}
	ctx.Log().Info(fileName, size)
	s, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("UploadFile出错4")
	}
	ctx.Response().Header("Content-Type", "application/json")
	return map[string]interface{}{
		"fileName": fileName,
		"size":     size,
		"body":     string(s),
	}
}

func DownloadFile(ctx hydra.IContext) (r interface{}) {
	var buffer bytes.Buffer
	f, err := os.Open("upload.test.txt")
	if err != nil {
		return err
	}
	buffer.ReadFrom(f)
	ctx.Log().Info("设置响应头")
	ctx.Response().Header("Accept-Ranges", "bytes")
	ctx.Response().Header("Content-Type", "text/plain")
	ctx.Response().Header("Content-Disposition", "attachment;filename=hello.txt")
	ctx.Response().Header("Charset", "utf-8")

	ctx.Log().Info("返回数据")
	return buffer.String()
}
