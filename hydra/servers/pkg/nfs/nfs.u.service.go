package nfs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/errs"
)

//SVSNOTExcludes 服务路径中排除的
var SVSNOTExcludes = []string{"/nfs/file/*", "/_/nfs/**"}

const (
	//SVSUpload 用户端上传文件
	SVSUpload = "/nfs/file/upload"

	//SVSDonwload 用户端下载文件
	SVSDonwload = "/nfs/file/download/:dir/:name"
)

//Upload 用户上传文件
func (c *cnfs) Upload(ctx context.IContext) interface{} {

	//读取文件
	name := ctx.Request().GetString(fileName, "file")
	name, reader, size, err := ctx.Request().GetFile(name)
	if err != nil {
		return err
	}

	//读取内容
	defer reader.Close()
	buff, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	//写入文件
	fp, err := c.module.SaveNewFile(name, buff)
	if err != nil {
		return err
	}
	ctx.Response().AddSpecial(fmt.Sprintf("nfs|%s|%d", name, size))
	return map[string]interface{}{
		"name": fp.Path,
	}
}

//Download 用户下载文件
func (c *cnfs) Download(ctx context.IContext) interface{} {
	dir := ctx.Request().Path().Params().GetString(dirName)
	name := ctx.Request().Path().Params().GetString(fileName)
	if dir == "" || name == "" {
		return errs.NewError(http.StatusNotAcceptable, "参数不能为空")
	}
	path := filepath.Join(dir, name)
	fmt.Println("path:", path)
	_, err := c.module.GetFile(path)
	if err != nil {
		return err
	}

	//写入文件
	ctx.Response().File(path, http.FS(c.module.local))
	return nil
}
