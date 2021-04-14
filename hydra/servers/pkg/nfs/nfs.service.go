package nfs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/errs"
)

//SVSNOTExcludes 服务路径中排除的
var SVSNOTExcludes = []string{"/nfs/**", "/_/nfs/**"}

const (
	//SVSUpload 用户端上传文件
	SVSUpload = "/nfs/upload"

	//SVSDonwload 用户端下载文件
	SVSDonwload = "/nfs/download/:dir/:name"

	fileName = "name"

	dirName = "dir"

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
	fp, err := c.module.SaveNewFile(name, buff)
	if err != nil {
		return err
	}

	// 处理返回结果
	ctx.Response().AddSpecial(fmt.Sprintf("nfs|%s|%d", name, size))
	return map[string]interface{}{
		"name": fmt.Sprintf("%s/%s", strings.Trim(c.c.Domain, "/"), strings.Trim(fp.Path, "/")),
	}
}

//Download 用户下载文件
func (c *cnfs) Download(ctx context.IContext) interface{} {

	//检查参数
	dir := ctx.Request().Path().Params().GetString(dirName)
	name := ctx.Request().Path().Params().GetString(fileName)
	if dir == "" || name == "" {
		return errs.NewError(http.StatusNotAcceptable, "参数不能为空")
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
