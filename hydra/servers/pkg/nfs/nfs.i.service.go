package nfs

import (
	"net/http"

	"github.com/micro-plat/hydra/context"
)

//GetFPList 获取每个机器所有文件
func (c *cnfs) GetFPList(ctx context.IContext) interface{} {
	reqLog(ctx)
	defer rspnsLog(ctx)
	return c.module.GetFPList()
}

//GetFPList 获取本机的指定文件的指纹信息，仅master提供对外查询功能
func (c *cnfs) GetFP(ctx context.IContext) interface{} {
	reqLog(ctx)
	defer rspnsLog(ctx)
	fp, err := c.module.GetLocalFP(ctx.Request().GetString("name"))
	if err != nil {
		return err
	}
	return fp
}

//RecvNotify 接收远程文件通知
func (c *cnfs) RecvNotify(ctx context.IContext) interface{} {
	reqLog(ctx)
	defer rspnsLog(ctx)
	fp := make(eFileFPLists)
	buff, err := ctx.Request().GetBody()
	if err != nil {
		return err
	}
	if err := ctx.Request().ToStruct(fp); err != nil {
		return err
	}
	ctx.Log().Debug("处理通知信息:", fp, string(buff))
	c.module.RecvNotify(fp)
	return "success"
}

//Download 用户下载文件
func (c *cnfs) GetFile(ctx context.IContext) interface{} {
	reqLog(ctx)
	defer rspnsLog(ctx)
	//根据路径查询文件
	path := ctx.Request().GetString("name")
	ctx.Log().Debugf("读取本地文件：%s", path)
	_, err := c.module.GetLocalFile(path)
	if err != nil {
		return err
	}
	ctx.Response().File(path, http.FS(c.module.local))
	return nil
}
