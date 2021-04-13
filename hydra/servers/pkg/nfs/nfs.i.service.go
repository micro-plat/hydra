package nfs

import (
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra/context"
)

const fileName = "name"
const dirName = "dir"

//QueryFP 获取每个机器所有文件
func (c *cnfs) QueryFP(ctx context.IContext) interface{} {
	list := c.module.GetFPList()
	ctx.Response().AddSpecial("nfs")
	ctx.Response().AddSpecial(fmt.Sprintf("%d", len(list)))
	return list
}

//GetFPList 获取本机的指定文件的指纹信息，仅master提供对外查询功能
func (c *cnfs) GetFP(ctx context.IContext) interface{} {
	if err := ctx.Request().Check(fileName); err != nil {
		return err
	}
	fp, err := c.module.GetLocalFP(ctx.Request().GetString(fileName))
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
	if err := ctx.Request().Check(fileName); err != nil {
		return err
	}
	path := ctx.Request().GetString(fileName)
	_, err := c.module.GetLocalFile(path)
	if err != nil {
		return err
	}
	ctx.Response().File(path, http.FS(c.module.local))
	return nil
}
