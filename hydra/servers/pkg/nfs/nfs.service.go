package nfs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/micro-plat/hydra/conf/server/auth"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
	"github.com/micro-plat/lib4go/errs"
)

//GetDirList 获取本机目录信息
func (c *cnfs) GetDirList(ctx context.IContext) interface{} {
	return c.infs.GetDirList(infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME,
		ctx.Request().GetString(infs.DIRNAME))),
		ctx.Request().GetInt("deep", 1))
}

//Upload 用户上传文件
func (c *cnfs) Upload(ctx context.IContext) interface{} {
	//读取文件
	name := ctx.Request().GetString(infs.FILENAME, "file")
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
	path := infs.MultiPath(ctx.Request().Path().Params().GetString("path", ctx.Request().GetString("path")))
	path, domain, err := c.infs.Save(path, name, buff)
	if err != nil {
		return err
	}

	// 处理返回结果
	ctx.Response().AddSpecial(fmt.Sprintf("nfs|%s|%d", name, size))
	return map[string]interface{}{
		"path": fmt.Sprintf("%s/%s", strings.Trim(domain, "/"), strings.Trim(path, "/")),
	}
}

//Download 用户下载文件
func (c *cnfs) Download(ctx context.IContext) interface{} {

	//检查参数
	dir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME))
	name := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.FILENAME))
	if name == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", infs.FILENAME)
	}

	//获取文件

	path := filepath.Join(dir, name)
	buff, tp, err := c.infs.Get(path)
	if err != nil {
		return err
	}

	//写入文件
	//未设置文件头
	ctx.Response().ContentType(tp)
	ctx.Response().GetHTTPReponse().Write(buff)
	return nil
}

//GetFileList 获取本机的指定文件的指纹信息，仅master提供对外查询功能
func (c *cnfs) GetFileList(ctx context.IContext) interface{} {
	return c.infs.GetFileList(infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME,
		ctx.Request().GetString(infs.DIRNAME))),
		ctx.Request().GetString("kw"),
		ctx.Request().GetBool("all", false),
		ctx.Request().GetInt("pi", 0),
		ctx.Request().GetInt("ps", 100))
}

//CreateDir 创建目录
func (c *cnfs) CreateDir(ctx context.IContext) interface{} {
	//检查参数
	dir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME, ctx.Request().GetString(infs.DIRNAME)))
	if dir == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", infs.DIRNAME)
	}
	return c.infs.CreateDir(c.c.Local, dir)
}

//RenameDir 重命名目录
func (c *cnfs) RenameDir(ctx context.IContext) interface{} {
	//检查参数
	dir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME, ctx.Request().GetString(infs.DIRNAME)))
	ndir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME, ctx.Request().GetString(infs.DIRNAME)))
	if dir == "" || ndir == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\",\":%s\"", infs.DIRNAME, infs.DIRNAME)
	}
	return c.infs.Rename(c.c.Local, dir, ndir)
}

//ImgScale 缩略图生成
func (c *cnfs) ImgScale(ctx context.IContext) interface{} {
	//检查参数
	dir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME, ctx.Request().GetString(infs.DIRNAME)))
	name := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.FILENAME, ctx.Request().GetString(infs.FILENAME)))
	if name == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", infs.FILENAME)
	}

	//获取文件
	path := filepath.Join(dir, name)
	width := ctx.Request().GetInt("w")
	height := ctx.Request().GetInt("h")
	quality := ctx.Request().GetInt("q")
	buff, ctp, err := c.infs.GetScaleImage(c.c.Local, path, width, height, quality)
	if err == nil {
		ctx.Response().ContentType(ctp)
		ctx.Response().GetHTTPReponse().Write(buff)
		return nil
	}
	return err
}

//View 获取PDF预览文件
func (c *cnfs) GetPDF4Preview(ctx context.IContext) interface{} {
	//检查参数
	dir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME, ctx.Request().GetString(infs.DIRNAME)))
	name := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.FILENAME, ctx.Request().GetString(infs.FILENAME)))
	if name == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", infs.FILENAME)
	}

	//获取文件
	path := filepath.Join(dir, name)

	buff, contentType, err := c.infs.Conver2PDF(c.c.Local, path)
	if err != nil {
		return err
	}
	ctx.Response().ContentType(contentType)
	ctx.Response().GetHTTPReponse().Write(buff)
	return nil
}

func init() {
	auth.AppendExcludes(infs.NOTEXCLUDES...)
}
