package nfs

import (
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

var allCnfs = map[string]*cnfs{}

func init() {

	//处理服务初始化
	services.Def.OnSetup(func(c app.IAPPConf) error {
		n, err := c.GetNFSConf()
		if err != nil {
			return err
		}
		if n.Disable {
			return nil
		}

		//构建对象
		cnfs := newNFS(c, n)
		allCnfs[c.GetServerConf().GetServerType()] = cnfs

		if c.GetServerConf().GetServerType() == global.API {
			//注册服务
			services.Def.API(SVS_Donwload, cnfs.Download)
			services.Def.API(SVS_Upload, cnfs.Upload)
			//内部服务
			services.Def.API(rmt_fp_get, cnfs.GetFP)
			services.Def.API(rmt_fp_push, cnfs.RecvNotify)
			services.Def.API(rmt_fp_list, cnfs.GetFPList)
			services.Def.API(rmt_file_pull, cnfs.GetFile)
		}

		if c.GetServerConf().GetServerType() == global.Web {
			services.Def.Web(SVS_Donwload, cnfs.Download)
			services.Def.Web(SVS_Upload, cnfs.Upload)

			//内部服务
			services.Def.Web(rmt_fp_get, cnfs.GetFP)
			services.Def.Web(rmt_fp_push, cnfs.RecvNotify)
			services.Def.Web(rmt_fp_list, cnfs.GetFPList)
			services.Def.Web(rmt_file_pull, cnfs.GetFile)
		}

		return nil

	}, global.API, global.Web)

	//处理服务启动完成
	services.Def.OnStarted(func(c app.IAPPConf) error {
		n, err := c.GetNFSConf()
		if err != nil {
			return err
		}
		if n.Disable {
			return nil
		}
		if cnfs, ok := allCnfs[c.GetServerConf().GetServerType()]; ok {
			return cnfs.Start()
		}
		return nil
	})

	//处理服务关闭
	services.Def.OnClosing(func(c app.IAPPConf) error {
		if cnfs, ok := allCnfs[c.GetServerConf().GetServerType()]; ok {
			cnfs.Close()
			delete(allCnfs, c.GetServerConf().GetServerType())
		}
		return nil
	})
}
