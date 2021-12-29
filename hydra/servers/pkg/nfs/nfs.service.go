package nfs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/micro-plat/hydra/conf/server/auth"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/internal"
	"github.com/micro-plat/lib4go/errs"
)

//notExcludes 服务路径中排除的
var notExcludes = []string{"/**/_/nfs/**"}

const (
	fileName = "name"

	dirName = "dir"

	ndirName = "ndir"
)

const (
	//SVSUpload 用户端上传文件
	SVSUpload = "/nfs/upload"

	SVSPreview = "/nfs/preview"

	//SVSDonwload 用户端下载文件
	SVSDonwload = "/nfs/file/:dir/:name"

	//SVSList 文件列表
	SVSList = "/nfs/file/list"

	//SVSDir 目录列表
	SVSDir = "/nfs/dir/list"

	//SVSScalrImage 压缩文件
	SVSScalrImage = "/nfs/scale/:name"

	//SVSCreateDir 创建文件目录
	SVSCreateDir = "/nfs/create/:dir"

	//SVSRenameDir 重命名文件目录
	SVSRenameDir = "/nfs/create/:dir/:ndir"

	//获取远程文件的指纹信息
	rmt_fp_get = "/_/nfs/fp/get"

	//推送指纹数据
	rmt_fp_notify = "/_/nfs/fp/notify"

	//拉取指纹列表
	rmt_fp_query = "/_/nfs/fp/query"

	//获取远程文件数据
	rmt_file_download = "/_/nfs/file/download"
)

//Query 获取每个机器所有文件
func (c *cnfs) Query(ctx context.IContext) interface{} {
	list := c.module.Query()
	ctx.Response().AddSpecial("nfs")
	ctx.Response().AddSpecial(fmt.Sprintf("%d", len(list)))
	return list
}

//GetFP 获取本机的指定文件的指纹信息，仅master提供对外查询功能
func (c *cnfs) GetFP(ctx context.IContext) interface{} {
	if err := ctx.Request().Check(fileName); err != nil {
		return err
	}
	fp, err := c.module.GetFP(ctx.Request().GetString(fileName))
	if err != nil {
		return err
	}
	return fp
}

//GetFileList 获取本机的指定文件的指纹信息，仅master提供对外查询功能
func (c *cnfs) GetFileList(ctx context.IContext) interface{} {
	return c.module.GetFileList(multiPath(ctx.Request().GetString(dirName)),
		ctx.Request().GetString("kw"),
		ctx.Request().GetBool("all"),
		ctx.Request().GetInt("pi"),
		ctx.Request().GetInt("ps", 100))
}

//GetDirList 获取本机目录信息
func (c *cnfs) GetDirList(ctx context.IContext) interface{} {
	return c.module.GetDirList(multiPath(ctx.Request().GetString(dirName)),
		ctx.Request().GetInt("deep", 1))
}

//RecvNotify 接收远程文件通知
func (c *cnfs) RecvNotify(ctx context.IContext) interface{} {
	fp := make(eFileFPLists)
	if err := ctx.Request().ToStruct(&fp); err != nil {
		return err
	}
	if err := c.module.RecvNotify(fp); err != nil {
		return err
	}
	return "success"
}

//Download 用户下载文件
func (c *cnfs) GetFile(ctx context.IContext) interface{} {
	//检查输入参数
	if err := ctx.Request().Check(fileName); err != nil {
		return err
	}

	//从本地获取文件
	path := ctx.Request().GetString(fileName)
	err := c.module.HasFile(path)
	if err != nil {
		return err
	}
	ctx.Response().File(path, http.FS(c.module.local))
	return nil
}

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

	// 保存文件
	path := multiPath(ctx.Request().Path().Params().GetString("path", ctx.Request().GetString("path")))
	fp, domain, err := c.module.SaveNewFile(path, name, buff)
	if err != nil {
		return err
	}

	// 处理返回结果
	ctx.Response().AddSpecial(fmt.Sprintf("nfs|%s|%d", name, size))
	return map[string]interface{}{
		"path": fmt.Sprintf("%s/%s", strings.Trim(domain, "/"), strings.Trim(fp.Path, "/")),
	}
}

//Download 用户下载文件
func (c *cnfs) Download(ctx context.IContext) interface{} {

	//检查参数
	dir := multiPath(ctx.Request().Path().Params().GetString(dirName))
	name := multiPath(ctx.Request().Path().Params().GetString(fileName))
	if name == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", fileName)
	}

	//获取文件

	path := filepath.Join(dir, name)
	err := c.module.checkAndDownload(path)
	if err != nil {
		return err
	}

	//写入文件
	ctx.Response().File(path, http.FS(c.module.local))
	return nil
}

//CreateDir 创建目录
func (c *cnfs) CreateDir(ctx context.IContext) interface{} {
	//检查参数
	dir := multiPath(ctx.Request().Path().Params().GetString(dirName, ctx.Request().GetString(dirName)))
	if dir == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", dirName)
	}
	return internal.CreateDir(c.c.Local, dir)
}

//RenameDir 重命名目录
func (c *cnfs) RenameDir(ctx context.IContext) interface{} {
	//检查参数
	dir := multiPath(ctx.Request().Path().Params().GetString(dirName, ctx.Request().GetString(dirName)))
	ndir := multiPath(ctx.Request().Path().Params().GetString(ndirName, ctx.Request().GetString(ndirName)))
	if dir == "" || ndir == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\",\":%s\"", dirName, ndirName)
	}
	return internal.Rename(c.c.Local, dir, ndir)
}

//ImgScale 缩略图生成
func (c *cnfs) ImgScale(ctx context.IContext) interface{} {
	//检查参数
	dir := multiPath(ctx.Request().Path().Params().GetString(dirName, ctx.Request().GetString(dirName)))
	name := multiPath(ctx.Request().Path().Params().GetString(fileName, ctx.Request().GetString(fileName)))
	if name == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", fileName)
	}

	//获取文件
	path := filepath.Join(dir, name)
	width := ctx.Request().GetInt("w")
	height := ctx.Request().GetInt("h")
	quality := ctx.Request().GetInt("q")
	buff, err := internal.ScaleImageByPath(c.c.Local, path, width, height, quality)
	if err == nil {
		ctx.Response().ContentType(internal.GetContentType(path))
		ctx.Response().GetHTTPReponse().Write(buff)
		return nil
	}

	ctx.Log().Error(fmt.Errorf("%w %s", err, path))
	buff, err = internal.ReadFile(filepath.Join(c.c.Local, path))
	if err != nil {
		return err
	}
	ctx.Response().ContentType(internal.GetContentType(path))
	ctx.Response().GetHTTPReponse().Write(buff)
	return nil
}

//View 获取PDF预览文件
func (c *cnfs) GetPDF4Preview(ctx context.IContext) interface{} {
	//检查参数
	dir := multiPath(ctx.Request().Path().Params().GetString(dirName, ctx.Request().GetString(dirName)))
	name := multiPath(ctx.Request().Path().Params().GetString(fileName, ctx.Request().GetString(fileName)))
	if name == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", fileName)
	}

	//获取文件
	path := filepath.Join(dir, name)
	contentType, buff, err := internal.Conver2PDF(c.c.Local, path)
	if err != nil {
		return err
	}
	ctx.Response().ContentType(contentType)
	ctx.Response().GetHTTPReponse().Write(buff)
	return nil
}

func init() {
	auth.AppendExcludes(notExcludes...)
}

//处理多级目录
func multiPath(path string) string {
	return strings.Trim(strings.ReplaceAll(path, "|", "/"), "/")
}
