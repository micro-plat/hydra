package nfs

import (
	"net/http"

	"github.com/micro-plat/hydra/context"
)

//SVS_NOT_Excludes 服务路径中排除的
var SVS_NOT_Excludes = []string{"/nfs/file/*", "/_/nfs/**"}

const (
	SVS_Upload   = "/nfs/file/upload"
	SVS_Donwload = "/nfs/file/download"
)

//Upload 用户上传文件
func (c *cnfs) Upload(ctx context.IContext) interface{} {

	//读取文件
	defer req(ctx).rspns()
	name, reader, size, err := ctx.Request().GetFile("file")
	if err != nil {
		return err
	}

	//读取内容
	defer reader.Close()
	buff := make([]byte, 0, size)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}

	//写入文件
	fp, err := c.module.SaveNewFile(name, buff)
	if err != nil {
		return err
	}
	return map[string]interface{}{
		"name": fp.Path,
	}
}

//Download 用户下载文件
func (c *cnfs) Download(ctx context.IContext) interface{} {
	//根据路径查询文件
	defer req(ctx).rspns()
	path := ctx.Request().Path().GetURL().Path
	_, err := c.module.GetFile(path)
	if err != nil {
		return err
	}

	//写入文件
	ctx.Response().File(path, http.FS(c.module.local))
	return nil
}
