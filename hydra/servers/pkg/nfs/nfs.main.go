package nfs

import (
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

var allCnfs = map[string]*cnfs{}

func init() {

	global.OnReady(func() {
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
				services.Def.API(SVSDonwload, cnfs.Download)
				services.Def.API(SVSUpload, cnfs.Upload)

				//内部服务
				services.Def.API(rmt_fp_get, cnfs.GetFP)
				services.Def.API(rmt_fp_notify, cnfs.RecvNotify)
				services.Def.API(rmt_fp_query, cnfs.QueryFP)
				services.Def.API(rmt_file_download, cnfs.GetFile)
			}

			if c.GetServerConf().GetServerType() == global.Web {
				services.Def.Web(SVSDonwload, cnfs.Download)
				services.Def.Web(SVSUpload, cnfs.Upload)

				//内部服务
				services.Def.Web(rmt_fp_get, cnfs.GetFP)
				services.Def.Web(rmt_fp_notify, cnfs.RecvNotify)
				services.Def.Web(rmt_fp_query, cnfs.QueryFP)
				services.Def.Web(rmt_file_download, cnfs.GetFile)
			}

			return nil

		}, global.API, global.Web)

		//处理服务启动完成
		services.Def.OnStarted(func(c app.IAPPConf) error {
			if cnfs, ok := allCnfs[c.GetServerConf().GetServerType()]; ok {
				err := cnfs.Start()
				if err != nil {
					return err
				}
				global.Def.Log().Info("nfs启动成功...")
				return nil
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

	})

}
