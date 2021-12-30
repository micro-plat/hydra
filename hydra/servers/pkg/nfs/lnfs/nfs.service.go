package lnfs

import (
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
)

//Query 获取每个机器所有文件
func (c *LNFS) Query(ctx context.IContext) interface{} {
	list := c.Module.Query()
	ctx.Response().AddSpecial("nfs")
	ctx.Response().AddSpecial(fmt.Sprintf("%d", len(list)))
	return list
}

//GetFP 获取本机的指定文件的指纹信息，仅master提供对外查询功能
func (c *LNFS) GetFP(ctx context.IContext) interface{} {
	if err := ctx.Request().Check(infs.FILENAME); err != nil {
		return err
	}
	fp, err := c.Module.GetFP(ctx.Request().GetString(infs.FILENAME))
	if err != nil {
		return err
	}
	return fp
}

//Download 用户下载文件
func (c *LNFS) GetFile(ctx context.IContext) interface{} {
	//检查输入参数
	if err := ctx.Request().Check(infs.FILENAME); err != nil {
		return err
	}

	//从本地获取文件
	path := ctx.Request().GetString(infs.FILENAME)
	err := c.Module.HasFile(path)
	if err != nil {
		return err
	}
	ctx.Response().File(path, http.FS(c.Module.Local))
	return nil
}

//RecvNotify 接收远程文件通知
func (c *LNFS) RecvNotify(ctx context.IContext) interface{} {
	fp := make(EFileFPLists)
	if err := ctx.Request().ToStruct(&fp); err != nil {
		return err
	}
	if err := c.Module.RecvNotify(fp); err != nil {
		return err
	}
	return "success"
}

var _ infs.Infs = &LNFS{}
