package nfs

import (
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/services"
)

func upload(context.IContext) interface{} {
	return nil
}
func init() {
	services.Def.OnSetup(func(c app.IAPPConf) error {
		nfs, err := c.GetNFSConf()
		if err != nil {
			return err
		}
		if nfs.Disable {
			return nil
		}
		services.Def.API("/file/upload", upload)
		services.Def.Web("/file/upload", upload)
		return nil
	}, global.API, global.Web)

	services.Def.OnClosing(func(c app.IAPPConf) error {
		return nil
	}, global.API, global.Web)
}
