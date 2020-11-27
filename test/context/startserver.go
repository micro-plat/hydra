package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	ghttp "net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var oncelock sync.Once

func startServer() {
	oncelock.Do(func() {

		app := hydra.NewApp(
			hydra.WithPlatName("hydra_context"),
			hydra.WithSystemName("test"),
			hydra.WithServerTypes(http.API),
			hydra.WithDebug(),
			hydra.WithClusterName("t"),
			hydra.WithRegistry("lm://."),
		)

		hydra.Conf.API(":9091")
		app.API("/getbodymap", GetBodyMap)
		app.API("/getbody/encoding", GetBodyEncoding)
		app.API("/getbody/encoding/gbk", GetBodyEncodingGBK, api.WithEncoding("gbk"))
		app.API("/getbody/encoding/utf8", GetBodyEncodingUTF8, api.WithEncoding("utf-8"))
		app.API("/getcookies/encoding", GetCookiesEncoding)
		app.API("/getheaders/encoding", GetHeaderEncoding)
		app.API("/getheaders/encoding/gbk", GetHeaderEncodingGBK, api.WithEncoding("gbk"))
		app.API("/getheaders/encoding/utf8", GetHeaderEncoding, api.WithEncoding("utf-8"))
		app.API("/form", GetBodyMapFormData)
		app.API("/upload", UploadFile)
		app.API("/download", DownloadFile)
		app.API("/request/bind", Bind)
		app.API("/request/check", Check)
		app.API("/response/redirect", Redirect)
		app.API("/response/redirect/dst", RedirectDst)

		os.Args = []string{"startserver", "run"}
		go app.Start()
		time.Sleep(time.Second * 2)
	})

}

func GetBodyMap(ctx hydra.IContext) interface{} {
	raw, err := ctx.Request().GetBodyMap()
	if err != nil {
		return fmt.Errorf("getBody出错")
	}
	header := ctx.Request().GetHeader("Content-Type")
	ctx.Response().Header("Content-Type", header)
	return raw
}

func GetBodyEncoding(ctx hydra.IContext) interface{} {
	raw, err := ctx.Request().GetBody()
	if err != nil {
		return fmt.Errorf("getBody出错")
	}
	return raw
}

func GetBodyEncodingGBK(ctx hydra.IContext) interface{} {
	raw, err := ctx.Request().GetBody()
	if err != nil {
		return fmt.Errorf("getBody出错")
	}
	return raw
}
func GetBodyEncodingUTF8(ctx hydra.IContext) interface{} {
	raw, err := ctx.Request().GetBody()
	if err != nil {
		return fmt.Errorf("getBody出错")
	}
	return raw
}

func GetCookiesEncoding(ctx hydra.IContext) interface{} {
	r := ctx.Request().GetCookies()
	if r == nil {
		return fmt.Errorf("GetCookies出错")
	}
	return r
}

func GetHeaderEncoding(ctx hydra.IContext) interface{} {
	r := ctx.Request().GetHeader("Hname")
	if r == "" {
		return fmt.Errorf("GetHeaders出错")
	}
	return r
}

func GetHeaderEncodingGBK(ctx hydra.IContext) interface{} {
	r := ctx.Request().GetHeader("Hname")
	if r == "" {
		return fmt.Errorf("GetHeaders出错")
	}
	return r
}

func GetHeaderEncodingUtf8(ctx hydra.IContext) interface{} {
	r := ctx.Request().GetHeader("Hname")
	if r == "" {
		return fmt.Errorf("GetHeaders出错")
	}
	e := ctx.Request().Path().GetEncoding()
	if e == "gbk" {
		r = GbkToUtf8(r)
	}
	return r
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

func Bind(ctx hydra.IContext) interface{} {
	type result struct {
		Key   string `json:"key" valid:"required" yaml:"key" xml:"key" form:"key"`
		Value string `json:"value" valid:"required" yaml:"value"  xml:"value" form:"value"`
	}
	s := &result{}
	err := ctx.Request().Bind(s)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Bind出错 %+v", err)
	}
	return *s
}

func Check(ctx hydra.IContext) interface{} {
	err := ctx.Request().Check("key", "value")
	if err != nil {
		fmt.Println("err:", err)
		return fmt.Errorf("Check出错 %+v", err)
	}
	return "success"
}

func Redirect(ctx hydra.IContext) interface{} {
	ctx.Response().Header("Location", "http://localhost:9091/response/redirect/dst")
	ctx.Response().Abort(ghttp.StatusFound)
	return nil
}
func RedirectDst(ctx hydra.IContext) interface{} {
	ctx.Response().Header("Content-Type", "application/json; charset=UTF-8")
	return "success"
}

func getTestUTF8Json(s map[string]string) string {
	for k, v := range s {
		s[k] = url.QueryEscape(v)
	}
	r, _ := json.Marshal(s)
	return string(r)
}

func getTestGBKJson(s map[string]string) string {
	for k, v := range s {
		s[k] = url.QueryEscape(Utf8ToGbk(v))
	}

	r, _ := json.Marshal(s)
	return string(r)
}

func Utf8ToGbk(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewEncoder())
	d, _ := ioutil.ReadAll(reader)
	return string(d)
}

func GbkToUtf8(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewDecoder())
	d, _ := ioutil.ReadAll(reader)
	return string(d)
}

func getTestMIMEMultipartPOSTForm() string {
	file, _ := os.Open("upload.test.txt")
	defer file.Close()
	body := &bytes.Buffer{}
	// 文件写入 body
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("upload", filepath.Base("upload.test.txt"))
	io.Copy(part, file)
	writer.Close()
	return body.String()
}


func getUploadBody() string {
	return "Content-Disposition: form-data; name=\"upload\"; filename=\"upload.test.txt\"\r\nContent-Type: application/octet-stream\r\n\r\nADASDASDASFHNOJM~!@#$%^&*"
}
