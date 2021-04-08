package nfs

import (
	"github.com/micro-plat/hydra/context"
)

//GetFPList 获取每个机器所有文件
func (c *cnfs) GetFPList(ctx context.IContext) interface{} {
	return c.module.GetFPList()
}

//GetFPList 获取本机的指定文件的指纹信息，仅master提供对外查询功能
func (c *cnfs) GetFP(ctx context.IContext) interface{} {
	fp, err := c.module.GetLocalFP(ctx.Request().GetString("name"))
	if err != nil {
		return err
	}
	return fp
}

//RecvNotify 接收远程文件通知
func (c *cnfs) RecvNotify(ctx context.IContext) interface{} {
	fp := &eFileFP{}
	if err := ctx.Request().Bind(fp); err != nil {
		return err
	}
	c.module.RecvNotify(fp)
	return "success"
}

//Download 用户下载文件
func (c *cnfs) GetFile(ctx context.IContext) interface{} {
	//根据路径查询文件
	path := ctx.Request().GetString("name")
	buff, err := c.module.GetFile(path)
	if err != nil {
		return err
	}
	return buff
}
